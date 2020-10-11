package handler

import (
	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/domain/user/command"

	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// LoggedinHandler to handle user.loggedin event
	LoggedinHandler struct {
		Repo *user.Repo
	}
)

// Handle will handle the event the process it
func (c *LoggedinHandler) Handle(ev event.Event) error {
	opt := command.LoginCommandOptions{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return err
	}
	u, err := c.Repo.Get(opt.UserID)
	if err != nil {
		return err
	}
	u.Spec.TokenHash = opt.HashedToken
	return c.Repo.Update(u)
}

func (c *LoggedinHandler) Name() string {
	return "domain_LoggedinHandler"
}
