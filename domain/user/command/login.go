package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/user/event/handler"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// LoginCommand to create task
	LoginCommand struct {
		Eventbus bus.Eventbus
	}

	// LoginCommandOptions add new token to allow api calls
	LoginCommandOptions struct {
		UserID      string
		HashedToken string
	}
)

// Handle runs LoginCommand to create new user
func (r *LoginCommand) Handle(ctx context.Context, done chan<- error, arguments interface{}) {
	fmt.Println("user.login command handler")
	opt, ok := arguments.(LoginCommandOptions)
	if !ok {
		done <- fmt.Errorf("Failed to convert arguments to LoginCommandOptions")
		return
	}
	r.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   opt.UserID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           "user.loggedin",
			CreatedAt:      time.Now(),
			AggregatorRoot: "user",
			AggregatorID:   opt.UserID,
		},
		Spec: handler.LoggedinSpec{
			ID:        opt.UserID,
			TokenHash: opt.HashedToken,
		},
	}, done)
}
