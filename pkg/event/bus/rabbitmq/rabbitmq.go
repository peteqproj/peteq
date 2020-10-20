package rabbitmq

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/gofrs/uuid"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/streadway/amqp"
)

const dbName = "event_log"

type (
	// Eventbus nats
	Eventbus struct {
		Logger           logger.Logger
		Lock             *sync.Mutex
		Handlers         map[string][]handler.EventHandler
		Channel          *amqp.Channel
		EventlogDB       *sql.DB
		RabbitMQHost     string
		RabbitMQPort     string
		RabbitMQAPIPort  string
		RabbitMQUsername string
		RabbitMQPassword string
	}
)

// Publish event
func (e *Eventbus) Publish(ctx context.Context, ev event.Event, done chan<- error) string {
	e.Logger.Info("Publishing event", "name", ev.Metadata.Name)
	id, err := uuid.NewV4()
	if err != nil {
		e.Logger.Info("Failed to create event id", "error", err.Error())
		done <- err
		return ""
	}
	ev.Metadata.ID = id.String()
	bytes, err := json.Marshal(ev)
	if err := e.persistEvent(context.Background(), ev.Tenant.ID, id.String(), ev.Metadata.Name, string(bytes)); err != nil {
		e.Logger.Info("Failed to persist event", "error", err.Error())
		done <- err
		return ""
	}
	q := e.getQueue(ev)
	if err := e.publish(ev.Metadata.Name, q, bytes); err != nil {
		e.Logger.Info("Failed to publish event", "event", ev.Metadata.Name, "error", err.Error())
		done <- err
		return ""
	}
	e.Logger.Info("Published", "event", ev.Metadata.Name)
	done <- nil
	return id.String()
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
	go e.watchQueues()
	return nil
}

func (e *Eventbus) Stop() {
	e.Logger.Info("Stopping eventbus")
	if e.Channel != nil {
		if err := e.Channel.Close(); err != nil {
			e.Logger.Info("Failed to close rabbitmq channel", "error", err.Error())
		}
	}
	if e.EventlogDB != nil {
		if err := e.EventlogDB.Close(); err != nil {
			e.Logger.Info("Failed to close conneciton to event log database", "error", err.Error())
		}
	}
}

func (e *Eventbus) Replay(ctx context.Context) error {
	e.Logger.Info("Replaying")
	err := e.start()
	if err != nil {
		return err
	}
	user := ctx.Value("UserID")
	if user == nil {
		return fmt.Errorf("UserID was not passed in contxt")
	}
	ex := exp.Ex{
		"userid": user.(string),
	}
	q, _, err := goqu.
		From(dbName).
		Where(ex).
		ToSQL()
	if err != nil {
		return err
	}
	rows, err := e.EventlogDB.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	set := []event.Event{}
	for rows.Next() {
		id := ""
		name := ""
		user := ""
		ev := ""
		if err := rows.Scan(&id, &name, &user, &ev); err != nil {
			return err
		}
		event := event.Event{}
		if err := json.Unmarshal([]byte(ev), &event); err != nil {
			return err
		}
		set = append(set, event)
	}
	for _, ev := range set {
		e.Logger.Info("Processing event", "name", ev.Metadata.Name)
		q := e.getQueue(ev)
		data, err := json.Marshal(ev)
		if err != nil {
			return err
		}
		if err := e.publish(ev.Metadata.Name, q, data); err != nil {
			return err
		}
	}
	return nil
}

func (e *Eventbus) start() error {
	u := fmt.Sprintf("amqp://%s:%s@%s:%s/", e.RabbitMQUsername, e.RabbitMQPassword, e.RabbitMQHost, e.RabbitMQPort)
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

func (e *Eventbus) publish(name string, queue string, data []byte) error {
	if err := e.ensureQueue(queue); err != nil {
		return fmt.Errorf("Failed to ensure queue: %w", err)
	}
	err := e.Channel.Publish("", queue, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	})
	if err != nil {
		return err
	}
	return nil
}

func (e *Eventbus) getQueue(ev event.Event) string {
	return fmt.Sprintf("user-%s", ev.Tenant.ID)
}

func (e *Eventbus) persistEvent(ctx context.Context, user string, id string, name string, data string) error {
	q, _, err := goqu.
		Insert(dbName).
		Cols("eventid", "eventname", "userid", "info").
		Vals(goqu.Vals{id, name, user, data}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = e.EventlogDB.ExecContext(ctx, q)
	if err != nil {
		return err
	}
	return nil
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
		list := e.queueList()
		for name := range list {
			if _, ok := knownQueues[name]; !ok {
				e.Logger.Info("Queue added, starting to watch", "name", name)
				knownQueues[name] = true
				go e.watchQueue(name, e.Logger.Fork("queue", name))
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (e *Eventbus) watchQueue(q string, lgr logger.Logger) {
	if err := e.ensureQueue(q); err != nil {
		lgr.Info(err.Error())
		// TODO: use channel to report back error
		return
	}
	if err := e.ensureQueue(q); err != nil {
		e.Logger.Info("Failed to ensure channel", "error", err.Error())
	}
	msgs, err := e.Channel.Consume(q, "", false, true, false, false, nil)
	if err != nil {
		lgr.Info("Failed to Consume", "error", err.Error())
	}
	for msg := range msgs {
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
		lgr.Info("Calling event handlers", "len", len(set))
		for _, h := range set {
			lgr.Info("Calling event handler", "name", h.Name())
			log := logger.New(logger.Options{})
			if err := h.Handle(context.Background(), ev, log.Fork("event", ev.Metadata.Name, "handler", h.Name())); err != nil {
				lgr.Info("Failed to handle event", "event", ev.Metadata.Name, "handler", h.Name(), "error", err.Error())
			}
		}
		if err := e.Channel.Ack(msg.DeliveryTag, false); err != nil {
			e.Logger.Info("Failed to ack evenet", "error", err.Error())
		}
	}
}

func (e *Eventbus) ensureQueue(name string) error {
	durable := false
	autoDelete := true
	exclusive := false
	noWait := false
	_, err := e.Channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, nil)
	if err != nil {
		return fmt.Errorf("Failed to QueueDeclare: %w", err)
	}
	return nil
}
