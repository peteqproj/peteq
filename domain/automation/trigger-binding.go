// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    triggerBinding, err := UnmarshalTriggerBinding(bytes)
//    bytes, err = triggerBinding.Marshal()

package automation

import "encoding/json"

func UnmarshalTriggerBinding(data []byte) (TriggerBinding, error) {
	var r TriggerBinding
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TriggerBinding) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Trigger binding aggregate
type TriggerBinding struct {
	TriggerBindingMetadata TriggerBindingMetadata `json:"triggerBindingMetadata"`
	TriggerBindingSpec     TriggerBindingSpec     `json:"triggerBindingSpec"`    
}

type TriggerBindingMetadata struct {
	Description *string           `json:"description,omitempty"`
	ID          string            `json:"id"`                   
	Labels      map[string]string `json:"labels,omitempty"`     
	Name        string            `json:"name"`                 
}

type TriggerBindingSpec struct {
	Automation string `json:"automation"`
	Trigger    string `json:"trigger"`   
}
