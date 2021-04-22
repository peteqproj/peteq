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

// SensorsSensorViewItem struct for SensorsSensorViewItem
type SensorsSensorViewItem struct {
	Description *string                 `json:"description,omitempty"`
	Id          *string                 `json:"id,omitempty"`
	Name        *string                 `json:"name,omitempty"`
	Spec        *map[string]interface{} `json:"spec,omitempty"`
	Type        *string                 `json:"type,omitempty"`
}

// NewSensorsSensorViewItem instantiates a new SensorsSensorViewItem object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSensorsSensorViewItem() *SensorsSensorViewItem {
	this := SensorsSensorViewItem{}
	return &this
}

// NewSensorsSensorViewItemWithDefaults instantiates a new SensorsSensorViewItem object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSensorsSensorViewItemWithDefaults() *SensorsSensorViewItem {
	this := SensorsSensorViewItem{}
	return &this
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *SensorsSensorViewItem) GetDescription() string {
	if o == nil || o.Description == nil {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SensorsSensorViewItem) GetDescriptionOk() (*string, bool) {
	if o == nil || o.Description == nil {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *SensorsSensorViewItem) HasDescription() bool {
	if o != nil && o.Description != nil {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *SensorsSensorViewItem) SetDescription(v string) {
	o.Description = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *SensorsSensorViewItem) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SensorsSensorViewItem) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *SensorsSensorViewItem) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *SensorsSensorViewItem) SetId(v string) {
	o.Id = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *SensorsSensorViewItem) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SensorsSensorViewItem) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *SensorsSensorViewItem) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *SensorsSensorViewItem) SetName(v string) {
	o.Name = &v
}

// GetSpec returns the Spec field value if set, zero value otherwise.
func (o *SensorsSensorViewItem) GetSpec() map[string]interface{} {
	if o == nil || o.Spec == nil {
		var ret map[string]interface{}
		return ret
	}
	return *o.Spec
}

// GetSpecOk returns a tuple with the Spec field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SensorsSensorViewItem) GetSpecOk() (*map[string]interface{}, bool) {
	if o == nil || o.Spec == nil {
		return nil, false
	}
	return o.Spec, true
}

// HasSpec returns a boolean if a field has been set.
func (o *SensorsSensorViewItem) HasSpec() bool {
	if o != nil && o.Spec != nil {
		return true
	}

	return false
}

// SetSpec gets a reference to the given map[string]interface{} and assigns it to the Spec field.
func (o *SensorsSensorViewItem) SetSpec(v map[string]interface{}) {
	o.Spec = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *SensorsSensorViewItem) GetType() string {
	if o == nil || o.Type == nil {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SensorsSensorViewItem) GetTypeOk() (*string, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *SensorsSensorViewItem) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *SensorsSensorViewItem) SetType(v string) {
	o.Type = &v
}

func (o SensorsSensorViewItem) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Description != nil {
		toSerialize["description"] = o.Description
	}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.Spec != nil {
		toSerialize["spec"] = o.Spec
	}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}
	return json.Marshal(toSerialize)
}

type NullableSensorsSensorViewItem struct {
	value *SensorsSensorViewItem
	isSet bool
}

func (v NullableSensorsSensorViewItem) Get() *SensorsSensorViewItem {
	return v.value
}

func (v *NullableSensorsSensorViewItem) Set(val *SensorsSensorViewItem) {
	v.value = val
	v.isSet = true
}

func (v NullableSensorsSensorViewItem) IsSet() bool {
	return v.isSet
}

func (v *NullableSensorsSensorViewItem) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSensorsSensorViewItem(val *SensorsSensorViewItem) *NullableSensorsSensorViewItem {
	return &NullableSensorsSensorViewItem{value: val, isSet: true}
}

func (v NullableSensorsSensorViewItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSensorsSensorViewItem) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
