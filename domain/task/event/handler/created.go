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

	// CreatedSpec is the event.spec for this event
	CreatedSpec struct {
		Name        string            `json:"name" yaml:"name"`
		ID          string            `json:"id" yaml:"id"`
		Description string            `json:"description" yaml:"description"`
		Labels      map[string]string `json:"labels"`
	}
)

// Handle will handle the event the process it
func (c *CreatedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := CreatedSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return c.Repo.Create(ev.Tenant.ID, task.Task{
		Metadata: task.Metadata{
			ID:          opt.ID,
			Name:        opt.Name,
			Description: opt.Description,
			Labels:      opt.Labels,
		},
		Spec: task.Spec{},
		Status: task.Status{
			Completed: false,
		},
	})
}

func (c *CreatedHandler) Name() string {
	return "domain_task"
}
