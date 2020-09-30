package project

import (
	"context"
	"io"

	"github.com/gofrs/uuid"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// CommandAPI for lists
	CommandAPI struct {
		Repo       *project.Repo
		Commandbus commandbus.CommandBus
		Logger     logger.Logger
	}

	// AddTasksRequestBody body spec of the request command AddTasks
	AddTasksRequestBody struct {
		TaskIDs []string `json:"tasks"`
		Project string   `json:"project"`
	}
)

// Create creates new project
func (ca *CommandAPI) Create(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	u := tenant.UserFromContext(ctx)
	proj := project.Project{}
	if err := api.UnmarshalInto(body, &proj); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	u2, err := uuid.NewV4()
	if err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}

	proj.Metadata.ID = u2.String()
	proj.Tenant = tenant.Tenant{
		ID:   u.Metadata.ID,
		Type: "User",
	}
	err = ca.Commandbus.ExecuteAndWait(ctx, "project.create", proj)
	if err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	return api.NewAcceptedCommandResponse("project", proj.Metadata.ID)
}

// AddTasks assign tasks to project
func (ca *CommandAPI) AddTasks(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	opt := AddTasksRequestBody{}
	if err := api.UnmarshalInto(body, &opt); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}

	for _, t := range opt.TaskIDs {
		err := ca.Commandbus.ExecuteAndWait(ctx, "project.add-task", command.AddTasksCommandOptions{
			Project: opt.Project,
			TaskID:  t,
		})
		if err != nil {
			return api.NewRejectedCommandResponse(err.Error())
		}
	}
	return api.NewAcceptedCommandResponse("project", opt.Project)
}
