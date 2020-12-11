package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/task/event/handler"
	"github.com/peteqproj/peteq/domain/task/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CompleteCommand to create task
	CompleteCommand struct {
		Eventbus bus.Eventbus
	}

	CompleteTaskArguments struct {
		TaskID string `json:"taskId"`
	}
)

// Handle runs CompleteCommand to create task
func (c *CompleteCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &CompleteTaskArguments{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to CompleteTaskArguments object")
	}
	u := tenant.UserFromContext(ctx)
	_, err = c.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.TaskStatusChanged,
			CreatedAt:      time.Now(),
			AggregatorRoot: "task",
			AggregatorID:   opt.TaskID,
		},
		Spec: handler.StatusChangedSpec{
			Completed: true,
		},
	})
	return err
}
