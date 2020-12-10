package task

import (
	"context"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	listCommand "github.com/peteqproj/peteq/domain/list/command"
	projectCommand "github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/domain/task/command"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CommandAPI for tasks
	CommandAPI struct {
		Repo        *task.Repo
		Commandbus  commandbus.CommandBus
		Logger      logger.Logger
		IDGenerator utils.IDGenerator
	}

	// completeReopenTaskRequestBody for request body of two commands:
	// CompleteTask
	// ReopenTask
	completeReopenTaskRequestBody struct {
		Task string `json:"task"`
	}

	// deleteTaskRequestBody for delete command
	deleteTaskRequestBody struct {
		ID string `json:"id"`
	}

	createTaskRequestBody struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		Project     string `json:"project"`
		List        string `json:"list"`
	}
)

// Create creats tasks
func (c *CommandAPI) Create(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	spec := createTaskRequestBody{}
	err := api.UnmarshalInto(body, &spec)
	if err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	tid, err := c.IDGenerator.GenerateV4()
	if err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	if err := c.Commandbus.Execute(ctx, "task.create", command.CreateCommandOptions{
		ID:   tid,
		Name: spec.Name,
	}); err != nil {
		return api.NewRejectedCommandResponse(err)
	}

	if spec.Project != "" {
		if err := c.Commandbus.Execute(ctx, "project.add-task", projectCommand.AddTasksCommandOptions{
			Project: spec.Project,
			TaskID:  tid,
		}); err != nil {
			return api.NewRejectedCommandResponse(err)
		}
	}

	if spec.List != "" {
		if err := c.Commandbus.Execute(ctx, "list.move-task", listCommand.MoveTaskArguments{
			Source:      "",
			Destination: spec.List,
			TaskID:      tid,
		}); err != nil {
			return api.NewRejectedCommandResponse(err)
		}
	}
	return api.NewAcceptedCommandResponse("task", tid)
}

// Update task
func (c *CommandAPI) Update(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	t := &task.Task{}
	err := api.UnmarshalInto(body, t)
	if err != nil {
		return api.NewRejectedCommandResponse(err)
	}

	if err := validator.New().Struct(t); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	if err := c.Commandbus.Execute(ctx, "task.update", *t); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	return api.NewAcceptedCommandResponse("task", t.Metadata.ID)
}

// Delete store new task
func (c *CommandAPI) Delete(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	u := tenant.UserFromContext(ctx)
	req := &deleteTaskRequestBody{}
	err := api.UnmarshalInto(body, req)
	t, err := c.Repo.Get(u.Metadata.ID, req.ID)
	if err != nil {
		return api.NewRejectedCommandResponse(err)
	}

	if err := c.Commandbus.Execute(ctx, "task.delete", t); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	return api.NewAcceptedCommandResponse("task", t.Metadata.ID)
}

// Complete task
func (c *CommandAPI) Complete(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	req := &completeReopenTaskRequestBody{}
	if err := api.UnmarshalInto(body, req); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	if err := c.Commandbus.Execute(ctx, "task.complete", req.Task); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	return api.NewAcceptedCommandResponse("task", req.Task)
}

// Reopen task
func (c *CommandAPI) Reopen(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	req := &completeReopenTaskRequestBody{}
	if err := api.UnmarshalInto(body, req); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	if err := c.Commandbus.Execute(ctx, "task.reopen", req.Task); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	return api.NewAcceptedCommandResponse("task", req.Task)
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}
