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

// TaskTask struct for TaskTask
type TaskTask struct {
	Id *string `json:"id,omitempty"`
	Metadata TaskMetadata `json:"metadata"`
	Spec TaskSpec `json:"spec"`
	Status *TaskStatus `json:"status,omitempty"`
	Type *string `json:"type,omitempty"`
}

// NewTaskTask instantiates a new TaskTask object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTaskTask(metadata TaskMetadata, spec TaskSpec, ) *TaskTask {
	this := TaskTask{}
	this.Metadata = metadata
	this.Spec = spec
	return &this
}

// NewTaskTaskWithDefaults instantiates a new TaskTask object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTaskTaskWithDefaults() *TaskTask {
	this := TaskTask{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *TaskTask) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskTask) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *TaskTask) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *TaskTask) SetId(v string) {
	o.Id = &v
}

// GetMetadata returns the Metadata field value
func (o *TaskTask) GetMetadata() TaskMetadata {
	if o == nil  {
		var ret TaskMetadata
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *TaskTask) GetMetadataOk() (*TaskMetadata, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *TaskTask) SetMetadata(v TaskMetadata) {
	o.Metadata = v
}

// GetSpec returns the Spec field value
func (o *TaskTask) GetSpec() TaskSpec {
	if o == nil  {
		var ret TaskSpec
		return ret
	}

	return o.Spec
}

// GetSpecOk returns a tuple with the Spec field value
// and a boolean to check if the value has been set.
func (o *TaskTask) GetSpecOk() (*TaskSpec, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Spec, true
}

// SetSpec sets field value
func (o *TaskTask) SetSpec(v TaskSpec) {
	o.Spec = v
}

// GetStatus returns the Status field value if set, zero value otherwise.
func (o *TaskTask) GetStatus() TaskStatus {
	if o == nil || o.Status == nil {
		var ret TaskStatus
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskTask) GetStatusOk() (*TaskStatus, bool) {
	if o == nil || o.Status == nil {
		return nil, false
	}
	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *TaskTask) HasStatus() bool {
	if o != nil && o.Status != nil {
		return true
	}

	return false
}

// SetStatus gets a reference to the given TaskStatus and assigns it to the Status field.
func (o *TaskTask) SetStatus(v TaskStatus) {
	o.Status = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *TaskTask) GetType() string {
	if o == nil || o.Type == nil {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskTask) GetTypeOk() (*string, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *TaskTask) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *TaskTask) SetType(v string) {
	o.Type = &v
}

func (o TaskTask) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if true {
		toSerialize["metadata"] = o.Metadata
	}
	if true {
		toSerialize["spec"] = o.Spec
	}
	if o.Status != nil {
		toSerialize["status"] = o.Status
	}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}
	return json.Marshal(toSerialize)
}

type NullableTaskTask struct {
	value *TaskTask
	isSet bool
}

func (v NullableTaskTask) Get() *TaskTask {
	return v.value
}

func (v *NullableTaskTask) Set(val *TaskTask) {
	v.value = val
	v.isSet = true
}

func (v NullableTaskTask) IsSet() bool {
	return v.isSet
}

func (v *NullableTaskTask) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTaskTask(val *TaskTask) *NullableTaskTask {
	return &NullableTaskTask{value: val, isSet: true}
}

func (v NullableTaskTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTaskTask) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

