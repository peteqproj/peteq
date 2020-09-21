package handler

import "github.com/peteqproj/peteq/pkg/event"

type (
	// EventHandler to handle events once occoured
	EventHandler interface {
		Handle(ev event.Event) error
	}
)
