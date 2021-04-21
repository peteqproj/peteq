// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    automation, err := UnmarshalAutomation(bytes)
//    bytes, err = automation.Marshal()
//
//    sensorBinding, err := UnmarshalSensorBinding(bytes)
//    bytes, err = sensorBinding.Marshal()

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

func UnmarshalSensorBinding(data []byte) (SensorBinding, error) {
	var r SensorBinding
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SensorBinding) Marshal() ([]byte, error) {
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

// Sensor binding aggregate
type SensorBinding struct {
	Metadata Metadata          `json:"metadata"`
	Spec     SensorBindingSpec `json:"spec"`
}

type SensorBindingSpec struct {
	Automation string `json:"automation"`
	Sensor     string `json:"sensor"`
}
