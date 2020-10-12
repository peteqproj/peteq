package handler

import (
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// CreatedHandler to handle task.created event
	CreatedHandler struct {
		Repo *task.Repo
	}
)

// Handle will handle the event the process it
func (c *CreatedHandler) Handle(ev event.Event) error {
	opt := task.Task{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return c.Repo.Create(ev.Tenant.ID, opt)
}

func (c *CreatedHandler) Name() string {
	return "task_domain_CreatedHandler"
}
