package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/automation"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// TriggerBindingCreatedHandler to handle task.created event
	TriggerBindingCreatedHandler struct {
		Repo *automation.Repo
	}

	// TriggerBindingCreatedSpec is the event.spec for this event
	TriggerBindingCreatedSpec struct {
		ID         string `json:"id" yaml:"id"`
		Name       string `json:"name" yaml:"name"`
		Trigger    string `json:"trigger" yaml:"trigger"`
		Automation string `json:"automation" yaml:"automation"`
	}
)

// Handle will process it the event
func (t *TriggerBindingCreatedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := TriggerBindingCreatedSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	spec := automation.TriggerBindingSpec{
		Automation: opt.Automation,
		Trigger:    opt.Trigger,
	}

	return t.Repo.CreateTriggerBinding(ev.Tenant.ID, automation.TriggerBinding{
		Tenant: ev.Tenant,
		Metadata: automation.TriggerBindingMetadata{
			ID:   opt.ID,
			Name: opt.Name,
		},
		Spec: spec,
	})
}

func (t *TriggerBindingCreatedHandler) Name() string {
	return "domain_automation"
}
