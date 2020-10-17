package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/list/command"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// TaskMovedHandler to handle list.task-moved event
	TaskMovedHandler struct {
		Repo *list.Repo
	}
)

// Handle will process it the event
func (t *TaskMovedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
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
