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
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// LoginCommand to create task
	LoginCommand struct {
		Eventbus bus.EventPublisher
		Repo     *user.Repo
	}

	// LoginCommandOptions add new token to allow api calls
	LoginCommandOptions struct {
		HashedToken    string `json:"hashedToken"`
		Email          string `json:"email"`
		HashedPassword string `json:"hashedPassword"`
	}
)

// Handle runs LoginCommand to create new user
func (r *LoginCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &LoginCommandOptions{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to LoginCommandOptions")
	}
	user, err := r.Repo.GetByEmail(ctx, opt.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if opt.HashedPassword != user.Spec.PasswordHash {
		return fmt.Errorf("Invalid credentials")
	}

	user.Spec.TokenHash = opt.HashedToken
	if err := r.Repo.UpdateUser(ctx, user); err != nil {
		return err
	}
	ectx := tenant.ContextWithUser(ctx, *user)
	_, err = r.Eventbus.Publish(ectx, event.Event{
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
			ID: user.Metadata.ID,
		},
	})
	return err
}
