package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/sensor"
	"github.com/peteqproj/peteq/domain/sensor/event/types"
	"github.com/peteqproj/peteq/internal/errors"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// TriggerCommand to create task
	TriggerCommand struct {
		Eventbus bus.EventPublisher
		Repo     *sensor.Repo
	}

	// SensorTriggerCommandOptions options to sensor the sensor
	SensorTriggerCommandOptions struct {
		ID   string      `json:"id"`
		Data interface{} `json:"data"`
	}
)

// Handle runs TriggerCommand to create task
func (m *TriggerCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &SensorTriggerCommandOptions{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to SensorTriggerCommandOptions object")
	}

	u := tenant.UserFromContext(ctx)
	if u == nil {
		return errors.ErrMissingUserInContext
	}
	_, err = m.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.SensorTriggeredEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "sensor",
			AggregatorID:   opt.ID,
		},
		Spec: opt.Data,
	})
	return err
}
