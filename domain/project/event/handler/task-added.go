package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// TaskAddedHandler to handle task.created event
	TaskAddedHandler struct {
		Repo *project.Repo
	}
)

// Handle will process it the event
func (t *TaskAddedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := command.AddTasksCommandOptions{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return t.Repo.AddTask(ev.Tenant.ID, opt.Project, opt.TaskID)
}

func (t *TaskAddedHandler) Name() string {
	return "domain_TaskAddedHandler"
}
