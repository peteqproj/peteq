package handler

import (
	"fmt"

	"github.com/peteqproj/peteq/domain/user"

	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// RegistredHandler to handle user.registred event
	RegistredHandler struct {
		Repo *user.Repo
	}
)

// Handle will handle the event the process it
func (c *RegistredHandler) Handle(ev event.Event) error {
	u, ok := ev.Spec.(user.User)
	if !ok {
		return fmt.Errorf("Failed to convert arguments to User")
	}
	return c.Repo.Create(u)
}
