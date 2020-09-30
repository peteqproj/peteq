package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// RegisterCommand to create task
	RegisterCommand struct {
		Eventbus bus.Eventbus
	}
)

// Handle runs RegisterCommand to create new user
func (r *RegisterCommand) Handle(ctx context.Context, done chan<- error, arguments interface{}) {
	fmt.Println("user.register command handler")
	u, ok := arguments.(user.User)
	if !ok {
		done <- fmt.Errorf("Failed to convert arguments to User")
		return
	}
	r.Eventbus.Publish(event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           "user.registred",
			CreatedAt:      time.Now(),
			AggregatorRoot: "user",
			AggregatorID:   u.Metadata.ID,
		},
		Spec: u,
	}, done)
}
