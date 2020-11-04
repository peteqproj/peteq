package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/trigger"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// TriggeredHandler to handle trigger.triggered event
	TriggeredHandler struct {
		Repo *trigger.Repo
	}

	// TriggeredSpec is the event.spec for this event
	TriggeredSpec struct {
		ID string `json:"id" yaml:"id"`
	}
)

// Handle will process it the event
func (t *TriggeredHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := TriggeredSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return nil
}

func (t *TriggeredHandler) Name() string {
	return "domain_trigger"
}
