package sensors

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	sensorEvent "github.com/peteqproj/peteq/domain/sensor/event/handler"
	sensorEventTypes "github.com/peteqproj/peteq/domain/sensor/event/types"
	userEventTypes "github.com/peteqproj/peteq/domain/user/event/types"
	"github.com/peteqproj/peteq/pkg/api/auth"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// ViewAPI for sensors view
	ViewAPI struct {
		DAL *DAL
	}

	sensorsView struct {
		Sensors []sensorViewItem `json:"sensors"`
	}

	sensorViewItem struct {
		ID          string      `json:"id"`
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Type        string      `json:"type"`
		Spec        interface{} `json:"spec"`
	}
)

// Get build sensors view
// @description Sensors View
// @tags View
// @produce  json
// @success 200 {object} sensorsView
// @router /q/sensors [get]
// @Security ApiKeyAuth
func (b *ViewAPI) Get(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())
	if u == nil {
		auth.UnauthorizedResponse(c)
		return
	}
	view, err := b.DAL.load(c.Request.Context(), u.Metadata.ID)
	if err != nil {
		handleError(400, err, c)
		return
	}
	c.JSON(200, view)
}

func (h *ViewAPI) EventHandlers() map[string]handler.EventHandler {
	return map[string]handler.EventHandler{
		userEventTypes.UserRegistredEvent:   h,
		sensorEventTypes.SensorCreatedEvent: h,
	}
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}

func (h *ViewAPI) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	if ev.Metadata.Name == userEventTypes.UserRegistredEvent {
		return h.handlerUserRegistration(ctx, ev, logger)
	}

	current, err := h.DAL.load(ctx, ev.Tenant.ID)
	if err != nil {
		return err
	}
	switch ev.Metadata.Name {
	case sensorEventTypes.SensorCreatedEvent:
		{
			updated, err := h.handlerSensorCreated(ctx, ev, current, logger)
			if err != nil {
				return err
			}
			return h.DAL.update(ctx, ev.Tenant.ID, updated)
		}
	}
	return nil
}
func (h *ViewAPI) Name() string {
	return "sensors_view"
}

func (h *ViewAPI) handlerUserRegistration(ctx context.Context, ev event.Event, logger logger.Logger) error {
	v := sensorsView{
		Sensors: []sensorViewItem{},
	}
	return h.DAL.create(ctx, ev.Tenant.ID, v)
}

func (h *ViewAPI) handlerSensorCreated(ctx context.Context, ev event.Event, view sensorsView, looger logger.Logger) (sensorsView, error) {
	spec := sensorEvent.CreatedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}
	sensorSpec := []string{}
	if spec.Cron != nil {
		sensorSpec = append(sensorSpec, *spec.Cron)
	}

	view.Sensors = append(view.Sensors, sensorViewItem{
		ID:          spec.ID,
		Name:        spec.Name,
		Description: spec.Description,
		Spec:        sensorSpec,
	})
	return view, nil
}
