package project

import "github.com/peteqproj/peteq/pkg/repo"

type (
	// Spec of a project
	Spec struct {
		Color    string   `json:"color" yaml:"color"`
		ImageURL string   `json:"imageUrl" yaml:"imageUrl"`
		Tasks    []string `json:"tasks" yaml:"tasks"`
	}
)

// NewProject build project resource
func NewProject(id string, name string, description string) repo.Resource {
	return repo.Resource{
		Metadata: repo.Metadata{
			Type:        "project",
			Name:        name,
			ID:          id,
			Description: description,
		},
		Spec: Spec{},
	}
}
