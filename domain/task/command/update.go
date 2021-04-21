package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/domain/task/event/handler"
	"github.com/peteqproj/peteq/domain/task/event/types"
	"github.com/peteqproj/peteq/internal/errors"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// UpdateCommand to create task
	UpdateCommand struct {
		Eventbus bus.EventPublisher
		Repo     *task.Repo
	}
)

// Handle runs UpdateCommand to create task
func (u *UpdateCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &task.Task{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to Task object")
	}
	if err := u.Repo.UpdateTask(ctx, opt); err != nil {
		return err
	}
	user := tenant.UserFromContext(ctx)
	if u == nil {
		return errors.ErrMissingUserInContext
	}
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
			ID:   opt.Metadata.ID,
			Name: opt.Metadata.Name,
		},
	})
	return err
}
