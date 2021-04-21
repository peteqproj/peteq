package command

import (
	"context"
	"time"

	"github.com/peteqproj/peteq/domain/automation"
	"github.com/peteqproj/peteq/domain/automation/event/handler"
	"github.com/peteqproj/peteq/domain/automation/event/types"
	errors "github.com/peteqproj/peteq/internal/errors"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CreateCommand to create task
	CreateCommand struct {
		Eventbus bus.EventPublisher
		Repo     *automation.Repo
	}

	// AutomationCreateCommandOptions options to create automation
	AutomationCreateCommandOptions struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Description     string `json:"description"`
		Type            string `json:"type"`
		JSONInputSchema string `json:"jsonInputSchema"`
	}
)

// Handle runs CreateCommand to create task
func (m *CreateCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &AutomationCreateCommandOptions{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return errors.ErrMissingUserInContext
	}

	u := tenant.UserFromContext(ctx)
	if u == nil {
		return errors.ErrMissingUserInContext
	}
	if err := m.Repo.Create(ctx, &automation.Automation{
		Metadata: automation.Metadata{
			ID:          opt.ID,
			Name:        opt.Name,
			Labels:      map[string]string{},
			Description: utils.PtrString(""),
		},
		Spec: automation.AutomationSpec{
			JSONInputSchema: opt.JSONInputSchema,
			Type:            opt.Type,
		},
	}); err != nil {
		return err
	}
	_, err = m.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.AutomationCreatedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "automation",
			AggregatorID:   opt.ID,
		},
		Spec: handler.CreatedSpec{
			ID:              opt.ID,
			Name:            opt.Name,
			Description:     opt.Description,
			Type:            opt.Type,
			JSONInputSchema: opt.JSONInputSchema,
		},
	})
	return err
}
