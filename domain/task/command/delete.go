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
	// DeleteCommand to create task
	DeleteCommand struct {
		Eventbus bus.Eventbus
	}
)

// Handle runs DeleteCommand to create task
func (c *DeleteCommand) Handle(ctx context.Context, done chan<- error, arguments interface{}) {
	t, ok := arguments.(task.Task)
	if !ok {
		done <- fmt.Errorf("Failed to convert arguments to Task object")
	}
	c.Eventbus.Publish(event.Event{
		Metadata: event.Metadata{
			Name:           "task.deleted",
			CreatedAt:      time.Now(),
			AggregatorRoot: "task",
			AggregatorID:   t.Metadata.ID,
		},
		Spec: t,
	}, done)
}
