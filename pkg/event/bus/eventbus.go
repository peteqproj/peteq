package bus

import (
	"context"
	"errors"
	"sync"

	socketio "github.com/googollee/go-socket.io"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus/rabbitmq"
	"github.com/peteqproj/peteq/pkg/event/handler"
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
		Type              string
		WS                *socketio.Server
		Logger            logger.Logger
		EventlogDB        db.Database
		RabbitMQ          RabbitMQOptions
		WatchQueues       bool
		ExtendContextFunc func(context.Context, event.Event) context.Context
	}

	// RabbitMQOptions to initiate rabbitmq
	RabbitMQOptions struct {
		Host     string
		Port     string
		APIPort  string
		Username string
		Password string
	}

	// ReplayOptions options to replay events
	ReplayOptions struct {
		User string
	}
)

// New is factory for eventbus
func New(options Options) (Eventbus, error) {

	if options.Type == "rabbitmq" {
		return &rabbitmq.Eventbus{
			Lock:              &sync.Mutex{},
			Logger:            options.Logger,
			Handlers:          map[string][]handler.EventHandler{},
			EventlogDB:        options.EventlogDB,
			RabbitMQHost:      options.RabbitMQ.Host,
			RabbitMQPassword:  options.RabbitMQ.Password,
			RabbitMQUsername:  options.RabbitMQ.Username,
			RabbitMQPort:      options.RabbitMQ.Port,
			RabbitMQAPIPort:   options.RabbitMQ.APIPort,
			IDGenerator:       utils.NewGenerator(),
			WatchQueues:       options.WatchQueues,
			ExtendContextFunc: options.ExtendContextFunc,
		}, nil
	}

	return nil, errors.New("Not found")
}
