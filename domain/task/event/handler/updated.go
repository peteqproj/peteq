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
	// UpdatedSpec is the event.spec for this event
	UpdatedSpec struct {
		ID          string `json:"id" yaml:"id"`
		Name        string `json:"name" yaml:"name"`
		Description string `json:"description" yaml:"description"`
	}
)

// Handle will handle the event the process it
func (c *UpdatedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := UpdatedSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return c.Repo.Update(ev.Tenant.ID, task.Task{
		Metadata: task.Metadata{
			ID:          opt.ID,
			Name:        opt.Name,
			Description: opt.Description,
		},
	})
}

func (c *UpdatedHandler) Name() string {
	return "domain_task"
}
