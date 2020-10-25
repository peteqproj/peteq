package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/project/event/handler"
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
func (m *AddTaskCommand) Handle(ctx context.Context, done chan<- error, arguments interface{}) {
	opt, ok := arguments.(AddTasksCommandOptions)
	if !ok {
		done <- fmt.Errorf("Failed to convert arguments to AddTasksCommandOptions object")
		return
	}
	u := tenant.UserFromContext(ctx)
	m.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           "project.task-added",
			CreatedAt:      time.Now(),
			AggregatorRoot: "project",
			AggregatorID:   opt.Project,
		},
		Spec: handler.TaskAddedSpec{
			Project: opt.Project,
			TaskID:  opt.TaskID,
		},
	}, done)
}
