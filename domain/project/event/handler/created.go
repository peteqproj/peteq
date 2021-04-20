package handler

type (
	// CreatedSpec is the event.spec for this event
	CreatedSpec struct {
		ID          string `json:"id" yaml:"id"`
		Name        string `json:"name" yaml:"name"`
		Description string `json:"description" yaml:"description"`
		Color       string `json:"color" yaml:"color"`
		ImageURL    string `json:"imageUrl" yaml:"imageUrl"`
	}
)
