package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// TaskAddedHandler to handle project.task-added event
	TaskAddedHandler struct {
		Repo *project.Repo
	}

	// TaskAddedSpec is the event.spec for this event
	TaskAddedSpec struct {
		TaskID  string `json:"taskId" yaml:"taskId"`
		Project string `json:"project" yaml:"project"`
	}
)

// Handle will process it the event
func (t *TaskAddedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := TaskAddedSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return t.Repo.AddTask(ev.Tenant.ID, opt.Project, opt.TaskID)
}

func (t *TaskAddedHandler) Name() string {
	return "domain_TaskAddedHandler"
}
