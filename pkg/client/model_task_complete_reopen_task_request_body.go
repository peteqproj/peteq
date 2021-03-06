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

// TaskCompleteReopenTaskRequestBody struct for TaskCompleteReopenTaskRequestBody
type TaskCompleteReopenTaskRequestBody struct {
	Task *string `json:"task,omitempty"`
}

// NewTaskCompleteReopenTaskRequestBody instantiates a new TaskCompleteReopenTaskRequestBody object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTaskCompleteReopenTaskRequestBody() *TaskCompleteReopenTaskRequestBody {
	this := TaskCompleteReopenTaskRequestBody{}
	return &this
}

// NewTaskCompleteReopenTaskRequestBodyWithDefaults instantiates a new TaskCompleteReopenTaskRequestBody object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTaskCompleteReopenTaskRequestBodyWithDefaults() *TaskCompleteReopenTaskRequestBody {
	this := TaskCompleteReopenTaskRequestBody{}
	return &this
}

// GetTask returns the Task field value if set, zero value otherwise.
func (o *TaskCompleteReopenTaskRequestBody) GetTask() string {
	if o == nil || o.Task == nil {
		var ret string
		return ret
	}
	return *o.Task
}

// GetTaskOk returns a tuple with the Task field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskCompleteReopenTaskRequestBody) GetTaskOk() (*string, bool) {
	if o == nil || o.Task == nil {
		return nil, false
	}
	return o.Task, true
}

// HasTask returns a boolean if a field has been set.
func (o *TaskCompleteReopenTaskRequestBody) HasTask() bool {
	if o != nil && o.Task != nil {
		return true
	}

	return false
}

// SetTask gets a reference to the given string and assigns it to the Task field.
func (o *TaskCompleteReopenTaskRequestBody) SetTask(v string) {
	o.Task = &v
}

func (o TaskCompleteReopenTaskRequestBody) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Task != nil {
		toSerialize["task"] = o.Task
	}
	return json.Marshal(toSerialize)
}

type NullableTaskCompleteReopenTaskRequestBody struct {
	value *TaskCompleteReopenTaskRequestBody
	isSet bool
}

func (v NullableTaskCompleteReopenTaskRequestBody) Get() *TaskCompleteReopenTaskRequestBody {
	return v.value
}

func (v *NullableTaskCompleteReopenTaskRequestBody) Set(val *TaskCompleteReopenTaskRequestBody) {
	v.value = val
	v.isSet = true
}

func (v NullableTaskCompleteReopenTaskRequestBody) IsSet() bool {
	return v.isSet
}

func (v *NullableTaskCompleteReopenTaskRequestBody) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTaskCompleteReopenTaskRequestBody(val *TaskCompleteReopenTaskRequestBody) *NullableTaskCompleteReopenTaskRequestBody {
	return &NullableTaskCompleteReopenTaskRequestBody{value: val, isSet: true}
}

func (v NullableTaskCompleteReopenTaskRequestBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTaskCompleteReopenTaskRequestBody) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
