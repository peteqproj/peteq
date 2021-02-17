package project

import (
	"context"
	"fmt"
	"io"

	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CommandAPI for lists
	CommandAPI struct {
		Repo        *project.Repo
		Commandbus  commandbus.CommandBus
		Logger      logger.Logger
		IDGenerator utils.IDGenerator
	}

	// AddTasksRequestBody body spec of the request command AddTasks
	AddTasksRequestBody struct {
		Project string   `json:"project" validate:"required"`
		TaskIDs []string `json:"tasks" validate:"required"`
	} //@name AddTasksRequestBody

	// CreateProjectRequestBody body spec of the request command Create
	CreateProjectRequestBody struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		Color       string `json:"color"`
		ImageURL    string `json:"imageUrl"`
	} //@name CreateProjectRequestBody
)

// Create creates new project
// @Description Create project
// @Tags Project Command API
// @Accept  json
// @Produce  json
// @Param body body CreateProjectRequestBody true "Add tasks into project"
// @Success 200 {object} api.CommandResponse
// @Success 400 {object} api.CommandResponse
// @Router /c/project/create [post]
// @Security ApiKeyAuth
func (ca *CommandAPI) Create(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	u := tenant.UserFromContext(ctx)
	spec := CreateProjectRequestBody{}
	if err := api.UnmarshalInto(body, &spec); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	uid, err := ca.IDGenerator.GenerateV4()
	if err != nil {
		return api.NewRejectedCommandResponse(err)
	}

	proj := project.Project{
		Tenant: tenant.Tenant{
			ID:   u.Metadata.ID,
			Type: "User",
		},
		Metadata: project.Metadata{
			ID:          uid,
			Name:        spec.Name,
			Description: spec.Description,
			Color:       spec.Color,
			ImageURL:    spec.ImageURL,
		},
		Tasks: []string{},
	}
	err = ca.Commandbus.Execute(ctx, "project.create", proj)
	if err != nil {
		ca.Logger.Info("Failed to execute project.create command", "error", err.Error())
		return api.NewRejectedCommandResponse(fmt.Errorf("Failed to create project"))
	}
	return api.NewAcceptedCommandResponse("project", proj.Metadata.ID)
}

// AddTasks assign tasks to project
// @Description Add tasks into project
// @Tags Project Command API
// @Accept  json
// @Produce  json
// @Param body body AddTasksRequestBody true "Move tasks from source to destination list"
// @Success 200 {object} api.CommandResponse
// @Success 400 {object} api.CommandResponse
// @Router /c/list/moveTasks [post]
// @Security ApiKeyAuth
func (ca *CommandAPI) AddTasks(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	opt := AddTasksRequestBody{}
	if err := api.UnmarshalInto(body, &opt); err != nil {
		return api.NewRejectedCommandResponse(err)
	}

	for _, t := range opt.TaskIDs {
		err := ca.Commandbus.Execute(ctx, "project.add-task", command.AddTasksCommandOptions{
			Project: opt.Project,
			TaskID:  t,
		})
		if err != nil {
			ca.Logger.Info("Failed to execute command project.add-task", "error", err.Error())
			return api.NewRejectedCommandResponse(fmt.Errorf("Failed to add task %s to project %s", t, opt.Project))
		}
	}
	return api.NewAcceptedCommandResponse("project", opt.Project)
}
