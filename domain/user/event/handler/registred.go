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

	// RegisteredSpec is the event.spec for this event
	RegisteredSpec struct {
		ID           string `json:"id" yaml:"id"`
		Email        string `json:"email" yaml:"email"`
		PasswordHash string `json:"passwordHash" yaml:"passwordHash"`
	}
)

// Handle will handle the event the process it
func (c *RegistredHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := RegisteredSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	return c.Repo.Create(user.User{
		Metadata: user.Metadata{
			ID:    opt.ID,
			Email: opt.Email,
		},
		Spec: user.Spec{
			PasswordHash: opt.PasswordHash,
		},
	})
}

func (c *RegistredHandler) Name() string {
	return "domain_user"
}
