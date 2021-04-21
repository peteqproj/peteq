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
	// CreateSensorBindingCommand to create task
	CreateSensorBindingCommand struct {
		Eventbus bus.EventPublisher
		Repo     *automation.Repo
	}

	// SensorBindingCreateCommandOptions options to create automation
	SensorBindingCreateCommandOptions struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Sensor     string `json:"sensor"`
		Automation string `json:"automation"`
	}
)

// Handle runs CreateSensorBindingCommand to create task
func (m *CreateSensorBindingCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &SensorBindingCreateCommandOptions{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to SensorBindingCreateCommandOptions object")
	}

	u := tenant.UserFromContext(ctx)
	if u == nil {
		return errors.ErrMissingUserInContext
	}
	if err := m.Repo.CreateSensorBinding(ctx, &automation.SensorBinding{
		Metadata: automation.Metadata{
			ID:          opt.ID,
			Name:        opt.Name,
			Labels:      map[string]string{},
			Description: utils.PtrString(""),
		},
		Spec: automation.SensorBindingSpec{
			Automation: opt.Automation,
			Sensor:     opt.Sensor,
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
			Name:           types.SensorBindingCreatedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "automation",
			AggregatorID:   opt.ID,
		},
		Spec: handler.SensorBindingCreatedSpec{
			ID:         opt.ID,
			Name:       opt.Name,
			Automation: opt.Automation,
			Sensor:     opt.Sensor,
		},
	})
	return err
}
