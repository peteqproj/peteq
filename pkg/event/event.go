package event

import "time"

type (
	// Event - something happend
	Event struct {
		Metadata Metadata    `json:"metadata" yaml:"metadata"`
		Spec     interface{} `json:"spec" yaml:"spec"`
	}

	// Metadata metadata on event
	Metadata struct {
		ID             string    `json:"id" yaml:"yaml"`
		Name           string    `json:"name" yaml:"name"`
		CreatedAt      time.Time `json:"createdAt" yaml:"createdAt"`
		AggregatorRoot string    `json:"aggregatorRoot" yaml:"aggregatorRoot"`
		AggregatorID   string    `json:"aggregatorId" yaml:"aggregatorId"`
	}
)
