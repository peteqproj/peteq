package handler

import (
	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/list/command"
	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// CreatedHandler to handle list.created event
	CreatedHandler struct {
		Repo *list.Repo
	}
)

// Handle will process it the event
func (t *CreatedHandler) Handle(ev event.Event) error {
	opt := command.CreateCommandOptions{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return t.Repo.Create(ev.Tenant.ID, list.List{
		Tenant: ev.Tenant,
		Metadata: list.Metadata{
			ID:    opt.ID,
			Name:  opt.Name,
			Index: opt.Index,
		},
		Tasks: []string{},
	})
}

func (t *CreatedHandler) Name() string {
	return "list_domain_CreatedHandler"
}
