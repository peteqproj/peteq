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
)

type (
	// CompleteCommand to create task
	CompleteCommand struct {
		Eventbus bus.Eventbus
	}
)

// Handle runs CompleteCommand to create task
func (c *CompleteCommand) Handle(ctx context.Context, arguments interface{}) error {
	t, ok := arguments.(string)
	if !ok {
		return fmt.Errorf("Failed to convert arguments to string")
	}
	u := tenant.UserFromContext(ctx)
	_, err := c.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.TaskStatusChanged,
			CreatedAt:      time.Now(),
			AggregatorRoot: "task",
			AggregatorID:   t,
		},
		Spec: handler.StatusChangedSpec{
			Completed: true,
		},
	})
	return err
}
