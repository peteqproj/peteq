package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/automation"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// CreatedHandler to handle task.created event
	CreatedHandler struct {
		Repo *automation.Repo
	}

	// CreatedSpec is the event.spec for this event
	CreatedSpec struct {
		ID              string `json:"id" yaml:"id"`
		Name            string `json:"name" yaml:"name"`
		Description     string `json:"description" yaml:"description"`
		Type            string `json:"type" yaml:"yaml"`
		JSONInputSchema string `json:"jsonInputSchema" yaml:"jsonInputSchema"`
	}
)

// Handle will process it the event
func (t *CreatedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := CreatedSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	spec := automation.AutomationSpec{
		Type:            opt.Type,
		JSONInputSchema: opt.JSONInputSchema,
	}

	return t.Repo.Create(ev.Tenant.ID, automation.Automation{
		Metadata: automation.Metadata{
			ID:          opt.ID,
			Name:        opt.Name,
			Description: &opt.Description,
		},
		Spec: spec,
	})
}

func (t *CreatedHandler) Name() string {
	return "domain_automation"
}
