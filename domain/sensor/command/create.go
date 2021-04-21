package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/sensor"
	"github.com/peteqproj/peteq/domain/sensor/event/handler"
	"github.com/peteqproj/peteq/domain/sensor/event/types"
	"github.com/peteqproj/peteq/internal/errors"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CreateCommand to create task
	CreateCommand struct {
		Eventbus bus.EventPublisher
		Repo     *sensor.Repo
	}

	// SensorCreateCommandOptions options to create sensor
	SensorCreateCommandOptions struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Cron        *string `json:"cron"`
		URL         *string `json:"url"`
	}
)

// Handle runs CreateCommand to create task
func (m *CreateCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &SensorCreateCommandOptions{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to SensorCreateCommandOptions object")
	}

	u := tenant.UserFromContext(ctx)
	if u == nil {
		return errors.ErrMissingUserInContext
	}
	if err := m.Repo.Create(ctx, &sensor.Sensor{
		Metadata: sensor.Metadata{
			ID:          opt.ID,
			Name:        opt.Name,
			Labels:      map[string]string{},
			Description: utils.PtrString(""),
		},
		Spec: sensor.Spec{
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
			Name:           types.SensorCreatedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "sensor",
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
