package list

import (
	"context"
	"io"

	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/list/command"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
)

type (
	// CommandAPI for lists
	CommandAPI struct {
		Repo       *list.Repo
		Commandbus commandbus.CommandBus
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
	opt := &MoveTasksRequestBody{}
	if err := api.UnmarshalInto(body, opt); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	var source *list.List
	var destination *list.List
	if opt.Source != "" {
		src, err := ca.Repo.Get(opt.Source)
		if err != nil {
			return api.NewRejectedCommandResponse(err.Error())
		}
		source = &src
	}
	if opt.Destination != "" {
		dst, err := ca.Repo.Get(opt.Destination)
		if err != nil {
			return api.NewRejectedCommandResponse(err.Error())
		}
		destination = &dst
	}
	for _, t := range opt.TaskIDs {
		err := ca.Commandbus.ExecuteAndWait(ctx, "list.move-task", command.MoveTaskArguments{
			Source:      opt.Source,
			Destination: opt.Destination,
			TaskID:      t,
		})
		if err != nil {
			return api.NewRejectedCommandResponse(err.Error())
		}
		if destination != nil && destination.Metadata.Name == "Done" {
			if err := ca.Commandbus.ExecuteAndWait(ctx, "task.complete", t); err != nil {
				return api.NewRejectedCommandResponse(err.Error())
			}
		}
		if source != nil && source.Metadata.Name == "Done" {
			if err := ca.Commandbus.ExecuteAndWait(ctx, "task.reopen", t); err != nil {
				return api.NewRejectedCommandResponse(err.Error())
			}
		}
	}
	return api.NewAcceptedCommandResponse("list", opt.Source)
}
