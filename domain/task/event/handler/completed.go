package handler

import (
	"fmt"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// CompletedHandler to handle task.created event
	CompletedHandler struct {
		Repo *task.Repo
	}
)

// Handle will handle the event the process it
func (c *CompletedHandler) Handle(ev event.Event) error {
	task, err := c.Repo.Get(ev.Tenant.ID, ev.Metadata.AggregatorID)
	if err != nil {
		return fmt.Errorf("Failed to get task %s: %v", ev.Metadata.AggregatorID, err)
	}
	task.Status.Completed = true
	return c.Repo.Update(ev.Tenant.ID, task)
}

func (c *CompletedHandler) Name() string {
	return "domain_CompletedHandler"
}
