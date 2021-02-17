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

// ListList struct for ListList
type ListList struct {
	Id *string `json:"id,omitempty"`
	Metadata *ListMetadata `json:"metadata,omitempty"`
	Tasks *[]string `json:"tasks,omitempty"`
	Type *string `json:"type,omitempty"`
}

// NewListList instantiates a new ListList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListList() *ListList {
	this := ListList{}
	return &this
}

// NewListListWithDefaults instantiates a new ListList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListListWithDefaults() *ListList {
	this := ListList{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *ListList) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListList) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ListList) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ListList) SetId(v string) {
	o.Id = &v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *ListList) GetMetadata() ListMetadata {
	if o == nil || o.Metadata == nil {
		var ret ListMetadata
		return ret
	}
	return *o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListList) GetMetadataOk() (*ListMetadata, bool) {
	if o == nil || o.Metadata == nil {
		return nil, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *ListList) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given ListMetadata and assigns it to the Metadata field.
func (o *ListList) SetMetadata(v ListMetadata) {
	o.Metadata = &v
}

// GetTasks returns the Tasks field value if set, zero value otherwise.
func (o *ListList) GetTasks() []string {
	if o == nil || o.Tasks == nil {
		var ret []string
		return ret
	}
	return *o.Tasks
}

// GetTasksOk returns a tuple with the Tasks field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListList) GetTasksOk() (*[]string, bool) {
	if o == nil || o.Tasks == nil {
		return nil, false
	}
	return o.Tasks, true
}

// HasTasks returns a boolean if a field has been set.
func (o *ListList) HasTasks() bool {
	if o != nil && o.Tasks != nil {
		return true
	}

	return false
}

// SetTasks gets a reference to the given []string and assigns it to the Tasks field.
func (o *ListList) SetTasks(v []string) {
	o.Tasks = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *ListList) GetType() string {
	if o == nil || o.Type == nil {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListList) GetTypeOk() (*string, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *ListList) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *ListList) SetType(v string) {
	o.Type = &v
}

func (o ListList) MarshalJSON() ([]byte, error) {
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

type NullableListList struct {
	value *ListList
	isSet bool
}

func (v NullableListList) Get() *ListList {
	return v.value
}

func (v *NullableListList) Set(val *ListList) {
	v.value = val
	v.isSet = true
}

func (v NullableListList) IsSet() bool {
	return v.isSet
}

func (v *NullableListList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListList(val *ListList) *NullableListList {
	return &NullableListList{value: val, isSet: true}
}

func (v NullableListList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

