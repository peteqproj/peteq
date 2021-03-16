package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/repo"
)

type (
	// DeleteHandler to handle task.deleted event
	DeleteHandler struct {
		Repo *repo.Repo
	}
	// DeletedSpec is the event.spec for this event
	DeletedSpec struct {
		ID string `json:"id"`
	}
)

// Handle will handle the event the process it
func (c *DeleteHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := DeletedSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return c.Repo.Delete(ctx, repo.Resource{
		Metadata: repo.Metadata{
			Type: "task",
			ID:   ev.ID,
		},
		Spec: task.Spec{},
	})
}

func (c *DeleteHandler) Name() string {
	return "domain_task"
}
