package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/project/event/handler"
	"github.com/peteqproj/peteq/domain/project/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
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

	u := tenant.UserFromContext(ctx)
	m.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.ProjectCreatedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "project",
			AggregatorID:   opt.Metadata.ID,
		},
		Spec: handler.CreatedSpec{
			ID:          opt.Metadata.ID,
			Name:        opt.Metadata.Name,
			Description: opt.Metadata.Description,
			Color:       opt.Metadata.Color,
			ImageURL:    opt.Metadata.ImageURL,
		},
	}, done)
}
