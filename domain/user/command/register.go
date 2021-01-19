package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/domain/user/event/handler"
	"github.com/peteqproj/peteq/domain/user/event/types"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// RegisterCommand to create task
	RegisterCommand struct {
		Eventbus    bus.Eventbus
		Repo        *user.Repo
		Commandbus  commandbus.CommandBus
		IDGenerator utils.IDGenerator
	}

	// RegisterCommandOptions to create new user
	RegisterCommandOptions struct {
		UserID       string `json:"userId"`
		Email        string `json:"email"`
		PasswordHash string `json:"passwordHash"`
	}
)

// Handle runs RegisterCommand to create new user
func (r *RegisterCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &RegisterCommandOptions{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to User")
	}
	usr, err := r.Repo.GetByEmail(opt.Email)
	if err != nil {
		if err.Error() != "User not found" {
			return err
		}
	}
	if usr != nil {
		return fmt.Errorf("Email already registred")
	}
	ectx := tenant.ContextWithUser(ctx, user.User{
		Metadata: user.Metadata{
			ID:    opt.UserID,
			Email: opt.Email,
		},
	})
	_, err = r.Eventbus.Publish(ectx, event.Event{
		Tenant: tenant.Tenant{
			ID:   opt.UserID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.UserRegistredEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "user",
			AggregatorID:   opt.UserID,
		},
		Spec: handler.RegisteredSpec{
			Email:        opt.Email,
			ID:           opt.UserID,
			PasswordHash: opt.PasswordHash,
		},
	})
	return err
}
