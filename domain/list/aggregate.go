// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    list, err := UnmarshalList(bytes)
//    bytes, err = list.Marshal()

package list

import "encoding/json"

func UnmarshalList(data []byte) (List, error) {
	var r List
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *List) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// List aggregate
type List struct {
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
	Index float64  `json:"index"`
	Tasks []string `json:"tasks"`
}
