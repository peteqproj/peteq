package repo

type (
	// Resource is a representation of any resource that is used across the app
	Resource struct {
		Metadata Metadata `json:"metadata" yaml:"metadata"`
		Spec     interface{}
	}

	// Metadata is all resources common metadata
	Metadata struct {
		Type        string
		Name        string
		Description string
		ID          string
	}
)
