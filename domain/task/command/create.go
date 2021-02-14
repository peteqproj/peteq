package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/task/event/handler"
	"github.com/peteqproj/peteq/domain/task/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CreateCommand to create task
	CreateCommand struct {
		Eventbus bus.EventPublisher
	}

	// CreateCommandOptions add new token to allow api calls
	CreateCommandOptions struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

// Handle runs CreateCommand to create task
func (c *CreateCommand) Handle(ctx context.Context, arguments interface{}) error {
	t := &CreateCommandOptions{}
	if err := utils.UnmarshalInto(arguments, t); err != nil {
		return fmt.Errorf("Failed to convert arguments to CreateCommandOptions object")
	}
	u := tenant.UserFromContext(ctx)
	tenant := tenant.Tenant{
		ID:   u.Metadata.ID,
		Type: tenant.User.String(),
	}
	_, err := c.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant,
		Metadata: event.Metadata{
			Name:           types.TaskCreatedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "task",
			AggregatorID:   t.ID,
		},
		Spec: handler.CreatedSpec{
			ID:   t.ID,
			Name: t.Name,
		},
	})
	return err
}
