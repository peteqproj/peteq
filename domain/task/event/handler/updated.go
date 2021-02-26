package handler

import (
	"context"

	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/repo"
)

type (
	// UpdatedHandler to handle task.created event
	UpdatedHandler struct {
		Repo *repo.Repo
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
	t, err := c.Repo.Get(ctx, repo.GetOptions{
		ID: opt.ID,
	})
	if err != nil {
		return err
	}
	t.Metadata.Name = opt.Name
	t.Metadata.Description = opt.Description
	return c.Repo.Update(ctx, *t)
}

func (c *UpdatedHandler) Name() string {
	return "domain_task"
}
