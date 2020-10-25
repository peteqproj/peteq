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
)

type (
	// CommandAPI for lists
	CommandAPI struct {
		Repo       *list.Repo
		Commandbus commandbus.CommandBus
		Logger     logger.Logger
	}

	// MoveTasksRequestBody passed from http client
	MoveTasksRequestBody struct {
		Source      string   `json:"source"`
		Destination string   `json:"destination"`
		TaskIDs     []string `json:"tasks"`
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
		return api.NewRejectedCommandResponse(err.Error())
	}
	var source *list.List
	var destination *list.List
	if opt.Source != "" {
		src, err := ca.Repo.Get(u.Metadata.ID, opt.Source)
		if err != nil {
			return api.NewRejectedCommandResponse(fmt.Sprintf("Source list: %v", err))
		}
		source = &src
	}
	if opt.Destination != "" {
		dst, err := ca.Repo.Get(u.Metadata.ID, opt.Destination)
		if err != nil {
			return api.NewRejectedCommandResponse(fmt.Sprintf("Destination list: %v", err))
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
			return api.NewRejectedCommandResponse(err.Error())
		}
		if destination != nil && destination.Metadata.Name == "Done" {
			ca.Logger.Info("Completing task", "name", t)
			if err := ca.Commandbus.Execute(ctx, "task.complete", t); err != nil {
				return api.NewRejectedCommandResponse(err.Error())
			}
		}
		if source != nil && source.Metadata.Name == "Done" {
			ca.Logger.Info("Reopenning task", "name", t)
			if err := ca.Commandbus.Execute(ctx, "task.reopen", t); err != nil {
				return api.NewRejectedCommandResponse(err.Error())
			}
		}
	}
	return api.NewAcceptedCommandResponse("list", opt.Source)
}
