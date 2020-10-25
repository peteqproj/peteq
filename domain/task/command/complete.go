package command

import (
	"context"
	"fmt"
	"time"

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
func (c *CompleteCommand) Handle(ctx context.Context, done chan<- error, arguments interface{}) {
	t, ok := arguments.(string)
	if !ok {
		done <- fmt.Errorf("Failed to convert arguments to string")
		return
	}
	u := tenant.UserFromContext(ctx)
	c.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.TaskCompletedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "task",
			AggregatorID:   t,
		},
	}, done)
}
