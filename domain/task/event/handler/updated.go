package handler

type (
	// UpdatedSpec is the event.spec for this event
	UpdatedSpec struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)
