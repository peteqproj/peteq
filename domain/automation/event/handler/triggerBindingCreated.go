package handler

type (
	// TriggerBindingCreatedSpec is the event.spec for this event
	TriggerBindingCreatedSpec struct {
		ID         string `json:"id" yaml:"id"`
		Name       string `json:"name" yaml:"name"`
		Trigger    string `json:"trigger" yaml:"trigger"`
		Automation string `json:"automation" yaml:"automation"`
	}
)
