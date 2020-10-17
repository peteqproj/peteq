package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// CreatedHandler to handle task.created event
	CreatedHandler struct {
		Repo *project.Repo
	}
)

// Handle will process it the event
func (t *CreatedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := project.Project{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return t.Repo.Create(ev.Tenant.ID, opt)
}

func (t *CreatedHandler) Name() string {
	return "project_domain_CreatedHandler"
}
