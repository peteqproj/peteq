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

// ProjectMetadata struct for ProjectMetadata
type ProjectMetadata struct {
	Color       *string `json:"color,omitempty"`
	Description *string `json:"description,omitempty"`
	Id          *string `json:"id,omitempty"`
	ImageUrl    *string `json:"imageUrl,omitempty"`
	Name        *string `json:"name,omitempty"`
}

// NewProjectMetadata instantiates a new ProjectMetadata object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewProjectMetadata() *ProjectMetadata {
	this := ProjectMetadata{}
	return &this
}

// NewProjectMetadataWithDefaults instantiates a new ProjectMetadata object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewProjectMetadataWithDefaults() *ProjectMetadata {
	this := ProjectMetadata{}
	return &this
}

// GetColor returns the Color field value if set, zero value otherwise.
func (o *ProjectMetadata) GetColor() string {
	if o == nil || o.Color == nil {
		var ret string
		return ret
	}
	return *o.Color
}

// GetColorOk returns a tuple with the Color field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectMetadata) GetColorOk() (*string, bool) {
	if o == nil || o.Color == nil {
		return nil, false
	}
	return o.Color, true
}

// HasColor returns a boolean if a field has been set.
func (o *ProjectMetadata) HasColor() bool {
	if o != nil && o.Color != nil {
		return true
	}

	return false
}

// SetColor gets a reference to the given string and assigns it to the Color field.
func (o *ProjectMetadata) SetColor(v string) {
	o.Color = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *ProjectMetadata) GetDescription() string {
	if o == nil || o.Description == nil {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectMetadata) GetDescriptionOk() (*string, bool) {
	if o == nil || o.Description == nil {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *ProjectMetadata) HasDescription() bool {
	if o != nil && o.Description != nil {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *ProjectMetadata) SetDescription(v string) {
	o.Description = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *ProjectMetadata) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectMetadata) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ProjectMetadata) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ProjectMetadata) SetId(v string) {
	o.Id = &v
}

// GetImageUrl returns the ImageUrl field value if set, zero value otherwise.
func (o *ProjectMetadata) GetImageUrl() string {
	if o == nil || o.ImageUrl == nil {
		var ret string
		return ret
	}
	return *o.ImageUrl
}

// GetImageUrlOk returns a tuple with the ImageUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectMetadata) GetImageUrlOk() (*string, bool) {
	if o == nil || o.ImageUrl == nil {
		return nil, false
	}
	return o.ImageUrl, true
}

// HasImageUrl returns a boolean if a field has been set.
func (o *ProjectMetadata) HasImageUrl() bool {
	if o != nil && o.ImageUrl != nil {
		return true
	}

	return false
}

// SetImageUrl gets a reference to the given string and assigns it to the ImageUrl field.
func (o *ProjectMetadata) SetImageUrl(v string) {
	o.ImageUrl = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ProjectMetadata) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectMetadata) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ProjectMetadata) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ProjectMetadata) SetName(v string) {
	o.Name = &v
}

func (o ProjectMetadata) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Color != nil {
		toSerialize["color"] = o.Color
	}
	if o.Description != nil {
		toSerialize["description"] = o.Description
	}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.ImageUrl != nil {
		toSerialize["imageUrl"] = o.ImageUrl
	}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	return json.Marshal(toSerialize)
}

type NullableProjectMetadata struct {
	value *ProjectMetadata
	isSet bool
}

func (v NullableProjectMetadata) Get() *ProjectMetadata {
	return v.value
}

func (v *NullableProjectMetadata) Set(val *ProjectMetadata) {
	v.value = val
	v.isSet = true
}

func (v NullableProjectMetadata) IsSet() bool {
	return v.isSet
}

func (v *NullableProjectMetadata) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableProjectMetadata(val *ProjectMetadata) *NullableProjectMetadata {
	return &NullableProjectMetadata{value: val, isSet: true}
}

func (v NullableProjectMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableProjectMetadata) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
