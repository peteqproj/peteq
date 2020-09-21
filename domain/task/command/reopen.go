package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
)

type (
	// ReopenCommand to create task
	ReopenCommand struct {
		Eventbus bus.Eventbus
	}
)

// Handle runs ReopenCommand to create task
func (c *ReopenCommand) Handle(ctx context.Context, done chan<- error, arguments interface{}) {
	t, ok := arguments.(string)
	if !ok {
		done <- fmt.Errorf("Failed to convert arguments to string")
		return
	}
	c.Eventbus.Publish(event.Event{
		Metadata: event.Metadata{
			Name:           "task.reopened",
			CreatedAt:      time.Now(),
			AggregatorRoot: "task",
			AggregatorID:   t,
		},
	}, done)
}
