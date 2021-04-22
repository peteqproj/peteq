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

	createRequestBody struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		Cron        string `json:"cron"`
	}
)

// Triggers sensor and executes all the bound automations
// @Description Sensor run
// @Tags Sensor Command API
// @Accept  json
// @Produce  json
// @Param body body sensorRunRequestBody true "Trigger sensor"
// @Success 200 {object} api.CommandResponse
// @Success 400 {object} api.CommandResponse
// @Router /c/sensor/trigger [post]
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

// Create sensor
// @Description returns the creates sensor id
// @Tags Sensor Command API
// @Accept  json
// @Produce  json
// @Param body body createRequestBody true "Create sensor"
// @Success 200 {object} api.CommandResponse
// @Success 400 {object} api.CommandResponse
// @Router /c/sensor/create [post]
// @Security ApiKeyAuth
func (c *CommandAPI) Create(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	spec := createRequestBody{}
	if err := api.UnmarshalInto(body, &spec); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	id, err := c.IDGenerator.GenerateV4()
	if err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	if err := c.Commandbus.Execute(ctx, "sensor.create", command.SensorCreateCommandOptions{
		ID:          id,
		Name:        spec.Name,
		Description: spec.Description,
		Cron:        &spec.Cron,
	}); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	return api.NewAcceptedCommandResponse("sensor", id)
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}
