package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/trigger"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// CreatedHandler to handle task.created event
	CreatedHandler struct {
		Repo *trigger.Repo
	}

	// CreatedSpec is the event.spec for this event
	CreatedSpec struct {
		ID              string            `json:"id" yaml:"id"`
		Name            string            `json:"name" yaml:"name"`
		Description     string            `json:"description" yaml:"description"`
		Cron            *string           `json:"cron,omitempty" yaml:"cron,omitempty"`
		URL             *string           `json:"url,omitempty" yaml:"url,omitempty"`
		RequiredHeaders map[string]string `json:"requiredHeaders" yaml:"requiredHeaders"`
	}
)

// Handle will process it the event
func (t *CreatedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := CreatedSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return t.Repo.Create(ev.Tenant.ID, trigger.Trigger{
		Metadata: trigger.Metadata{
			ID:          opt.ID,
			Name:        opt.Name,
			Description: opt.Description,
		},
	})
}

func (t *CreatedHandler) Name() string {
	return "domain_trigger"
}
