package handler

import (
	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/list/command"
	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// TaskMovedHandler to handle list.task-moved event
	TaskMovedHandler struct {
		Repo *list.Repo
	}
)

// Handle will process it the event
func (t *TaskMovedHandler) Handle(ev event.Event) error {
	opt := command.MoveTaskArguments{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return t.Repo.MoveTask(ev.Tenant.ID, opt.Source, opt.Destination, opt.TaskID)
}

func (t *TaskMovedHandler) Name() string {
	return "domain_TaskMovedHandler"
}
