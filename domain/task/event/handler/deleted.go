package handler

import (
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// DeleteHandler to handle task.deleted event
	DeleteHandler struct {
		Repo *task.Repo
	}
)

// Handle will handle the event the process it
func (c *DeleteHandler) Handle(ev event.Event) error {
	opt := task.Task{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return c.Repo.Delete(opt.Tenant.ID, opt.Metadata.ID)
}

func (c *DeleteHandler) Name() string {
	return "domain_DeleteHandler"
}
