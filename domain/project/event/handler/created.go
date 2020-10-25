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

	// CreatedSpec is the event.spec for this event
	CreatedSpec struct {
		ID          string `json:"id" yaml:"id"`
		Name        string `json:"name" yaml:"name"`
		Description string `json:"description" yaml:"description"`
		Color       string `json:"color" yaml:"color"`
		ImageURL    string `json:"imageUrl" yaml:"imageUrl"`
	}
)

// Handle will process it the event
func (t *CreatedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := CreatedSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return t.Repo.Create(ev.Tenant.ID, project.Project{
		Metadata: project.Metadata{
			ID:          opt.ID,
			Description: opt.Description,
			Name:        opt.Name,
			Color:       opt.Color,
			ImageURL:    opt.ImageURL,
		},
	})
}

func (t *CreatedHandler) Name() string {
	return "project_domain_CreatedHandler"
}
