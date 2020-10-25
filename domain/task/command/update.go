package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/domain/task/event/handler"
	"github.com/peteqproj/peteq/domain/task/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// UpdateCommand to create task
	UpdateCommand struct {
		Eventbus bus.Eventbus
	}
)

// Handle runs UpdateCommand to create task
func (u *UpdateCommand) Handle(ctx context.Context, done chan<- error, arguments interface{}) {
	t, ok := arguments.(task.Task)
	if !ok {
		done <- fmt.Errorf("Failed to convert arguments to Task object")
	}
	user := tenant.UserFromContext(ctx)
	u.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   user.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.TaskUpdatedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "task",
			AggregatorID:   t.Metadata.ID,
		},
		Spec: handler.UpdatedSpec{
			ID:          t.Metadata.ID,
			Name:        t.Metadata.Name,
			Description: t.Metadata.Description,
		},
	}, done)
}
