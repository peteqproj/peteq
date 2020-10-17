package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// RegistredHandler to handle user.registred event
	RegistredHandler struct {
		Repo *user.Repo
	}
)

// Handle will handle the event the process it
func (c *RegistredHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := user.User{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return c.Repo.Create(opt)
}

func (c *RegistredHandler) Name() string {
	return "domain_RegistredHandler"
}
