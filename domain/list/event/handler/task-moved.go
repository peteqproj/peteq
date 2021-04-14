package handler

type (
	// TaskMovedSpec is the event.spec for this event
	TaskMovedSpec struct {
		TaskID      string `json:"taskId" yaml:"taskId"`
		Source      string `json:"source" yaml:"source"`
		Destination string `json:"destination" yaml:"destination"`
	}
)
