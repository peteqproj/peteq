package handler

import (
	"context"

	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// LoggedinHandler to handle user.loggedin event
	LoggedinHandler struct {
		Repo *user.Repo
	}

	// LoggedinSpec is the event.spec for this event
	LoggedinSpec struct {
		ID        string `json:"id" yaml:"id"`
		TokenHash string `json:"tokenHash" yaml:"tokenHash"`
	}
)

// Handle will handle the event the process it
func (c *LoggedinHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	opt := LoggedinSpec{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	u, err := c.Repo.Get(opt.ID)
	if err != nil {
		return err
	}
	u.Spec.TokenHash = opt.TokenHash
	return c.Repo.Update(u)
}

func (c *LoggedinHandler) Name() string {
	return "domain_LoggedinHandler"
}
