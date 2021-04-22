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

// ProjectProject struct for ProjectProject
type ProjectProject struct {
	Metadata *ProjectMetadata `json:"metadata,omitempty"`
	Spec     *ProjectSpec     `json:"spec,omitempty"`
}

// NewProjectProject instantiates a new ProjectProject object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewProjectProject() *ProjectProject {
	this := ProjectProject{}
	return &this
}

// NewProjectProjectWithDefaults instantiates a new ProjectProject object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewProjectProjectWithDefaults() *ProjectProject {
	this := ProjectProject{}
	return &this
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *ProjectProject) GetMetadata() ProjectMetadata {
	if o == nil || o.Metadata == nil {
		var ret ProjectMetadata
		return ret
	}
	return *o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectProject) GetMetadataOk() (*ProjectMetadata, bool) {
	if o == nil || o.Metadata == nil {
		return nil, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *ProjectProject) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given ProjectMetadata and assigns it to the Metadata field.
func (o *ProjectProject) SetMetadata(v ProjectMetadata) {
	o.Metadata = &v
}

// GetSpec returns the Spec field value if set, zero value otherwise.
func (o *ProjectProject) GetSpec() ProjectSpec {
	if o == nil || o.Spec == nil {
		var ret ProjectSpec
		return ret
	}
	return *o.Spec
}

// GetSpecOk returns a tuple with the Spec field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectProject) GetSpecOk() (*ProjectSpec, bool) {
	if o == nil || o.Spec == nil {
		return nil, false
	}
	return o.Spec, true
}

// HasSpec returns a boolean if a field has been set.
func (o *ProjectProject) HasSpec() bool {
	if o != nil && o.Spec != nil {
		return true
	}

	return false
}

// SetSpec gets a reference to the given ProjectSpec and assigns it to the Spec field.
func (o *ProjectProject) SetSpec(v ProjectSpec) {
	o.Spec = &v
}

func (o ProjectProject) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}
	if o.Spec != nil {
		toSerialize["spec"] = o.Spec
	}
	return json.Marshal(toSerialize)
}

type NullableProjectProject struct {
	value *ProjectProject
	isSet bool
}

func (v NullableProjectProject) Get() *ProjectProject {
	return v.value
}

func (v *NullableProjectProject) Set(val *ProjectProject) {
	v.value = val
	v.isSet = true
}

func (v NullableProjectProject) IsSet() bool {
	return v.isSet
}

func (v *NullableProjectProject) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableProjectProject(val *ProjectProject) *NullableProjectProject {
	return &NullableProjectProject{value: val, isSet: true}
}

func (v NullableProjectProject) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableProjectProject) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
