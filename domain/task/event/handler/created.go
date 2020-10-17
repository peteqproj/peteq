package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// CreatedHandler to handle task.created event
	CreatedHandler struct {
		Repo *task.Repo
	}
)

// Handle will handle the event the process it
func (c *CreatedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := task.Task{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return c.Repo.Create(ev.Tenant.ID, opt)
}

func (c *CreatedHandler) Name() string {
	return "task_domain_CreatedHandler"
}
