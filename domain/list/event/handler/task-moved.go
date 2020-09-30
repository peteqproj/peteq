package handler

import (
	"fmt"

	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/list/command"
	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// TaskMovedHandler to handle task.created event
	TaskMovedHandler struct {
		Repo *list.Repo
	}
)

// Handle will process it the event
func (t *TaskMovedHandler) Handle(ev event.Event) error {
	opt, ok := ev.Spec.(command.MoveTaskArguments)
	if !ok {
		return fmt.Errorf("Failed to cast to task object")
	}
	return t.Repo.MoveTask(ev.Tenant.ID, opt.Source, opt.Destination, opt.TaskID)
}
