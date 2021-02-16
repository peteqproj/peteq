package triggers

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	triggerEvent "github.com/peteqproj/peteq/domain/trigger/event/handler"
	triggerEventTypes "github.com/peteqproj/peteq/domain/trigger/event/types"
	userEventTypes "github.com/peteqproj/peteq/domain/user/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// ViewAPI for triggers view
	ViewAPI struct {
		DAL *DAL
	}

	triggersView struct {
		Triggers []triggerViewItem `json:"triggers"`
	}

	triggerViewItem struct {
		ID          string                        `json:"id"`
		Name        string                        `json:"name"`
		Description string                        `json:"description"`
		Type        string                        `json:"type"`
		Spec        interface{}                   `json:"spec"`
		History     []triggerExecutionHistoryItem `json:"history"`
	}

	triggerExecutionHistoryItem struct {
		TriggeredAt time.Time `json:"triggeredAt"`
		Manual      bool      `json:"manual"`
	}
)

// Get build triggers view
// @description Triggers View
// @tags View
// @produce  json
// @success 200 {object} triggersView
// @router /q/triggers [get]
// @Security ApiKeyAuth
func (b *ViewAPI) Get(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())
	view, err := b.DAL.load(c.Request.Context(), u.Metadata.ID)
	if err != nil {
		handleError(400, err, c)
		return
	}
	c.JSON(200, view)
}

func (h *ViewAPI) EventHandlers() map[string]handler.EventHandler {
	return map[string]handler.EventHandler{
		userEventTypes.UserRegistredEvent:     h,
		triggerEventTypes.TriggerCreatedEvent: h,
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
	case triggerEventTypes.TriggerCreatedEvent:
		{
			updated, err := h.handlerTriggerCreated(ctx, ev, current, logger)
			if err != nil {
				return err
			}
			return h.DAL.update(ctx, ev.Tenant.ID, updated)
		}
	}
	return nil
}
func (h *ViewAPI) Name() string {
	return "triggers_view"
}

func (h *ViewAPI) handlerUserRegistration(ctx context.Context, ev event.Event, logger logger.Logger) error {
	v := triggersView{
		Triggers: []triggerViewItem{},
	}
	return h.DAL.create(ctx, ev.Tenant.ID, v)
}

func (h *ViewAPI) handlerTriggerCreated(ctx context.Context, ev event.Event, view triggersView, looger logger.Logger) (triggersView, error) {
	spec := triggerEvent.CreatedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}
	triggerSpec := []string{}
	if spec.Cron != nil {
		triggerSpec = append(triggerSpec, *spec.Cron)
	}

	if spec.URL != nil {
		triggerSpec = append(triggerSpec, *spec.URL)
	}
	view.Triggers = append(view.Triggers, triggerViewItem{
		ID:          spec.ID,
		Name:        spec.Name,
		Description: spec.Description,
		Spec:        triggerSpec,
		History:     []triggerExecutionHistoryItem{},
	})
	return view, nil
}
