package handler

import (
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// UpdatedHandler to handle task.created event
	UpdatedHandler struct {
		Repo *task.Repo
	}
)

// Handle will handle the event the process it
func (c *UpdatedHandler) Handle(ev event.Event) error {
	opt := task.Task{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return c.Repo.Update(ev.Tenant.ID, opt)
}

func (c *UpdatedHandler) Name() string {
	return "domain_UpdatedHandler"
}
