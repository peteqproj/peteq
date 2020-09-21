package handler

import (
	"fmt"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// ReopenedHandler to handle task.created event
	ReopenedHandler struct {
		Repo *task.Repo
	}
)

// Handle will handle the event the process it
func (c *ReopenedHandler) Handle(ev event.Event) error {
	task, err := c.Repo.Get(ev.Metadata.AggregatorID)
	if err != nil {
		return fmt.Errorf("Failed to get task %s: %v", ev.Metadata.AggregatorID, err)
	}
	task.Status.Completed = false
	return c.Repo.Update(task)
}
