package handler

import (
	"fmt"

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
	opt, ok := ev.Spec.(command.LoginCommandOptions)
	if !ok {
		return fmt.Errorf("Failed to convert arguments to User")
	}
	u, err := c.Repo.Get(opt.UserID)
	if err != nil {
		return err
	}
	u.Spec.TokenHash = opt.HashedToken
	return c.Repo.Update(u)
}
