package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/project/event/handler"
	"github.com/peteqproj/peteq/domain/project/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// AddTaskCommand to create task
	AddTaskCommand struct {
		Eventbus bus.Eventbus
	}

	// AddTasksCommandOptions options to add tasks to project
	AddTasksCommandOptions struct {
		Project string `json:"project"`
		TaskID  string `json:"task"`
	}
)

// Handle runs AddTaskCommand to create task
func (m *AddTaskCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt, ok := arguments.(AddTasksCommandOptions)
	if !ok {
		return fmt.Errorf("Failed to convert arguments to AddTasksCommandOptions object")
	}
	u := tenant.UserFromContext(ctx)
	_, err := m.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.TaskAddedToProjectEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "project",
			AggregatorID:   opt.Project,
		},
		Spec: handler.TaskAddedSpec{
			Project: opt.Project,
			TaskID:  opt.TaskID,
		},
	})
	return err
}
