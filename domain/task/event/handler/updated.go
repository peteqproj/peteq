package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// UpdatedHandler to handle task.created event
	UpdatedHandler struct {
		Repo *task.Repo
	}
)

// Handle will handle the event the process it
func (c *UpdatedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := task.Task{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return c.Repo.Update(ev.Tenant.ID, opt)
}

func (c *UpdatedHandler) Name() string {
	return "domain_UpdatedHandler"
}
