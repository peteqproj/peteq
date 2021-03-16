// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    automation, err := UnmarshalAutomation(bytes)
//    bytes, err = automation.Marshal()
//
//    triggerBinding, err := UnmarshalTriggerBinding(bytes)
//    bytes, err = triggerBinding.Marshal()

package automation

import "encoding/json"

func UnmarshalAutomation(data []byte) (Automation, error) {
	var r Automation
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Automation) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalTriggerBinding(data []byte) (TriggerBinding, error) {
	var r TriggerBinding
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TriggerBinding) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Automation aggregate
type Automation struct {
	Metadata Metadata       `json:"metadata"`
	Spec     AutomationSpec `json:"spec"`    
}

type Metadata struct {
	Description *string           `json:"description,omitempty"`
	ID          string            `json:"id"`                   
	Labels      map[string]string `json:"labels,omitempty"`     
	Name        string            `json:"name"`                 
}

type AutomationSpec struct {
	JSONInputSchema string `json:"jsonInputSchema"`
	Type            string `json:"type"`           
}

// Trigger binding aggregate
type TriggerBinding struct {
	Metadata Metadata           `json:"metadata"`
	Spec     TriggerBindingSpec `json:"spec"`    
}

type TriggerBindingSpec struct {
	Automation string `json:"automation"`
	Trigger    string `json:"trigger"`   
}
