package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/user/event/handler"
	"github.com/peteqproj/peteq/domain/user/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// RegisterCommand to create task
	RegisterCommand struct {
		Eventbus bus.Eventbus
	}

	// RegisterCommandOptions to create new user
	RegisterCommandOptions struct {
		UserID       string
		Email        string
		PasswordHash string
	}
)

// Handle runs RegisterCommand to create new user
func (r *RegisterCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt, ok := arguments.(RegisterCommandOptions)
	if !ok {
		return fmt.Errorf("Failed to convert arguments to User")
	}
	_, err := r.Eventbus.Publish(ctx, event.Event{
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
