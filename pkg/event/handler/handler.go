package handler

import (
	"context"

	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// EventHandler to handle events once occoured
	EventHandler interface {
		Handle(ctx context.Context, ev event.Event, logger logger.Logger) error
		Name() string
	}
)
