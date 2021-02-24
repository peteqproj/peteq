package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/repo"
)

type (
	// TaskAddedHandler to handle project.task-added event
	TaskAddedHandler struct {
		Repo *repo.Repo
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
	prj, err := t.Repo.Get(ctx, repo.GetOptions{ID: opt.Project})
	if err != nil {
		return err
	}
	if s, ok := prj.Spec.(project.Spec); ok {
		s.Tasks = append(s.Tasks, opt.TaskID)
	}
	return t.Repo.Update(ctx, *prj)
}

func (t *TaskAddedHandler) Name() string {
	return "domain_list"
}
