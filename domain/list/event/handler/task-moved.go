package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// TaskMovedHandler to handle list.task-moved event
	TaskMovedHandler struct {
		Repo *list.Repo
	}
	// TaskMovedSpec is the event.spec for this event
	TaskMovedSpec struct {
		TaskID      string `json:"taskId" yaml:"taskId"`
		Source      string `json:"source" yaml:"source"`
		Destination string `json:"destination" yaml:"destination"`
	}
)

// Handle will process it the event
func (t *TaskMovedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := TaskMovedSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return t.Repo.MoveTask(ev.Tenant.ID, opt.Source, opt.Destination, opt.TaskID)
}

func (t *TaskMovedHandler) Name() string {
	return "domain_TaskMovedHandler"
}
