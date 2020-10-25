package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// CreatedHandler to handle list.created event
	CreatedHandler struct {
		Repo *list.Repo
	}
	// CreatedSpec is the event.spec for this event
	CreatedSpec struct {
		ID    string `json:"id" yaml:"id"`
		Name  string `json:"name" yaml:"name"`
		Index int    `json:"index" yaml:"index"`
	}
)

// Handle will process it the event
func (t *CreatedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := CreatedSpec{}
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
	return "domain_list"
}
