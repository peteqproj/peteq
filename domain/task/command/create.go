package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
)

type (
	// CreateCommand to create task
	CreateCommand struct {
		Eventbus bus.Eventbus
	}
)

// Handle runs CreateCommand to create task
func (c *CreateCommand) Handle(ctx context.Context, done chan<- error, arguments interface{}) {
	t, ok := arguments.(task.Task)
	if !ok {
		done <- fmt.Errorf("Failed to convert arguments to Task object")
		return
	}
	c.Eventbus.Publish(event.Event{
		Metadata: event.Metadata{
			Name:           "task.created",
			CreatedAt:      time.Now(),
			AggregatorRoot: "task",
			AggregatorID:   t.Metadata.ID,
		},
		Spec: t,
	}, done)
}
