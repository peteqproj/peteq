// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    project, err := UnmarshalProject(bytes)
//    bytes, err = project.Marshal()

package project

import "encoding/json"

func UnmarshalProject(data []byte) (Project, error) {
	var r Project
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Project) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// project
type Project struct {
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
	Color    *string  `json:"color,omitempty"`   
	ImageURL *string  `json:"imageUrl,omitempty"`
	Tasks    []string `json:"tasks,omitempty"`   
}
