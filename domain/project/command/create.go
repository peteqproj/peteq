package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/project"
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
func (m *CreateCommand) Handle(ctx context.Context, done chan<- error, arguments interface{}) {
	opt, ok := arguments.(project.Project)
	if !ok {
		done <- fmt.Errorf("Failed to convert arguments to Project object")
		return
	}

	m.Eventbus.Publish(event.Event{
		Metadata: event.Metadata{
			Name:           "project.created",
			CreatedAt:      time.Now(),
			AggregatorRoot: "project",
			AggregatorID:   opt.Metadata.ID,
		},
		Spec: opt,
	}, done)
}
