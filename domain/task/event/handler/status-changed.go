package handler

type (
	// StatusChangedSpec is the event.spec for this event
	StatusChangedSpec struct {
		Completed bool `json:"completed"`
	}
)
