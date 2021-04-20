package command

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/list/event/handler"
	"github.com/peteqproj/peteq/domain/list/event/types"
	"github.com/peteqproj/peteq/internal/errors"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// MoveTaskCommand to create task
	MoveTaskCommand struct {
		Eventbus bus.EventPublisher
		Repo     *list.Repo
	}

	// MoveTaskArguments is the arguments the command expects
	MoveTaskArguments struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		TaskID      string `json:"taskId"`
	}
)

// Handle runs MoveTaskCommand to create task
func (m *MoveTaskCommand) Handle(ctx context.Context, arguments interface{}) error {
	usr := tenant.UserFromContext(ctx)
	if usr == nil {
		return errors.ErrMissingUserInContext
	}
	opt := &MoveTaskArguments{}
	err := utils.UnmarshalInto(arguments, opt)
	if err != nil {
		return fmt.Errorf("Failed to convert arguments to MoveTaskArguments object")
	}

	if opt.Source != "" {
		list, err := m.Repo.GetById(ctx, opt.Source)
		if err != nil {
			return fmt.Errorf("Failed to get source list: %v", err)
		}
		index := -1
		for i, t := range list.Spec.Tasks {
			if opt.TaskID == t {
				index = i
				break
			}
		}
		if index == -1 {
			list.Spec.Tasks = append(list.Spec.Tasks[:index], list.Spec.Tasks[index+1:]...)
			if err := m.Repo.UpdateList(ctx, list); err != nil {
				return fmt.Errorf("Failed to remove task from source list: %v", err)
			}
		}
	}

	if opt.Destination != "" {
		list, err := m.Repo.GetById(ctx, opt.Destination)
		if err != nil {
			return fmt.Errorf("Failed to get destination list: %v", err)
		}
		list.Spec.Tasks = append(list.Spec.Tasks, opt.TaskID)
		if err := m.Repo.UpdateList(ctx, list); err != nil {
			return fmt.Errorf("Failed to add task to destination list %v", err)
		}
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
			Name:           types.TaskMovedIntoListEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "list",
			AggregatorID:   opt.Source,
		},
		Spec: handler.TaskMovedSpec{
			TaskID:      opt.TaskID,
			Source:      opt.Source,
			Destination: opt.Destination,
		},
	})
	return err
}
