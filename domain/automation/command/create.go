package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/automation/event/handler"
	"github.com/peteqproj/peteq/domain/automation/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// CreateCommand to create task
	CreateCommand struct {
		Eventbus bus.Eventbus
	}

	// AutomationCreateCommandOptions options to create automation
	AutomationCreateCommandOptions struct {
		ID              string
		Name            string
		Description     string
		Type            string
		JSONInputSchema string
	}
)

// Handle runs CreateCommand to create task
func (m *CreateCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt, ok := arguments.(AutomationCreateCommandOptions)
	if !ok {
		return fmt.Errorf("Failed to convert arguments to AutomationCreateCommandOptions object")
	}

	u := tenant.UserFromContext(ctx)
	_, err := m.Eventbus.Publish(ctx, event.Event{
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
