package handler

import (
	"fmt"

	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/pkg/event"
)

type (
	// TaskAddedHandler to handle task.created event
	TaskAddedHandler struct {
		Repo *project.Repo
	}
)

// Handle will process it the event
func (t *TaskAddedHandler) Handle(ev event.Event) error {
	fmt.Println("Handling event")
	opt, ok := ev.Spec.(command.AddTasksCommandOptions)
	if !ok {
		return fmt.Errorf("Failed to cast to Project object")
	}
	return t.Repo.AddTask(opt.Project, opt.TaskID)
}
