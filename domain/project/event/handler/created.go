package handler

import (
	"fmt"

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
	opt, ok := ev.Spec.(project.Project)
	if !ok {
		return fmt.Errorf("Failed to cast to Project object")
	}
	return t.Repo.Create(opt)
}
