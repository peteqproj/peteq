package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// UpdatedHandler to handle task.created event
	UpdatedHandler struct {
		Repo *task.Repo
	}
	// UpdatedSpec is the event.spec for this event
	UpdatedSpec struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

// Handle will handle the event the process it
func (c *UpdatedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := UpdatedSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	t, err := c.Repo.GetById(ctx, opt.ID)
	if err != nil {
		return err
	}
	t.Metadata.Name = opt.Name
	t.Metadata.Description = utils.PtrString(opt.Description)
	return c.Repo.UpdateTask(ctx, t)
}

func (c *UpdatedHandler) Name() string {
	return "domain_task"
}
