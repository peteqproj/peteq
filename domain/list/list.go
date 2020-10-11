package list

import "github.com/peteqproj/peteq/pkg/tenant"

type (
	// List holds set of tasks
	List struct {
		tenant.Tenant `json:"tenant" yaml:"tenant"`
		Metadata      Metadata `json:"metadata" yaml:"metadata"`
		Tasks         []string `json:"tasks" yaml:"tasks"`
	}

	// Metadata of list
	Metadata struct {
		ID    string `json:"id" yaml:"id"`
		Name  string `json:"name" yaml:"name"`
		Index int    `json:"index" yaml:"index"`
	}
)
