package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
)

type (
	// MoveTaskCommand to create task
	MoveTaskCommand struct {
		Eventbus bus.Eventbus
	}

	// MoveTaskArguments is the arguments the command expects
	MoveTaskArguments struct {
		Source      string
		Destination string
		TaskID      string
	}
)

// Handle runs MoveTaskCommand to create task
func (m *MoveTaskCommand) Handle(ctx context.Context, done chan<- error, arguments interface{}) {
	opt, ok := arguments.(MoveTaskArguments)
	if !ok {
		done <- fmt.Errorf("Failed to convert arguments to MoveTaskArguments object")
		return
	}

	m.Eventbus.Publish(event.Event{
		Metadata: event.Metadata{
			Name:           "list.task-moved",
			CreatedAt:      time.Now(),
			AggregatorRoot: "list",
			AggregatorID:   opt.Source,
		},
		Spec: opt,
	}, done)
}
