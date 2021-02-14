package bus

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
	"google.golang.org/api/iterator"
)

type (
	GoogleEventbus struct {
		Logger            logger.Logger
		Ps                *pubsub.Client
		Lock              *sync.Mutex
		Handlers          map[string][]handler.EventHandler
		EventStorage      EventStorage
		IDGenerator       utils.IDGenerator
		ExtendContextFunc func(context.Context, event.Event) context.Context
	}
)

func (e *GoogleEventbus) Publish(ctx context.Context, ev event.Event) (string, error) {
	id, err := e.IDGenerator.GenerateV4()
	if err != nil {
		e.Logger.Info("Failed to create event id", "error", err.Error())
		return "", err
	}
	ev.Metadata.ID = id
	if err := e.EventStorage.Persist(ctx, ev); err != nil {
		return "", fmt.Errorf("Failed to persist event: %w", err)
	}
	k := e.getTopic(ev)
	err, t := e.ensureTopic(k)
	if err != nil {
		return "", err
	}

	e.createSubscriptionIfNotExists(k, t)
	tenant := tenant.UserFromContext(ctx)
	data, err := json.Marshal(ev)
	if err != nil {
		return "", err
	}
	res := e.Ps.Topic(k).Publish(ctx, &pubsub.Message{
		ID: ev.Metadata.ID,
		Attributes: map[string]string{
			"tenant":  tenant.Metadata.ID,
			"handler": ev.Metadata.Name,
		},
		Data: data,
	})
	sID, err := res.Get(ctx)
	if err != nil {
		return "", err
	}

	e.Logger.Info("Event published", "id", ev.Metadata.ID, "server-id", sID, "name", ev.Metadata.Name)
	return "", nil
}

func (e *GoogleEventbus) Subscribe(name string, h handler.EventHandler) {
	e.Lock.Lock()
	defer e.Lock.Unlock()
	if _, ok := e.Handlers[name]; ok {
		e.Logger.Info("Similar handler with same name exists, adding to set", "name", name, "handler", h.Name())
		e.Handlers[name] = append(e.Handlers[name], h)
		return
	}
	e.Handlers[name] = []handler.EventHandler{h}

}

func (e *GoogleEventbus) Start() error {
	go e.watchSubscriptions()
	return nil
}

func (e *GoogleEventbus) Stop() {
}

func (e *GoogleEventbus) getTopic(ev event.Event) string {
	return fmt.Sprintf("user-%s", ev.Tenant.ID)
}
func (e *GoogleEventbus) watchSubscriptions() {
	knownSubscriptions := map[string]bool{}
	for {
		it := e.Ps.Subscriptions(context.Background())
		for {
			s, err := it.Next()
			if err == iterator.Done {
				e.Logger.Info("No more subscriptions found")
				break
			}
			if err != nil {
				e.Logger.Info("Failed to request next subscription", "error", err.Error())
				continue
			}
			if _, ok := knownSubscriptions[s.ID()]; !ok {
				knownSubscriptions[s.ID()] = true
				if strings.HasPrefix(s.ID(), "user") {
					e.Logger.Info("Starting to watch suscripsion", "name", s.ID())
					go e.watchSubscription(s)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}
func (e *GoogleEventbus) ensureTopic(name string) (error, *pubsub.Topic) {
	ctx := context.Background()
	t := e.Ps.Topic(name)
	ok, err := t.Exists(ctx)
	if err != nil {
		return err, nil
	}
	if ok {
		return nil, t
	}
	t, err = e.Ps.CreateTopic(ctx, name)
	if err != nil {
		return err, nil
	}
	return nil, t
}
func (e *GoogleEventbus) createSubscriptionIfNotExists(name string, topic *pubsub.Topic) error {
	ctx := context.Background()
	_, err := e.Ps.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
		Topic:                 topic,
		ExpirationPolicy:      time.Duration(0),
		AckDeadline:           20 * time.Second,
		EnableMessageOrdering: true,
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "rpc error: code = AlreadyExists") {
			return nil
		}
		return err
	}
	return nil
}

func (e *GoogleEventbus) watchSubscription(sub *pubsub.Subscription) {
	err := sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		e.Logger.Info("Received message", "msg", string(msg.Data))
		name, exists := msg.Attributes["handler"]
		if !exists {
			e.Logger.Info("Event handler not set")
			return
		}
		handlers, exists := e.Handlers[name]
		if !exists {
			e.Logger.Info("Event handler not found")
			return
		}
		ev := event.Event{}
		err := json.Unmarshal(msg.Data, &ev)
		if err != nil {
			e.Logger.Info("Failed to unmarshal into event", "error", err.Error())
			return
		}
		_, exists = msg.Attributes["tenant"]
		if !exists {
			e.Logger.Info("UserID is not set", "handler", name)
			return
		}
		for _, h := range handlers {
			e.Logger.Info("Calling event handler", "handler", h.Name())
			if err := h.Handle(e.ExtendContextFunc(ctx, ev), ev, e.Logger.Fork("event", ev.Metadata.Name)); err != nil {
				e.Logger.Info("Failed to handler event", "handler", name, "error", err.Error())
				return
			}
		}
	})
	if err != nil {
		e.Logger.Info("Failed to start receiving", "error", err.Error())
	}
}
