package nats

import (
	"sync"

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
		handlers map[string][]handler.EventHandler
	}
)

// Publish event
func (e *Eventbus) Publish(ev event.Event, done chan<- error) string {
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
		e.handlers = map[string][]handler.EventHandler{}
	}
	if _, ok := e.handlers[name]; ok {
		e.handlers[name] = append(e.handlers[name], h)
		return
	}

	e.handlers[name] = []handler.EventHandler{h}
	e.Store.QueueSubscribe(name, "worker", e.onMsg(name))
}

func (e *Eventbus) onMsg(name string) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		event := event.FromBytes(msg.Data)
		for _, h := range e.handlers[name] {
			e.Logger.Info("Recieved event -> calling handler", "event", event.Metadata.Name, "handler", h.Name())
			err := h.Handle(event)
			if err != nil {
				e.Logger.Info("Failed to handle event", "event", event.Metadata.Name, "error", err.Error(), "subscriber", h.Name())
			}
		}
	}
}
