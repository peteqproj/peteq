package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/project/event/handler"
	"github.com/peteqproj/peteq/domain/project/event/types"
	"github.com/peteqproj/peteq/internal/errors"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// AddTaskCommand to create task
	AddTaskCommand struct {
		Eventbus bus.EventPublisher
		Repo     *project.Repo
	}

	// AddTasksCommandOptions options to add tasks to project
	AddTasksCommandOptions struct {
		Project string `json:"project"`
		TaskID  string `json:"taskId"`
	}
)

// Handle runs AddTaskCommand to create task
func (m *AddTaskCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &AddTasksCommandOptions{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to AddTasksCommandOptions object")
	}
	prj, err := m.Repo.GetById(ctx, opt.Project)
	if err != nil {
		return err
	}
	prj.Spec.Tasks = append(prj.Spec.Tasks, opt.TaskID)
	if err := m.Repo.UpdateProject(ctx, prj); err != nil {
		return err
	}
	u := tenant.UserFromContext(ctx)
	if u == nil {
		return errors.ErrMissingUserInContext
	}
	_, err = m.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.TaskAddedToProjectEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "project",
			AggregatorID:   opt.Project,
		},
		Spec: handler.TaskAddedSpec{
			Project: opt.Project,
			TaskID:  opt.TaskID,
		},
	})
	return err
}
