package google

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/peteqproj/peteq/pkg/command/handler"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	Bus struct {
		Ps                *pubsub.Client
		Topic             string
		IDGenerator       utils.IDGenerator
		Handlers          map[string]handler.CommandHandler
		Lock              *sync.Mutex
		Logger            logger.Logger
		ExtendContextFunc func(context.Context, string) context.Context
		Subscribtion      string
	}
)

// Start local do nothing
func (b *Bus) Start() error {
	b.Logger.Info("Starting commandbus")
	sub := b.Ps.Subscription(b.Subscribtion)
	t, err := b.createTopicIfNotExists()
	b.Logger.Info("Topic created")
	if err != nil {
		return err
	}
	if err := b.createSubscriptionIfNotExists(b.Subscribtion, t); err != nil {
		return err
	}
	b.Logger.Info("Starting to receive messages")
	go func() {
		err = sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
			msg.Ack()
			b.Logger.Info("Received message", "msg", string(msg.Data))
			name, exists := msg.Attributes["handler"]
			if !exists {
				b.Logger.Info("Command handler not set")
				return
			}
			h, exists := b.Handlers[name]
			if !exists {
				b.Logger.Info("Command handler not found")
				return
			}
			uid, exists := msg.Attributes["user"]
			if !exists {
				b.Logger.Info("UserID is not set", "handler", name)
			}
			if err := h.Handle(b.ExtendContextFunc(ctx, uid), string(msg.Data)); err != nil {
				b.Logger.Info("Failed to handler command", "handler", name, "error", err.Error())
			}
		})
		if err != nil {
			b.Logger.Info("Failed to start receiving", "error", err.Error())
		}
	}()
	return nil
}

// Execute runs command handler
func (b *Bus) Execute(ctx context.Context, name string, arguments interface{}) error {
	id, err := b.IDGenerator.GenerateV4()
	if err != nil {
		return fmt.Errorf("Failed to generate UUID: %w", err)
	}
	data, err := json.Marshal(arguments)
	if err != nil {
		return err
	}
	u := tenant.UserFromContext(ctx)
	uid := ""
	if u != nil {
		uid = u.Metadata.ID
	}
	res := b.Ps.Topic(b.Topic).Publish(ctx, &pubsub.Message{
		ID:   id,
		Data: data,
		Attributes: map[string]string{
			"handler": name,
			"user":    uid,
		},
	})
	sID, err := res.Get(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get response from Google: %w", err)
	}
	b.Logger.Info("Command published", "id", id, "server-id", sID, "command", name)
	return nil
}
func (b *Bus) RegisterHandler(name string, ch handler.CommandHandler) error {
	b.Lock.Lock()
	defer b.Lock.Unlock()
	_, ok := b.Handlers[name]
	if ok {
		return fmt.Errorf("Handler already exist")
	}
	b.Handlers[name] = ch
	return nil
}

func (b *Bus) createTopicIfNotExists() (*pubsub.Topic, error) {
	ctx := context.Background()
	t := b.Ps.Topic(b.Topic)
	ok, err := t.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if ok {
		return t, nil
	}
	t, err = b.Ps.CreateTopic(ctx, b.Topic)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (b *Bus) createSubscriptionIfNotExists(name string, topic *pubsub.Topic) error {
	ctx := context.Background()
	_, err := b.Ps.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
		Topic:            topic,
		ExpirationPolicy: time.Duration(0),
		AckDeadline:      20 * time.Second,
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "rpc error: code = AlreadyExists") {
			return nil
		}
		return err
	}
	return nil
}
