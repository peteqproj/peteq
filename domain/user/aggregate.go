// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    user, err := UnmarshalUser(bytes)
//    bytes, err = user.Marshal()

package user

import "encoding/json"

func UnmarshalUser(data []byte) (User, error) {
	var r User
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *User) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// user
type User struct {
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
	Email        string `json:"email"`       
	PasswordHash string `json:"passwordHash"`
	TokenHash    string `json:"tokenHash"`   
}
