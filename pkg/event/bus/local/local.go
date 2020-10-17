package local

import (
	"context"
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
	socketio "github.com/googollee/go-socket.io"
	"github.com/peteqproj/peteq/pkg/db/local"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
	"gopkg.in/yaml.v2"
)

type (
	// Eventbus local, in-memory eventbus
	Eventbus struct {
		Store       *local.DB
		Subscribers map[string]chan<- EventChan
		Lock        *sync.Mutex
		WS          *socketio.Server
		Logger      logger.Logger
	}

	// EventChan is pair between event and a channel to report done
	EventChan struct {
		event event.Event
		done  chan<- error
	}
)

// Publish event
func (e *Eventbus) Publish(ctx context.Context, ev event.Event, done chan<- error) string {

	eID, err := uuid.NewV4()
	if err != nil {
		done <- fmt.Errorf("Failed to generate event-id: %w", err)
		return ""
	}
	ev.Metadata.ID = eID.String()

	eventBytes, err := e.Store.Read()
	if err != nil {
		done <- fmt.Errorf("Failed to read current eventlog: %w", err)
		return ""
	}
	events := []event.Event{}
	if err := yaml.Unmarshal(eventBytes, &events); err != nil {
		done <- err
		return ""
	}
	events = append(events, ev)
	bytes, err := yaml.Marshal(events)
	if err != nil {
		done <- fmt.Errorf("Failed to marshal event: %w", err)
		return ""
	}
	if err := e.Store.Write(bytes); err != nil {
		done <- fmt.Errorf("Failed to store event: %w", err)
		return ""
	}
	e.Lock.Lock()
	defer e.Lock.Unlock()
	for name, subscriber := range e.Subscribers {
		if name == ev.Metadata.Name {
			e.Logger.Info("Publishing event", "name", ev.Metadata.Name, "tenantId", ev.Tenant.ID)
			subscriber <- EventChan{
				event: ev,
				done:  done,
			}
		}
	}
	for _, room := range e.WS.Rooms("/") {
		e.WS.BroadcastToRoom("/", room, ev.Metadata.Name, ev)
	}
	return eID.String()
}

// Subscribe to event
// should be called with go Subscribe as this function is creating
// a channel and waits on it to receive event in order
// to call the handler
func (e *Eventbus) Subscribe(name string, h handler.EventHandler) {
	e.Lock.Lock()
	cn := make(chan EventChan)
	e.Subscribers[name] = cn
	e.Lock.Unlock()
	log := logger.New(logger.Options{})
	for {
		select {
		case e := <-cn:
			if err := h.Handle(context.Background(), e.event, log); err != nil {
				e.done <- err
			}
			e.done <- nil
		}
	}
}

func (e *Eventbus) Start() error {
	return nil
}
func (e *Eventbus) Stop() {
}
func (e *Eventbus) Replay(ctx context.Context) error {
	return nil
}
