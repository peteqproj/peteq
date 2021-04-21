package handler

type (
	// SensorBindingCreatedSpec is the event.spec for this event
	SensorBindingCreatedSpec struct {
		ID         string `json:"id" yaml:"id"`
		Name       string `json:"name" yaml:"name"`
		Sensor     string `json:"sensor" yaml:"sensor"`
		Automation string `json:"automation" yaml:"automation"`
	}
)
