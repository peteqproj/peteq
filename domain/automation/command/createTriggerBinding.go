package command

import (
	"context"
	"fmt"
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
	// CreateTriggerBindingCommand to create task
	CreateTriggerBindingCommand struct {
		Eventbus bus.EventPublisher
		Repo     *automation.Repo
	}

	// TriggerBindingCreateCommandOptions options to create automation
	TriggerBindingCreateCommandOptions struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Trigger    string `json:"trigger"`
		Automation string `json:"automation"`
	}
)

// Handle runs CreateTriggerBindingCommand to create task
func (m *CreateTriggerBindingCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &TriggerBindingCreateCommandOptions{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to TriggerBindingCreateCommandOptions object")
	}

	u := tenant.UserFromContext(ctx)
	if u == nil {
		return errors.ErrMissingUserInContext
	}
	if err := m.Repo.CreateTriggerBinding(ctx, &automation.TriggerBinding{
		Metadata: automation.Metadata{
			ID:          opt.ID,
			Name:        opt.Name,
			Labels:      map[string]string{},
			Description: utils.PtrString(""),
		},
		Spec: automation.TriggerBindingSpec{
			Automation: opt.Automation,
			Trigger:    opt.Trigger,
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
			Name:           types.TriggerBindingCreatedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "automation",
			AggregatorID:   opt.ID,
		},
		Spec: handler.TriggerBindingCreatedSpec{
			ID:         opt.ID,
			Name:       opt.Name,
			Automation: opt.Automation,
			Trigger:    opt.Trigger,
		},
	})
	return err
}
