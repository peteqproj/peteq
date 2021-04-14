package handler

type (
	// CreatedSpec is the event.spec for this event
	CreatedSpec struct {
		ID    string `json:"id" yaml:"id"`
		Name  string `json:"name" yaml:"name"`
		Index int    `json:"index" yaml:"index"`
	}
)
