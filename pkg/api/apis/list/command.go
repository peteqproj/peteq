package list

import (
	"context"
	"fmt"
	"io"

	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/list/command"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CommandAPI for lists
	CommandAPI struct {
		Repo        *list.Repo
		Commandbus  commandbus.CommandBus
		Logger      logger.Logger
		IDGenerator utils.IDGenerator
	}

	// MoveTasksRequestBody passed from http client
	MoveTasksRequestBody struct {
		Source      string   `json:"source" validate:"required"`
		Destination string   `json:"destination" validate:"required"`
		TaskIDs     []string `json:"tasks" validate:"required"`
	}

	// AddTaskRequestBody passed from http client
	AddTaskRequestBody struct {
		Destination string `json:"destination"`
		TaskID      string `json:"task"`
	}
)

// MoveTasks move multiple tasks from one list to another
func (ca *CommandAPI) MoveTasks(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	u := tenant.UserFromContext(ctx)
	opt := &MoveTasksRequestBody{}
	if err := api.UnmarshalInto(body, opt); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	var source *list.List
	var destination *list.List
	if opt.Source != "" {
		src, err := ca.Repo.Get(u.Metadata.ID, opt.Source)
		if err != nil {
			return api.NewRejectedCommandResponse(fmt.Errorf("Source list: %v", err))
		}
		source = &src
	}
	if opt.Destination != "" {
		dst, err := ca.Repo.Get(u.Metadata.ID, opt.Destination)
		if err != nil {
			return api.NewRejectedCommandResponse(fmt.Errorf("Destination list: %v", err))
		}
		destination = &dst
	}
	for _, t := range opt.TaskIDs {
		ca.Logger.Info("Moving task", "source", opt.Source, "destination", opt.Destination, "task", t)
		err := ca.Commandbus.Execute(ctx, "list.move-task", command.MoveTaskArguments{
			Source:      opt.Source,
			Destination: opt.Destination,
			TaskID:      t,
		})
		if err != nil {
			ca.Logger.Info("Failed to execute command list.move-task", "error", err.Error())
			return api.NewRejectedCommandResponse(fmt.Errorf("Failed to move task %s", t))
		}
		if destination != nil && destination.Metadata.Name == "Done" {
			ca.Logger.Info("Completing task", "name", t)
			if err := ca.Commandbus.Execute(ctx, "task.complete", t); err != nil {
				ca.Logger.Info("Failed to execute command task.complete", "error", err.Error())
				return api.NewRejectedCommandResponse(fmt.Errorf("Failed to complete task %s", t))
			}
		}
		if source != nil && source.Metadata.Name == "Done" {
			ca.Logger.Info("Reopenning task", "name", t)
			if err := ca.Commandbus.Execute(ctx, "task.reopen", t); err != nil {
				ca.Logger.Info("Failed to execute command task.reopen", "error", err.Error())
				return api.NewRejectedCommandResponse(fmt.Errorf("Failed to reopen task %s", t))
			}
		}
	}
	return api.NewAcceptedCommandResponse("list", opt.Source)
}
