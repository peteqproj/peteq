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

// HomeHomeTask struct for HomeHomeTask
type HomeHomeTask struct {
	Id *string `json:"id,omitempty"`
	Metadata TaskMetadata `json:"metadata"`
	Project *ProjectProject `json:"project,omitempty"`
	Spec TaskSpec `json:"spec"`
	Status *TaskStatus `json:"status,omitempty"`
	Type *string `json:"type,omitempty"`
}

// NewHomeHomeTask instantiates a new HomeHomeTask object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewHomeHomeTask(metadata TaskMetadata, spec TaskSpec, ) *HomeHomeTask {
	this := HomeHomeTask{}
	this.Metadata = metadata
	this.Spec = spec
	return &this
}

// NewHomeHomeTaskWithDefaults instantiates a new HomeHomeTask object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewHomeHomeTaskWithDefaults() *HomeHomeTask {
	this := HomeHomeTask{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *HomeHomeTask) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HomeHomeTask) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *HomeHomeTask) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *HomeHomeTask) SetId(v string) {
	o.Id = &v
}

// GetMetadata returns the Metadata field value
func (o *HomeHomeTask) GetMetadata() TaskMetadata {
	if o == nil  {
		var ret TaskMetadata
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *HomeHomeTask) GetMetadataOk() (*TaskMetadata, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *HomeHomeTask) SetMetadata(v TaskMetadata) {
	o.Metadata = v
}

// GetProject returns the Project field value if set, zero value otherwise.
func (o *HomeHomeTask) GetProject() ProjectProject {
	if o == nil || o.Project == nil {
		var ret ProjectProject
		return ret
	}
	return *o.Project
}

// GetProjectOk returns a tuple with the Project field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HomeHomeTask) GetProjectOk() (*ProjectProject, bool) {
	if o == nil || o.Project == nil {
		return nil, false
	}
	return o.Project, true
}

// HasProject returns a boolean if a field has been set.
func (o *HomeHomeTask) HasProject() bool {
	if o != nil && o.Project != nil {
		return true
	}

	return false
}

// SetProject gets a reference to the given ProjectProject and assigns it to the Project field.
func (o *HomeHomeTask) SetProject(v ProjectProject) {
	o.Project = &v
}

// GetSpec returns the Spec field value
func (o *HomeHomeTask) GetSpec() TaskSpec {
	if o == nil  {
		var ret TaskSpec
		return ret
	}

	return o.Spec
}

// GetSpecOk returns a tuple with the Spec field value
// and a boolean to check if the value has been set.
func (o *HomeHomeTask) GetSpecOk() (*TaskSpec, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Spec, true
}

// SetSpec sets field value
func (o *HomeHomeTask) SetSpec(v TaskSpec) {
	o.Spec = v
}

// GetStatus returns the Status field value if set, zero value otherwise.
func (o *HomeHomeTask) GetStatus() TaskStatus {
	if o == nil || o.Status == nil {
		var ret TaskStatus
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HomeHomeTask) GetStatusOk() (*TaskStatus, bool) {
	if o == nil || o.Status == nil {
		return nil, false
	}
	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *HomeHomeTask) HasStatus() bool {
	if o != nil && o.Status != nil {
		return true
	}

	return false
}

// SetStatus gets a reference to the given TaskStatus and assigns it to the Status field.
func (o *HomeHomeTask) SetStatus(v TaskStatus) {
	o.Status = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *HomeHomeTask) GetType() string {
	if o == nil || o.Type == nil {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HomeHomeTask) GetTypeOk() (*string, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *HomeHomeTask) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *HomeHomeTask) SetType(v string) {
	o.Type = &v
}

func (o HomeHomeTask) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if true {
		toSerialize["metadata"] = o.Metadata
	}
	if o.Project != nil {
		toSerialize["project"] = o.Project
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

type NullableHomeHomeTask struct {
	value *HomeHomeTask
	isSet bool
}

func (v NullableHomeHomeTask) Get() *HomeHomeTask {
	return v.value
}

func (v *NullableHomeHomeTask) Set(val *HomeHomeTask) {
	v.value = val
	v.isSet = true
}

func (v NullableHomeHomeTask) IsSet() bool {
	return v.isSet
}

func (v *NullableHomeHomeTask) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableHomeHomeTask(val *HomeHomeTask) *NullableHomeHomeTask {
	return &NullableHomeHomeTask{value: val, isSet: true}
}

func (v NullableHomeHomeTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableHomeHomeTask) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

