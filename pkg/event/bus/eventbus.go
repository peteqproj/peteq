package bus

import (
	"sync"

	socketio "github.com/googollee/go-socket.io"
	"github.com/peteqproj/peteq/pkg/db/local"
	"github.com/peteqproj/peteq/pkg/event"
	localbus "github.com/peteqproj/peteq/pkg/event/bus/local"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// Eventbus to publish and subscribe events
	Eventbus interface {
		Publish(ev event.Event, done chan<- error) string
		Subscribe(name string, handler handler.EventHandler)
	}

	// Options to create eventbus
	Options struct {
		Type            string
		LocalEventStore *local.DB
		WS              *socketio.Server
		Logger          logger.Logger
	}
)

// New is factory for eventbus
func New(options Options) Eventbus {
	if options.Type == "local" {
		return &localbus.Eventbus{
			Store:       options.LocalEventStore,
			Subscribers: map[string]chan<- localbus.EventChan{},
			Lock:        &sync.Mutex{},
			WS:          options.WS,
			Logger:      options.Logger,
		}
	}

	return nil
}
