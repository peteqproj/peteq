package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/list/event/handler"
	"github.com/peteqproj/peteq/domain/list/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// Create to create task
	Create struct {
		Eventbus bus.EventPublisher
	}

	// CreateCommandOptions is the arguments the command expects
	CreateCommandOptions struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		Index int    `json:"index"`
	}
)

// Handle runs Create to create task
func (m *Create) Handle(ctx context.Context, arguments interface{}) error {
	opt := &CreateCommandOptions{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to CreateCommandOptions object")
	}

	u := tenant.UserFromContext(ctx)
	_, err = m.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.ListCreatedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "list",
			AggregatorID:   opt.ID,
		},
		Spec: handler.CreatedSpec{
			ID:    opt.ID,
			Name:  opt.Name,
			Index: opt.Index,
		},
	})
	return err

}
