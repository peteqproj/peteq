package sensor

import (
	"context"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/sensor"
	"github.com/peteqproj/peteq/domain/sensor/command"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CommandAPI for tasks
	CommandAPI struct {
		Repo        *sensor.Repo
		Commandbus  commandbus.CommandBus
		Logger      logger.Logger
		IDGenerator utils.IDGenerator
	}

	sensorRunRequestBody struct {
		ID   string      `json:"id" validate:"required"`
		Data interface{} `json:"data"`
	} //@name SensorRunRequestBody
)

// Run runs sensor
// @Description Sensor run
// @Tags Sensor Command API
// @Accept  json
// @Produce  json
// @Param body body sensorRunRequestBody true "Sensor trigger"
// @Success 200 {object} api.CommandResponse
// @Success 400 {object} api.CommandResponse
// @Router /c/sensor/run [post]
// @Security ApiKeyAuth
func (c *CommandAPI) Run(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	spec := sensorRunRequestBody{}
	err := api.UnmarshalInto(body, &spec)
	if err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	if err := c.Commandbus.Execute(ctx, "sensor.trigger", command.SensorTriggerCommandOptions{
		ID:   spec.ID,
		Data: spec.Data,
	}); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	return api.NewAcceptedCommandResponse("sensor", spec.ID)
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}
