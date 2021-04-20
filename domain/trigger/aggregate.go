// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    trigger, err := UnmarshalTrigger(bytes)
//    bytes, err = trigger.Marshal()

package trigger

import "encoding/json"

func UnmarshalTrigger(data []byte) (Trigger, error) {
	var r Trigger
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Trigger) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// trigger
type Trigger struct {
	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec"`
}

type Metadata struct {
	Description *string           `json:"description,omitempty"`
	ID          string            `json:"id"`
	Labels      map[string]string `json:"labels,omitempty"`
	Name        string            `json:"name"`
}

type Spec struct {
	Cron *string `json:"cron,omitempty"`
}
