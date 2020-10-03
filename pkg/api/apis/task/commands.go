package task

import (
	"context"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// CommandAPI for tasks
	CommandAPI struct {
		Repo       *task.Repo
		Commandbus commandbus.CommandBus
		Logger     logger.Logger
	}

	// CompleteReopenTaskRequestBody for request body of two commands:
	// CompleteTask
	// ReopenTask
	CompleteReopenTaskRequestBody struct {
		Task string `json:"task"`
	}

	// DeleteTaskRequestBody for delete command
	DeleteTaskRequestBody struct {
		ID string `json:"id"`
	}
)

// Create creats tasks
func (c *CommandAPI) Create(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	u := tenant.UserFromContext(ctx)
	t := &task.Task{}
	err := api.UnmarshalInto(body, t)
	if err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	u2, err := uuid.NewV4()
	if err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	t.Metadata.ID = u2.String()
	t.Tenant = tenant.Tenant{
		ID:   u.Metadata.ID,
		Type: "User",
	}

	if err := validator.New().Struct(t); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	if err := c.Commandbus.ExecuteAndWait(ctx, "task.create", *t); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	return api.NewAcceptedCommandResponse("task", t.Metadata.ID)
}

// Update task
func (c *CommandAPI) Update(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	t := &task.Task{}
	err := api.UnmarshalInto(body, t)
	if err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}

	if err := validator.New().Struct(t); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	if err := c.Commandbus.ExecuteAndWait(ctx, "task.update", *t); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	return api.NewAcceptedCommandResponse("task", t.Metadata.ID)
}

// Delete store new task
func (c *CommandAPI) Delete(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	u := tenant.UserFromContext(ctx)
	req := &DeleteTaskRequestBody{}
	err := api.UnmarshalInto(body, req)
	t, err := c.Repo.Get(u.Metadata.ID, req.ID)
	if err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}

	if err := c.Commandbus.ExecuteAndWait(ctx, "task.delete", t); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	return api.NewAcceptedCommandResponse("task", t.Metadata.ID)
}

// Complete task
func (c *CommandAPI) Complete(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	req := &CompleteReopenTaskRequestBody{}
	if err := api.UnmarshalInto(body, req); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	if err := c.Commandbus.ExecuteAndWait(ctx, "task.complete", req.Task); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	return api.NewAcceptedCommandResponse("task", req.Task)
}

// Reopen task
func (c *CommandAPI) Reopen(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	req := &CompleteReopenTaskRequestBody{}
	if err := api.UnmarshalInto(body, req); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	if err := c.Commandbus.ExecuteAndWait(ctx, "task.reopen", req.Task); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	return api.NewAcceptedCommandResponse("task", req.Task)
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}