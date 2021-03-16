// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    task, err := UnmarshalTask(bytes)
//    bytes, err = task.Marshal()

package task

import "encoding/json"

func UnmarshalTask(data []byte) (Task, error) {
	var r Task
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Task) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// task
type Task struct {
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
	Completed bool `json:"completed"`
}
