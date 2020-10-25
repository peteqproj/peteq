package handler

import (
	"context"
	"fmt"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// StatusChangedHandler to handle task.created event
	StatusChangedHandler struct {
		Repo *task.Repo
	}

	// StatusChangedSpec is the event.spec for this event
	StatusChangedSpec struct {
		Completed bool `json:"completed" yaml:"completed"`
	}
)

// Handle will handle the event the process it
func (c *StatusChangedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := StatusChangedSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	task, err := c.Repo.Get(ev.Tenant.ID, ev.Metadata.AggregatorID)
	if err != nil {
		return fmt.Errorf("Failed to get task %s: %v", ev.Metadata.AggregatorID, err)
	}
	task.Status.Completed = opt.Completed
	return c.Repo.Update(ev.Tenant.ID, task)
}

func (c *StatusChangedHandler) Name() string {
	return "domain_task"
}
