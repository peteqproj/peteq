/*
 * Peteq API
 *
 * Peteq OpenAPI spec.
 *
 * API version: 1.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
)

// MoveTasksRequestBody struct for MoveTasksRequestBody
type MoveTasksRequestBody struct {
	Destination *string `json:"destination,omitempty"`
	Source *string `json:"source,omitempty"`
	Tasks []string `json:"tasks"`
}

// NewMoveTasksRequestBody instantiates a new MoveTasksRequestBody object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMoveTasksRequestBody(tasks []string, ) *MoveTasksRequestBody {
	this := MoveTasksRequestBody{}
	this.Tasks = tasks
	return &this
}

// NewMoveTasksRequestBodyWithDefaults instantiates a new MoveTasksRequestBody object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMoveTasksRequestBodyWithDefaults() *MoveTasksRequestBody {
	this := MoveTasksRequestBody{}
	return &this
}

// GetDestination returns the Destination field value if set, zero value otherwise.
func (o *MoveTasksRequestBody) GetDestination() string {
	if o == nil || o.Destination == nil {
		var ret string
		return ret
	}
	return *o.Destination
}

// GetDestinationOk returns a tuple with the Destination field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MoveTasksRequestBody) GetDestinationOk() (*string, bool) {
	if o == nil || o.Destination == nil {
		return nil, false
	}
	return o.Destination, true
}

// HasDestination returns a boolean if a field has been set.
func (o *MoveTasksRequestBody) HasDestination() bool {
	if o != nil && o.Destination != nil {
		return true
	}

	return false
}

// SetDestination gets a reference to the given string and assigns it to the Destination field.
func (o *MoveTasksRequestBody) SetDestination(v string) {
	o.Destination = &v
}

// GetSource returns the Source field value if set, zero value otherwise.
func (o *MoveTasksRequestBody) GetSource() string {
	if o == nil || o.Source == nil {
		var ret string
		return ret
	}
	return *o.Source
}

// GetSourceOk returns a tuple with the Source field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MoveTasksRequestBody) GetSourceOk() (*string, bool) {
	if o == nil || o.Source == nil {
		return nil, false
	}
	return o.Source, true
}

// HasSource returns a boolean if a field has been set.
func (o *MoveTasksRequestBody) HasSource() bool {
	if o != nil && o.Source != nil {
		return true
	}

	return false
}

// SetSource gets a reference to the given string and assigns it to the Source field.
func (o *MoveTasksRequestBody) SetSource(v string) {
	o.Source = &v
}

// GetTasks returns the Tasks field value
func (o *MoveTasksRequestBody) GetTasks() []string {
	if o == nil  {
		var ret []string
		return ret
	}

	return o.Tasks
}

// GetTasksOk returns a tuple with the Tasks field value
// and a boolean to check if the value has been set.
func (o *MoveTasksRequestBody) GetTasksOk() (*[]string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Tasks, true
}

// SetTasks sets field value
func (o *MoveTasksRequestBody) SetTasks(v []string) {
	o.Tasks = v
}

func (o MoveTasksRequestBody) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Destination != nil {
		toSerialize["destination"] = o.Destination
	}
	if o.Source != nil {
		toSerialize["source"] = o.Source
	}
	if true {
		toSerialize["tasks"] = o.Tasks
	}
	return json.Marshal(toSerialize)
}

type NullableMoveTasksRequestBody struct {
	value *MoveTasksRequestBody
	isSet bool
}

func (v NullableMoveTasksRequestBody) Get() *MoveTasksRequestBody {
	return v.value
}

func (v *NullableMoveTasksRequestBody) Set(val *MoveTasksRequestBody) {
	v.value = val
	v.isSet = true
}

func (v NullableMoveTasksRequestBody) IsSet() bool {
	return v.isSet
}

func (v *NullableMoveTasksRequestBody) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMoveTasksRequestBody(val *MoveTasksRequestBody) *NullableMoveTasksRequestBody {
	return &NullableMoveTasksRequestBody{value: val, isSet: true}
}

func (v NullableMoveTasksRequestBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMoveTasksRequestBody) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


