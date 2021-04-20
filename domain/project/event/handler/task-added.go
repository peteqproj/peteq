package handler

type (
	// TaskAddedSpec is the event.spec for this event
	TaskAddedSpec struct {
		TaskID  string `json:"taskId" yaml:"taskId"`
		Project string `json:"project" yaml:"project"`
	}
)
