package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/repo"
)

type (
	// CreatedHandler to handle task.created event
	CreatedHandler struct {
		Repo *repo.Repo
	}

	// CreatedSpec is the event.spec for this event
	CreatedSpec struct {
		Name        string            `json:"name"`
		ID          string            `json:"id"`
		Description string            `json:"description"`
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
	t := task.NewTask(opt.ID, opt.Name, opt.Description)
	t.Metadata.Labels = opt.Labels
	t.Spec = task.Spec{Completed: false}
	return c.Repo.Create(ctx, t)
}

func (c *CreatedHandler) Name() string {
	return "domain_task"
}
