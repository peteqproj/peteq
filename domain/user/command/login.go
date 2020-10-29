package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/domain/user/event/handler"
	"github.com/peteqproj/peteq/domain/user/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// LoginCommand to create task
	LoginCommand struct {
		Eventbus bus.Eventbus
		Repo     *user.Repo
	}

	// LoginCommandOptions add new token to allow api calls
	LoginCommandOptions struct {
		HashedToken    string
		Email          string
		HashedPassword string
	}
)

// Handle runs LoginCommand to create new user
func (r *LoginCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt, ok := arguments.(LoginCommandOptions)
	if !ok {
		return fmt.Errorf("Failed to convert arguments to LoginCommandOptions")
	}
	user, err := r.Repo.GetByEmail(opt.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if opt.HashedPassword != user.Spec.PasswordHash {
		return fmt.Errorf("Invalid credentials")
	}

	_, err = r.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   user.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.UserLoggedIn,
			CreatedAt:      time.Now(),
			AggregatorRoot: "user",
			AggregatorID:   user.Metadata.ID,
		},
		Spec: handler.LoggedinSpec{
			ID:        user.Metadata.ID,
			TokenHash: opt.HashedToken,
		},
	})
	return err
}
