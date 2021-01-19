package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	retry "github.com/avast/retry-go"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/streadway/amqp"
)

const dbName = "event_log"

type (
	// Eventbus nats
	Eventbus struct {
		Logger            logger.Logger
		Lock              *sync.Mutex
		Handlers          map[string][]handler.EventHandler
		Channel           *amqp.Channel
		RabbitMQHost      string
		RabbitMQPort      string
		RabbitMQAPIPort   string
		RabbitMQUsername  string
		RabbitMQPassword  string
		IDGenerator       utils.IDGenerator
		WatchQueues       bool
		EventStorage      EventStorage
		ExtendContextFunc func(context.Context, event.Event) context.Context
	}
	EventStorage interface {
		Persist(context.Context, event.Event) error
	}
)

// Publish event
func (e *Eventbus) Publish(ctx context.Context, ev event.Event) (string, error) {
	e.Logger.Info("Publishing event", "name", ev.Metadata.Name)
	if err := e.ensureQueue(e.getKey(ev), false); err != nil {
		return "", fmt.Errorf("Failed to ensure queue: %w", err)
	}
	id, err := e.IDGenerator.GenerateV4()
	if err != nil {
		e.Logger.Info("Failed to create event id", "error", err.Error())
		return "", err
	}
	ev.Metadata.ID = id
	bytes, err := json.Marshal(ev)
	if err := e.EventStorage.Persist(ctx, ev); err != nil {
		e.Logger.Info("Failed to persist event", "error", err.Error())
		return "", err
	}
	if err := e.publish(ev.Metadata.Name, e.getKey(ev), bytes); err != nil {
		e.Logger.Info("Failed to publish event", "event", ev.Metadata.Name, "error", err.Error())
		return "", err
	}
	e.Logger.Info("Published", "event", ev.Metadata.Name)
	return id, nil
}

// Subscribe to event
// should be called with go Subscribe as this function is creating
// a channel and waits on it to receive event in order
// to call the handler
func (e *Eventbus) Subscribe(name string, h handler.EventHandler) {
	e.Lock.Lock()
	defer e.Lock.Unlock()
	if _, ok := e.Handlers[name]; ok {
		e.Logger.Info("Similar handler with same name exists, adding to set", "name", name, "handler", h.Name())
		e.Handlers[name] = append(e.Handlers[name], h)
		return
	}
	e.Handlers[name] = []handler.EventHandler{h}
}

func (e *Eventbus) Start() error {
	err := e.start()
	if err != nil {
		return err
	}
	if e.WatchQueues {
		go e.watchQueues()
	}
	return nil
}

func (e *Eventbus) Stop() {
	e.Logger.Info("Stopping eventbus")
	if e.Channel != nil {
		if err := e.Channel.Close(); err != nil {
			e.Logger.Info("Failed to close rabbitmq channel", "error", err.Error())
		}
	}
}

func (e *Eventbus) Replay(ctx context.Context) error {
	return nil
}

func (e *Eventbus) start() error {
	u := fmt.Sprintf("amqp://%s:%s@%s:%s", e.RabbitMQUsername, e.RabbitMQPassword, e.RabbitMQHost, e.RabbitMQPort)
	client, err := amqp.Dial(u)
	if err != nil {
		return err
	}
	ch, err := client.Channel()
	if err != nil {
		return err
	}
	e.Channel = ch
	return nil
}

func (e *Eventbus) publish(name string, key string, data []byte) error {
	defaultExchange := ""
	err := e.Channel.Publish(defaultExchange, key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	})
	if err != nil {
		return err
	}
	return nil
}

func (e *Eventbus) getKey(ev event.Event) string {
	return fmt.Sprintf("user-%s", ev.Tenant.ID)
}

func (e *Eventbus) queueList() map[string]bool {
	type Queue struct {
		Name  string `json:name`
		VHost string `json:vhost`
	}
	res := map[string]bool{}
	manager := fmt.Sprintf("http://%s:%s/api/queues/", e.RabbitMQHost, e.RabbitMQAPIPort)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", manager, nil)
	req.SetBasicAuth(e.RabbitMQUsername, e.RabbitMQPassword)
	resp, err := client.Do(req)
	if err != nil {
		e.Logger.Info("Failed to request RabbitMQ API server", "error", err.Error())
		return res
	}

	value := make([]Queue, 0)
	json.NewDecoder(resp.Body).Decode(&value)
	for _, q := range value {
		res[q.Name] = true
	}
	return res
}

func (e *Eventbus) watchQueues() {
	knownQueues := map[string]bool{}
	for {
		e.Logger.Info("Listing queues")
		list := e.queueList()
		for name := range list {
			e.Logger.Info("Testing queue for consumptions", "name", name)
			if _, ok := knownQueues[name]; !ok {
				e.Logger.Info("Queue added, starting to watch", "name", name)
				replayQueue := false
				msgs, err := e.Channel.Consume(name, "", false, replayQueue, false, false, nil)
				if err != nil {
					e.Logger.Info("Failed to Consume", "error", err.Error())
					continue
				}
				knownQueues[name] = true
				go e.watchQueue(msgs, e.Logger.Fork("queue", name))
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (e *Eventbus) watchQueue(ch <-chan amqp.Delivery, lgr logger.Logger) {

	for msg := range ch {
		ev := event.Event{}
		err := json.Unmarshal(msg.Body, &ev)
		if err != nil {
			lgr.Info("Failed to unmarshal into event", "error", err.Error())
			continue
		}
		lgr.Info("Received event", "name", ev.Metadata.Name)
		set, ok := e.Handlers[ev.Metadata.Name]
		if !ok {
			lgr.Info("Handler not found", "event", ev.Metadata.Name)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, h := range set {
			wg.Add(1)
			lgr.Info("Calling event handler", "name", h.Name())
			log := logger.New(logger.Options{})
			elgr := log.Fork("event", ev.Metadata.Name, "id", ev.Metadata.ID, "handler", h.Name())
			go e.deliverEvent(context.Background(), wg, ev, h, elgr)
		}
		wg.Wait()
		if err := e.Channel.Ack(msg.DeliveryTag, false); err != nil {
			e.Logger.Info("Failed to ack event", "error", err.Error())
		}
	}
}

func (e *Eventbus) deliverEvent(ctx context.Context, wg *sync.WaitGroup, ev event.Event, handler handler.EventHandler, lgr logger.Logger) {
	delay := retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
		lgr.Info("Failed to handle event, retrying")
		return retry.BackOffDelay(n, err, config)
	})
	if err := retry.Do(func() error {
		return handler.Handle(e.ExtendContextFunc(ctx, ev), ev, lgr)
	}, retry.Delay(1*time.Second), retry.Attempts(5), delay); err != nil {
		lgr.Info("Failed to handle event", "error", err.Error())
	}
	wg.Done()
}

func (e *Eventbus) ensureQueue(name string, replayQueue bool) error {
	durable := false
	autoDelete := true
	exclusive := false
	if replayQueue == true {
		autoDelete = false
		exclusive = false
	}
	noWait := false
	e.Logger.Info("Creating queue", "autoDelete", autoDelete, "exclusive", exclusive)
	_, err := e.Channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, nil)
	if err != nil {
		return fmt.Errorf("Failed to QueueDeclare: %w", err)
	}
	return nil
}
