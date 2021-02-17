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

// ProjectProjectView struct for ProjectProjectView
type ProjectProjectView struct {
	Id *string `json:"id,omitempty"`
	Metadata *ProjectMetadata `json:"metadata,omitempty"`
	Tasks *[]TaskTask `json:"tasks,omitempty"`
	Type *string `json:"type,omitempty"`
}

// NewProjectProjectView instantiates a new ProjectProjectView object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewProjectProjectView() *ProjectProjectView {
	this := ProjectProjectView{}
	return &this
}

// NewProjectProjectViewWithDefaults instantiates a new ProjectProjectView object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewProjectProjectViewWithDefaults() *ProjectProjectView {
	this := ProjectProjectView{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *ProjectProjectView) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectProjectView) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ProjectProjectView) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ProjectProjectView) SetId(v string) {
	o.Id = &v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *ProjectProjectView) GetMetadata() ProjectMetadata {
	if o == nil || o.Metadata == nil {
		var ret ProjectMetadata
		return ret
	}
	return *o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectProjectView) GetMetadataOk() (*ProjectMetadata, bool) {
	if o == nil || o.Metadata == nil {
		return nil, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *ProjectProjectView) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given ProjectMetadata and assigns it to the Metadata field.
func (o *ProjectProjectView) SetMetadata(v ProjectMetadata) {
	o.Metadata = &v
}

// GetTasks returns the Tasks field value if set, zero value otherwise.
func (o *ProjectProjectView) GetTasks() []TaskTask {
	if o == nil || o.Tasks == nil {
		var ret []TaskTask
		return ret
	}
	return *o.Tasks
}

// GetTasksOk returns a tuple with the Tasks field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectProjectView) GetTasksOk() (*[]TaskTask, bool) {
	if o == nil || o.Tasks == nil {
		return nil, false
	}
	return o.Tasks, true
}

// HasTasks returns a boolean if a field has been set.
func (o *ProjectProjectView) HasTasks() bool {
	if o != nil && o.Tasks != nil {
		return true
	}

	return false
}

// SetTasks gets a reference to the given []TaskTask and assigns it to the Tasks field.
func (o *ProjectProjectView) SetTasks(v []TaskTask) {
	o.Tasks = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *ProjectProjectView) GetType() string {
	if o == nil || o.Type == nil {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectProjectView) GetTypeOk() (*string, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *ProjectProjectView) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *ProjectProjectView) SetType(v string) {
	o.Type = &v
}

func (o ProjectProjectView) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}
	if o.Tasks != nil {
		toSerialize["tasks"] = o.Tasks
	}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}
	return json.Marshal(toSerialize)
}

type NullableProjectProjectView struct {
	value *ProjectProjectView
	isSet bool
}

func (v NullableProjectProjectView) Get() *ProjectProjectView {
	return v.value
}

func (v *NullableProjectProjectView) Set(val *ProjectProjectView) {
	v.value = val
	v.isSet = true
}

func (v NullableProjectProjectView) IsSet() bool {
	return v.isSet
}

func (v *NullableProjectProjectView) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableProjectProjectView(val *ProjectProjectView) *NullableProjectProjectView {
	return &NullableProjectProjectView{value: val, isSet: true}
}

func (v NullableProjectProjectView) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableProjectProjectView) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


