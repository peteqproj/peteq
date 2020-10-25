package bus

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	socketio "github.com/googollee/go-socket.io"
	"github.com/peteqproj/peteq/pkg/db/local"
	"github.com/peteqproj/peteq/pkg/event"
	localbus "github.com/peteqproj/peteq/pkg/event/bus/local"
	"github.com/peteqproj/peteq/pkg/event/bus/rabbitmq"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// Eventbus to publish and subscribe events
	Eventbus interface {
		Publish(ctx context.Context, ev event.Event, done chan<- error) string
		Subscribe(name string, handler handler.EventHandler)
		Start() error
		Stop()
		Replay(ctx context.Context) error
	}

	// Options to create eventbus
	Options struct {
		Type            string
		LocalEventStore *local.DB
		WS              *socketio.Server
		Logger          logger.Logger
		EventlogDB      *sql.DB
		RabbitMQ        RabbitMQOptions
	}

	// Options to initiate rabbitmq
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
	if options.Type == "local" {
		return &localbus.Eventbus{
			Store:       options.LocalEventStore,
			Subscribers: map[string]chan<- localbus.EventChan{},
			Lock:        &sync.Mutex{},
			WS:          options.WS,
			Logger:      options.Logger,
		}, nil
	}

	if options.Type == "rabbitmq" {
		return &rabbitmq.Eventbus{
			Lock:             &sync.Mutex{},
			Logger:           options.Logger,
			Handlers:         map[string][]handler.EventHandler{},
			EventlogDB:       options.EventlogDB,
			RabbitMQHost:     options.RabbitMQ.Host,
			RabbitMQPassword: options.RabbitMQ.Password,
			RabbitMQUsername: options.RabbitMQ.Username,
			RabbitMQPort:     options.RabbitMQ.Port,
			RabbitMQAPIPort:  options.RabbitMQ.APIPort,
		}, nil
	}

	return nil, errors.New("Not found")
}
