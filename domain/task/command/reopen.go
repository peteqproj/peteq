package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// ReopenCommand to create task
	ReopenCommand struct {
		Eventbus bus.Eventbus
	}
)

// Handle runs ReopenCommand to create task
func (r *ReopenCommand) Handle(ctx context.Context, done chan<- error, arguments interface{}) {
	t, ok := arguments.(string)
	if !ok {
		done <- fmt.Errorf("Failed to convert arguments to string")
		return
	}
	u := tenant.UserFromContext(ctx)
	r.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           "task.reopened",
			CreatedAt:      time.Now(),
			AggregatorRoot: "task",
			AggregatorID:   t,
		},
	}, done)
}
