package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
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
	fmt.Println("publishing event")
	m.Eventbus.Publish(event.Event{
		Metadata: event.Metadata{
			Name:           "project.task-added",
			CreatedAt:      time.Now(),
			AggregatorRoot: "project",
			AggregatorID:   opt.Project,
		},
		Spec: opt,
	}, done)
}
