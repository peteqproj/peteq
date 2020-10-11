package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// Event - something happend
	Event struct {
		tenant.Tenant `json:"tenant" yaml:"tenant"`
		Metadata      Metadata    `json:"metadata" yaml:"metadata"`
		Spec          interface{} `json:"spec" yaml:"spec"`
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

// ToBytes returns bytes for the event
func (e Event) ToBytes() []byte {
	d, err := json.Marshal(e)
	if err != nil {
		return []byte{}
	}
	return d
}

// UnmarshalSpecInto cast spec into target object
func (e Event) UnmarshalSpecInto(t interface{}) error {
	spec, ok := e.Spec.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Failed to cast event.spec")
	}

	d, err := json.Marshal(spec)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, t)

}

// FromBytes builds event from bytes
func FromBytes(b []byte) Event {
	e := Event{}
	json.Unmarshal(b, &e)
	return e
}
