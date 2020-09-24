package project

type (
	// Project holds set of tasks related to common goal
	Project struct {
		Metadata Metadata `json:"metadata" yaml:"metadata"`
		Tasks    []string `json:"tasks" yaml:"tasks"`
	}

	// Metadata of project
	Metadata struct {
		ID          string `json:"id" yaml:"id"`
		Name        string `json:"name" yaml:"name"`
		Description string `json:"description" yaml:"description" validate:"required"`
		Color       string `json:"color" yaml:"color"`
		ImageURL    string `json:"imageUrl" yaml:"imageUrl"`
	}
)
