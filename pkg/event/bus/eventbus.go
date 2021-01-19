package bus

import (
	"context"
	"errors"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus/google"
	"github.com/peteqproj/peteq/pkg/event/bus/rabbitmq"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/event/storage"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// Eventbus to publish and subscribe events
	Eventbus interface {
		EventPublisher
		EventWatcher
		Start() error
		Stop()
		Replay(ctx context.Context) error
	}

	EventPublisher interface {
		Publish(ctx context.Context, ev event.Event) (string, error)
	}

	EventWatcher interface {
		Subscribe(name string, handler handler.EventHandler)
	}

	// Options to create eventbus
	Options struct {
		Logger            logger.Logger
		RabbitMQ          *RabbitMQOptions
		GooglePubSub      *GooglePubSubOptions
		ExtendContextFunc func(context.Context, event.Event) context.Context
		EventStorage      *storage.Storage
	}

	// RabbitMQOptions to initiate rabbitmq
	RabbitMQOptions struct {
		Host        string
		Port        string
		APIPort     string
		Username    string
		Password    string
		WatchQueues bool
	}
	// GooglePubSubOptions to initiate Google pubsub
	GooglePubSubOptions struct {
		Client *pubsub.Client
	}

	// ReplayOptions options to replay events
	ReplayOptions struct {
		User string
	}
)

// New is factory for eventbus
func New(options Options) (Eventbus, error) {

	if options.RabbitMQ != nil {
		return &rabbitmq.Eventbus{
			Lock:              &sync.Mutex{},
			Logger:            options.Logger,
			Handlers:          map[string][]handler.EventHandler{},
			RabbitMQHost:      options.RabbitMQ.Host,
			RabbitMQPassword:  options.RabbitMQ.Password,
			RabbitMQUsername:  options.RabbitMQ.Username,
			RabbitMQPort:      options.RabbitMQ.Port,
			RabbitMQAPIPort:   options.RabbitMQ.APIPort,
			IDGenerator:       utils.NewGenerator(),
			WatchQueues:       options.RabbitMQ.WatchQueues,
			ExtendContextFunc: options.ExtendContextFunc,
			EventStorage:      options.EventStorage,
		}, nil
	}

	if options.GooglePubSub != nil {
		return &google.Eventbus{
			Logger:            options.Logger,
			ExtendContextFunc: options.ExtendContextFunc,
			IDGenerator:       utils.NewGenerator(),
			Ps:                options.GooglePubSub.Client,
			EventStorage:      options.EventStorage,
			Lock:              &sync.Mutex{},
			Handlers:          map[string][]handler.EventHandler{},
		}, nil
	}

	return nil, errors.New("Not found")
}
