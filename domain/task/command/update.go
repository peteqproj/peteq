package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/task/event/handler"
	"github.com/peteqproj/peteq/domain/task/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/repo"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// UpdateCommand to create task
	UpdateCommand struct {
		Eventbus bus.EventPublisher
	}
)

// Handle runs UpdateCommand to create task
func (u *UpdateCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &repo.Resource{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to Task object")
	}
	user := tenant.UserFromContext(ctx)
	_, err = u.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   user.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.TaskUpdatedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "task",
			AggregatorID:   opt.Metadata.ID,
		},
		Spec: handler.UpdatedSpec{
			ID:          opt.Metadata.ID,
			Name:        opt.Metadata.Name,
			Description: opt.Metadata.Description,
		},
	})
	return err
}
