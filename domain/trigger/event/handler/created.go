package handler

type (

	// CreatedSpec is the event.spec for this event
	CreatedSpec struct {
		ID          string  `json:"id" yaml:"id"`
		Name        string  `json:"name" yaml:"name"`
		Description string  `json:"description" yaml:"description"`
		Cron        *string `json:"cron,omitempty" yaml:"cron,omitempty"`
	}
)
