package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/domain/task/event/handler"
	"github.com/peteqproj/peteq/domain/task/event/types"
	"github.com/peteqproj/peteq/internal/errors"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// ReopenCommand to create task
	ReopenCommand struct {
		Eventbus bus.EventPublisher
		Repo     *task.Repo
	}

	ReopenTaskArguments struct {
		TaskID string `json:"taskId"`
	}
)

// Handle runs ReopenCommand to create task
func (r *ReopenCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt := &ReopenTaskArguments{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to ReopenTaskArguments object")
	}
	u := tenant.UserFromContext(ctx)
	if u == nil {
		return errors.ErrMissingUserInContext
	}
	t, err := r.Repo.GetById(ctx, opt.TaskID)
	if err != nil {
		return err
	}
	t.Spec.Completed = false
	if err := r.Repo.UpdateTask(ctx, t); err != nil {
		return err
	}
	_, err = r.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.TaskStatusChanged,
			CreatedAt:      time.Now(),
			AggregatorRoot: "task",
			AggregatorID:   opt.TaskID,
		},
		Spec: handler.StatusChangedSpec{
			Completed: false,
		},
	})
	return err
}
