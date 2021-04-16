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
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CreateCommand to create task
	CreateCommand struct {
		Eventbus bus.EventPublisher
		Repo     *project.Repo
	}

	// CreateProjectCommandOptions to create new project
	CreateProjectCommandOptions struct {
		ID          string            `json:"id"`
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Labels      map[string]string `json:"labels"`
		Color       string            `json:"color"`
		ImageURL    string            `json:"imageUrl"`
	}
)

// Handle runs CreateCommand to create task
func (m *CreateCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &CreateProjectCommandOptions{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to Project object")
	}
	prj := &project.Project{
		Metadata: project.Metadata{
			ID:          opt.ID,
			Name:        opt.Name,
			Labels:      opt.Labels,
			Description: utils.PtrString(opt.Description),
		},
		Spec: project.Spec{
			Color:    utils.PtrString(opt.Color),
			ImageURL: utils.PtrString(opt.ImageURL),
			Tasks:    []string{},
		},
	}
	if err := m.Repo.Create(ctx, prj); err != nil {
		return err
	}
	u := tenant.UserFromContext(ctx)
	_, err = m.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.ProjectCreatedEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "project",
			AggregatorID:   opt.ID,
		},
		Spec: handler.CreatedSpec{
			ID:          opt.ID,
			Name:        opt.Name,
			Description: opt.Description,
			Color:       opt.Color,
			ImageURL:    opt.ImageURL,
		},
	})
	return err
}
