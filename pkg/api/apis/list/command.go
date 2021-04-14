package list

import (
	"context"
	"fmt"
	"io"

	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/list/command"
	taskCommands "github.com/peteqproj/peteq/domain/task/command"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
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
		Source      string   `json:"source"`
		Destination string   `json:"destination"`
		TaskIDs     []string `json:"tasks" validate:"required"`
	} //@name MoveTasksRequestBody
)

// MoveTasks move multiple tasks from one list to another
// Complete task that moved into done list and opens task that moved from done list
// this should be done in different place
// @Description Move tasks from source to destination list
// @Tags List Command API
// @Accept  json
// @Produce  json
// @Param body body MoveTasksRequestBody true "Move tasks from source to destination list"
// @Success 200 {object} api.CommandResponse
// @Success 400 {object} api.CommandResponse
// @Router /c/list/moveTasks [post]
// @Security ApiKeyAuth
func (ca *CommandAPI) MoveTasks(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	opt := &MoveTasksRequestBody{}
	if err := api.UnmarshalInto(body, opt); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	var source *list.List
	var destination *list.List
	if opt.Source != "" {
		src, err := ca.Repo.GetById(ctx, opt.Source)
		if err != nil {
			return api.NewRejectedCommandResponse(fmt.Errorf("Source list: %v", err))
		}
		source = src
	}
	if opt.Destination != "" {
		dst, err := ca.Repo.GetById(ctx, opt.Destination)
		if err != nil {
			return api.NewRejectedCommandResponse(fmt.Errorf("Destination list: %v", err))
		}
		destination = dst
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
			if err := ca.Commandbus.Execute(ctx, "task.complete", taskCommands.CompleteTaskArguments{
				TaskID: t,
			}); err != nil {
				ca.Logger.Info("Failed to execute command task.complete", "error", err.Error())
				return api.NewRejectedCommandResponse(fmt.Errorf("Failed to complete task %s", t))
			}
		}
		if source != nil && source.Metadata.Name == "Done" {
			ca.Logger.Info("Reopenning task", "name", t)
			if err := ca.Commandbus.Execute(ctx, "task.reopen", taskCommands.ReopenTaskArguments{
				TaskID: t,
			}); err != nil {
				ca.Logger.Info("Failed to execute command task.reopen", "error", err.Error())
				return api.NewRejectedCommandResponse(fmt.Errorf("Failed to reopen task %s", t))
			}
		}
	}
	return api.NewAcceptedCommandResponse("list", opt.Source)
}
