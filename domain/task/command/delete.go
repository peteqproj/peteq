package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/domain/task/event/handler"
	"github.com/peteqproj/peteq/domain/task/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// DeleteCommand to create task
	DeleteCommand struct {
		Eventbus bus.EventPublisher
	}
)

// Handle runs DeleteCommand to create task
func (c *DeleteCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &task.Task{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to Task object")
	}
	u := tenant.UserFromContext(ctx)
	_, err = c.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.TaskDeletedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "task",
			AggregatorID:   opt.Metadata.ID,
		},
		Spec: handler.DeletedSpec{
			ID: opt.Metadata.ID,
		},
	})
	return err
}
