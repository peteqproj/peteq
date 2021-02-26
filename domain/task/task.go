package task

import (
	"time"

	"github.com/peteqproj/peteq/pkg/repo"
)

type (
	// Spec of task
	Spec struct {
		DueDate   time.Time
		Completed bool `json:"completed" `
	}
)

// NewTask build repo.Resource
func NewTask(id string, name string, description string) repo.Resource {
	return repo.Resource{
		Metadata: repo.Metadata{
			Type:        "task",
			ID:          id,
			Name:        name,
			Description: description,
		},
	}
}
