package list

type (
	// List holds set of tasks
	List struct {
		Metadata Metadata `json:"metadata" yaml:"metadata"`
		Tasks    []string `json:"tasks" yaml:"tasks"`
	}

	// Metadata of list
	Metadata struct {
		ID   string `json:"id" yaml:"id"`
		Name string `json:"name" yaml:"name"`
	}
)
