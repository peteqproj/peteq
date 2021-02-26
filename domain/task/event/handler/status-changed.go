package handler

import (
	"context"
	"fmt"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/repo"
)

type (
	// StatusChangedHandler to handle task.created event
	StatusChangedHandler struct {
		Repo *repo.Repo
	}

	// StatusChangedSpec is the event.spec for this event
	StatusChangedSpec struct {
		Completed bool `json:"completed"`
	}
)

// Handle will handle the event the process it
func (c *StatusChangedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := StatusChangedSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	t, err := c.Repo.Get(ctx, repo.GetOptions{ID: ev.Metadata.AggregatorID})
	if err != nil {
		return fmt.Errorf("Failed to get task %s: %v", ev.Metadata.AggregatorID, err)
	}
	spec := t.Spec.(task.Spec)
	spec.Completed = opt.Completed
	t.Spec = spec
	return c.Repo.Update(ctx, *t)
}

func (c *StatusChangedHandler) Name() string {
	return "domain_task"
}
