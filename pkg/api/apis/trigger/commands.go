package trigger

import (
	"context"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/trigger"
	"github.com/peteqproj/peteq/domain/trigger/command"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CommandAPI for tasks
	CommandAPI struct {
		Repo        *trigger.Repo
		Commandbus  commandbus.CommandBus
		Logger      logger.Logger
		IDGenerator utils.IDGenerator
	}

	triggerRunRequestBody struct {
		ID string `json:"id" validate:"required"`
	}
)

// Run creats tasks
func (c *CommandAPI) Run(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	spec := triggerRunRequestBody{}
	err := api.UnmarshalInto(body, &spec)
	if err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	if err := c.Commandbus.Execute(ctx, "trigger.run", command.TriggerRunCommandOptions{
		ID: spec.ID,
	}); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	return api.NewAcceptedCommandResponse("trigger", spec.ID)
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}
