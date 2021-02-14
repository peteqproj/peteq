package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/list/event/handler"
	"github.com/peteqproj/peteq/domain/list/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// MoveTaskCommand to create task
	MoveTaskCommand struct {
		Eventbus bus.EventPublisher
	}

	// MoveTaskArguments is the arguments the command expects
	MoveTaskArguments struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		TaskID      string `json:"taskId"`
	}
)

// Handle runs MoveTaskCommand to create task
func (m *MoveTaskCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &MoveTaskArguments{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to MoveTaskArguments object")
	}

	u := tenant.UserFromContext(ctx)
	_, err = m.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.TaskMovedIntoListEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "list",
			AggregatorID:   opt.Source,
		},
		Spec: handler.TaskMovedSpec{
			TaskID:      opt.TaskID,
			Source:      opt.Source,
			Destination: opt.Destination,
		},
	})
	return err
}
