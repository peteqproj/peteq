package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/domain/task/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// DeleteCommand to create task
	DeleteCommand struct {
		Eventbus bus.Eventbus
	}
)

// Handle runs DeleteCommand to create task
func (c *DeleteCommand) Handle(ctx context.Context, arguments interface{}) error {
	t, ok := arguments.(task.Task)
	if !ok {
		return fmt.Errorf("Failed to convert arguments to Task object")
	}
	u := tenant.UserFromContext(ctx)
	_, err := c.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.TaskDeletedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "task",
			AggregatorID:   t.Metadata.ID,
		},
		Spec: t,
	})
	return err
}
