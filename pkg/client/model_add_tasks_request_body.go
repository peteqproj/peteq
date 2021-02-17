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

// AddTasksRequestBody struct for AddTasksRequestBody
type AddTasksRequestBody struct {
	Project string `json:"project"`
	Tasks []string `json:"tasks"`
}

// NewAddTasksRequestBody instantiates a new AddTasksRequestBody object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAddTasksRequestBody(project string, tasks []string, ) *AddTasksRequestBody {
	this := AddTasksRequestBody{}
	this.Project = project
	this.Tasks = tasks
	return &this
}

// NewAddTasksRequestBodyWithDefaults instantiates a new AddTasksRequestBody object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAddTasksRequestBodyWithDefaults() *AddTasksRequestBody {
	this := AddTasksRequestBody{}
	return &this
}

// GetProject returns the Project field value
func (o *AddTasksRequestBody) GetProject() string {
	if o == nil  {
		var ret string
		return ret
	}

	return o.Project
}

// GetProjectOk returns a tuple with the Project field value
// and a boolean to check if the value has been set.
func (o *AddTasksRequestBody) GetProjectOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Project, true
}

// SetProject sets field value
func (o *AddTasksRequestBody) SetProject(v string) {
	o.Project = v
}

// GetTasks returns the Tasks field value
func (o *AddTasksRequestBody) GetTasks() []string {
	if o == nil  {
		var ret []string
		return ret
	}

	return o.Tasks
}

// GetTasksOk returns a tuple with the Tasks field value
// and a boolean to check if the value has been set.
func (o *AddTasksRequestBody) GetTasksOk() (*[]string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Tasks, true
}

// SetTasks sets field value
func (o *AddTasksRequestBody) SetTasks(v []string) {
	o.Tasks = v
}

func (o AddTasksRequestBody) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["project"] = o.Project
	}
	if true {
		toSerialize["tasks"] = o.Tasks
	}
	return json.Marshal(toSerialize)
}

type NullableAddTasksRequestBody struct {
	value *AddTasksRequestBody
	isSet bool
}

func (v NullableAddTasksRequestBody) Get() *AddTasksRequestBody {
	return v.value
}

func (v *NullableAddTasksRequestBody) Set(val *AddTasksRequestBody) {
	v.value = val
	v.isSet = true
}

func (v NullableAddTasksRequestBody) IsSet() bool {
	return v.isSet
}

func (v *NullableAddTasksRequestBody) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAddTasksRequestBody(val *AddTasksRequestBody) *NullableAddTasksRequestBody {
	return &NullableAddTasksRequestBody{value: val, isSet: true}
}

func (v NullableAddTasksRequestBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAddTasksRequestBody) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


