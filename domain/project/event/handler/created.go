package handler

import (
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// CreatedHandler to handle task.created event
	CreatedHandler struct {
		Repo *project.Repo
	}
)

// Handle will process it the event
func (t *CreatedHandler) Handle(ev event.Event) error {
	opt := project.Project{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return t.Repo.Create(ev.Tenant.ID, opt)
}

func (t *CreatedHandler) Name() string {
	return "project_domain_CreatedHandler"
}
