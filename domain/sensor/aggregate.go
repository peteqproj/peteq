// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    sensor, err := UnmarshalSensor(bytes)
//    bytes, err = sensor.Marshal()

package sensor

import "encoding/json"

func UnmarshalSensor(data []byte) (Sensor, error) {
	var r Sensor
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Sensor) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// sensor
type Sensor struct {
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
