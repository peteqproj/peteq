package task

import "time"

type (
	// Task entity
	Task struct {
		Metadata Metadata `json:"metadata" yaml:"metadata" validate:"required"`
		Spec     Spec     `json:"spec" yaml:"spec" validate:"required"`
		Status   Status   `json:"status" yaml:"status"`
	}

	// Metadata of task
	Metadata struct {
		ID          string `json:"id" yaml:"id" validate:"required"`
		Name        string `json:"name" yaml:"name" validate:"required"`
		Description string `json:"description" yaml:"description"`
	}

	// Spec of task
	Spec struct {
		DueDate time.Time
	}

	// Status of task
	Status struct {
		Completed bool `json:"completed" yaml:"completed"`
	}
)
