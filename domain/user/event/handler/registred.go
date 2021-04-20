package handler

type (
	// RegisteredSpec is the event.spec for this event
	RegisteredSpec struct {
		ID    string `json:"id" yaml:"id"`
		Email string `json:"email" yaml:"email"`
	}
)
