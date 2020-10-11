package handler

import (
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// TaskAddedHandler to handle task.created event
	TaskAddedHandler struct {
		Repo *project.Repo
	}
)

// Handle will process it the event
func (t *TaskAddedHandler) Handle(ev event.Event) error {
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
