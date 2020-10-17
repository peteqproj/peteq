package nats

import (
	"context"
	"fmt"
	"sync"
	"time"

	stan "github.com/nats-io/stan.go"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// Eventbus nats
	Eventbus struct {
		Store    stan.Conn
		Logger   logger.Logger
		Lock     *sync.Mutex
		handlers []eventHandlers
		started  bool
	}

	eventHandlers struct {
		name     string
		handlers []handler.EventHandler
	}
)

// Publish event
func (e *Eventbus) Publish(ctx context.Context, ev event.Event, done chan<- error) string {
	e.Logger.Info("Publishing event", "name", ev.Metadata.Name)
	err := e.Store.Publish(ev.Metadata.Name, ev.ToBytes())
	if err != nil {
		e.Logger.Info("Failed to publish event", "name", ev.Metadata.Name, "error", err.Error())
		done <- err
		return ""
	}
	done <- nil
	return ""
}

// Subscribe to event
// should be called with go Subscribe as this function is creating
// a channel and waits on it to receive event in order
// to call the handler
func (e *Eventbus) Subscribe(name string, h handler.EventHandler) {
	e.Lock.Lock()
	defer e.Lock.Unlock()
	if e.handlers == nil {
		e.handlers = []eventHandlers{}
	}
	handlerIndex := -1
	for i, h := range e.handlers {
		if h.name == name {
			handlerIndex = i
			break
		}
	}
	if handlerIndex != -1 {
		e.Logger.Info("Similar handler with same name exists, adding to set", "name", name, "handler", h.Name())
		e.handlers[handlerIndex].handlers = append(e.handlers[handlerIndex].handlers, h)
		return
	}
	e.handlers = append(e.handlers, eventHandlers{
		name:     name,
		handlers: []handler.EventHandler{h},
	})
}

func (e *Eventbus) Start() error {
	if e.started {
		return fmt.Errorf("Already running")
	}
	d, _ := time.ParseDuration("2h")
	for _, h := range e.handlers {
		e.Logger.Info("Starting queue", "name", h.name)
		_, err := e.Store.QueueSubscribe(h.name, "worker", e.onMsg(h.name, h.handlers), stan.StartAtTimeDelta(d), stan.AckWait(time.Second*30), stan.MaxInflight(1))
		if err != nil {
			return err
		}
	}
	e.started = true
	return nil
}

func (e *Eventbus) Replay(ctx context.Context) error {
	var onMsg = func(msg *stan.Msg) {
		event := event.FromBytes(msg.Data)
		if event.Metadata.Name != "user.registred" && event.Tenant.ID != "ae1703ec-51f3-4c01-b316-dea0fcf742fc" {
			e.Logger.Info("Skipping event", "tenant", event.Tenant.ID, "name", event.Metadata.Name)
			return
		}
		for _, he := range e.handlers {
			if he.name != event.Metadata.Name {
				continue
			}
			for _, h := range he.handlers {
				e.Logger.Info("Recieved event", "event", event.Metadata.Name, "handler", h.Name())
				log := logger.New(logger.Options{})
				err := h.Handle(context.Background(), event, log)
				if err != nil {
					e.Logger.Info("Failed to handle event", "event", event.Metadata.Name, "error", err.Error(), "subscriber", h.Name())
				}
			}
		}
	}
	d, _ := time.ParseDuration("3h")
	_, err := e.Store.Subscribe("", onMsg, stan.StartAtTimeDelta(d), stan.AckWait(time.Second*30), stan.MaxInflight(1))
	if err != nil {
		return err
	}
	return nil
}

func (e *Eventbus) Stop() {
}

func (e *Eventbus) onMsg(name string, handlers []handler.EventHandler) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		event := event.FromBytes(msg.Data)
		if event.Metadata.Name != "user.registred" && event.Tenant.ID != "ae1703ec-51f3-4c01-b316-dea0fcf742fc" {
			e.Logger.Info("Skipping event", "tenant", event.Tenant.ID, "name", event.Metadata.Name)
			return
		}
		for _, h := range handlers {
			e.Logger.Info("Recieved event", "event", event.Metadata.Name, "handler", h.Name())
			log := logger.New(logger.Options{})
			err := h.Handle(context.Background(), event, log)
			if err != nil {
				e.Logger.Info("Failed to handle event", "event", event.Metadata.Name, "error", err.Error(), "subscriber", h.Name())
			}
		}
		msg.Ack()
	}
}
