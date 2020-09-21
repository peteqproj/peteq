package handler

import (
	"fmt"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// DeleteHandler to handle task.deleted event
	DeleteHandler struct {
		Repo *task.Repo
	}
)

// Handle will handle the event the process it
func (c *DeleteHandler) Handle(ev event.Event) error {
	t, ok := ev.Spec.(task.Task)
	if !ok {
		return fmt.Errorf("Failed to cast to task object")
	}
	return c.Repo.Delete(t.Metadata.ID)
}
