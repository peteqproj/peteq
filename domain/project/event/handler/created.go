package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/project"
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

	return t.Repo.Create(ctx, repo.Resource{
		Metadata: repo.Metadata{
			Type: "project",
			Name: opt.Name,
			ID:   opt.ID,
		},
		Spec: project.Spec{
			Color:    opt.Color,
			ImageURL: opt.ImageURL,
			Tasks:    []string{},
		},
	})
}

func (t *CreatedHandler) Name() string {
	return "domain_project"
}
