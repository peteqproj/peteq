package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// Create to create task
	Create struct {
		Eventbus bus.Eventbus
	}

	// CreateCommandOptions is the arguments the command expects
	CreateCommandOptions struct {
		Name  string
		ID    string
		Index int
	}
)

// Handle runs Create to create task
func (m *Create) Handle(ctx context.Context, done chan<- error, arguments interface{}) {
	opt, ok := arguments.(CreateCommandOptions)
	if !ok {
		done <- fmt.Errorf("Failed to convert arguments to CreateCommandOptions object")
		return
	}

	u := tenant.UserFromContext(ctx)
	m.Eventbus.Publish(event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           "list.created",
			CreatedAt:      time.Now(),
			AggregatorRoot: "list",
			AggregatorID:   opt.ID,
		},
		Spec: opt,
	}, done)
}
