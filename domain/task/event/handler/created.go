package handler

type (
	// CreatedSpec is the event.spec for this event
	CreatedSpec struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}
)
