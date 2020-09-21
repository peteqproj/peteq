package handler

import (
	"fmt"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// UpdatedHandler to handle task.created event
	UpdatedHandler struct {
		Repo *task.Repo
	}
)

// Handle will handle the event the process it
func (c *UpdatedHandler) Handle(ev event.Event) error {
	t, ok := ev.Spec.(task.Task)
	if !ok {
		return fmt.Errorf("Failed to cast to task object")
	}
	return c.Repo.Update(t)
}
