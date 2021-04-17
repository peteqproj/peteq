package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/trigger"
	"github.com/peteqproj/peteq/domain/trigger/event/handler"
	"github.com/peteqproj/peteq/domain/trigger/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CreateCommand to create task
	CreateCommand struct {
		Eventbus bus.EventPublisher
		Repo     *trigger.Repo
	}

	// TriggerCreateCommandOptions options to create trigger
	TriggerCreateCommandOptions struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Cron        *string `json:"cron"`
		URL         *string `json:"url"`
	}
)

// Handle runs CreateCommand to create task
func (m *CreateCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &TriggerCreateCommandOptions{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to TriggerCreateCommandOptions object")
	}

	u := tenant.UserFromContext(ctx)
	if err := m.Repo.Create(ctx, &trigger.Trigger{
		Metadata: trigger.Metadata{
			ID:          opt.ID,
			Name:        opt.Name,
			Labels:      map[string]string{},
			Description: utils.PtrString(""),
		},
		Spec: trigger.Spec{
			Cron: opt.Cron,
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
			Name:           types.TriggerCreatedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "trigger",
			AggregatorID:   opt.ID,
		},
		Spec: handler.CreatedSpec{
			ID:          opt.ID,
			Name:        opt.Name,
			Description: opt.Description,
			Cron:        opt.Cron,
		},
	})
	return err
}
