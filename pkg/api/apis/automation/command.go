package automation

import (
	"context"
	"io"

	"github.com/peteqproj/peteq/domain/automation"
	automationCommands "github.com/peteqproj/peteq/domain/automation/command"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CommandAPI for automations
	CommandAPI struct {
		Repo        *automation.Repo
		Commandbus  commandbus.CommandBus
		Logger      logger.Logger
		IDGenerator utils.IDGenerator
	}

	CreateAutomationRequestBody struct {
		Name string `json:"name" validate:"required"`
	} //@name CreateAutomationRequestBody

	CreateSensorAutomationBindingRequestBody struct {
		Name       string `json:"name" validate:"required"`
		Automation string `json:"automation" validate:"required"`
		Sensor     string `json:"sensor" validate:"required"`
	} //@name CreateSensorAutomationBindingRequestBody
)

// Creates automation
// @Tags Automation Command API
// @Accept  json
// @Produce  json
// @Param body body CreateAutomationRequestBody true "Creates automation"
// @Success 200 {object} api.CommandResponse
// @Success 400 {object} api.CommandResponse
// @Router /c/automation/create [post]
// @Security ApiKeyAuth
func (ca *CommandAPI) Create(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	opt := &CreateAutomationRequestBody{}
	if err := api.UnmarshalInto(body, opt); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	id, err := ca.IDGenerator.GenerateV4()
	if err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	if err := ca.Commandbus.Execute(ctx, "automation.create", automationCommands.AutomationCreateCommandOptions{
		ID:   id,
		Name: opt.Name,
		Type: "rss-importer",
	}); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	return api.NewAcceptedCommandResponse("automation", id)
}

// Binds Sensor to Automation
// @Tags Automation Command API
// @Accept  json
// @Produce  json
// @Param body body CreateSensorAutomationBindingRequestBody true "Binds Sensor to Automation"
// @Success 200 {object} api.CommandResponse
// @Success 400 {object} api.CommandResponse
// @Router /c/automation/bindSensor [post]
// @Security ApiKeyAuth
func (ca *CommandAPI) CreateSensorAutomationBinding(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	opt := &CreateSensorAutomationBindingRequestBody{}
	if err := api.UnmarshalInto(body, opt); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	id, err := ca.IDGenerator.GenerateV4()
	if err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	if err := ca.Commandbus.Execute(ctx, "automation.bindSensor", automationCommands.SensorBindingCreateCommandOptions{
		ID:         id,
		Name:       opt.Name,
		Sensor:     opt.Sensor,
		Automation: opt.Automation,
	}); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	return api.NewAcceptedCommandResponse("sensor-automation-binding", id)
}
