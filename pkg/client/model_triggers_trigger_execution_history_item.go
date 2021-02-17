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

// TriggersTriggerExecutionHistoryItem struct for TriggersTriggerExecutionHistoryItem
type TriggersTriggerExecutionHistoryItem struct {
	Manual *bool `json:"manual,omitempty"`
	TriggeredAt *string `json:"triggeredAt,omitempty"`
}

// NewTriggersTriggerExecutionHistoryItem instantiates a new TriggersTriggerExecutionHistoryItem object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTriggersTriggerExecutionHistoryItem() *TriggersTriggerExecutionHistoryItem {
	this := TriggersTriggerExecutionHistoryItem{}
	return &this
}

// NewTriggersTriggerExecutionHistoryItemWithDefaults instantiates a new TriggersTriggerExecutionHistoryItem object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTriggersTriggerExecutionHistoryItemWithDefaults() *TriggersTriggerExecutionHistoryItem {
	this := TriggersTriggerExecutionHistoryItem{}
	return &this
}

// GetManual returns the Manual field value if set, zero value otherwise.
func (o *TriggersTriggerExecutionHistoryItem) GetManual() bool {
	if o == nil || o.Manual == nil {
		var ret bool
		return ret
	}
	return *o.Manual
}

// GetManualOk returns a tuple with the Manual field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TriggersTriggerExecutionHistoryItem) GetManualOk() (*bool, bool) {
	if o == nil || o.Manual == nil {
		return nil, false
	}
	return o.Manual, true
}

// HasManual returns a boolean if a field has been set.
func (o *TriggersTriggerExecutionHistoryItem) HasManual() bool {
	if o != nil && o.Manual != nil {
		return true
	}

	return false
}

// SetManual gets a reference to the given bool and assigns it to the Manual field.
func (o *TriggersTriggerExecutionHistoryItem) SetManual(v bool) {
	o.Manual = &v
}

// GetTriggeredAt returns the TriggeredAt field value if set, zero value otherwise.
func (o *TriggersTriggerExecutionHistoryItem) GetTriggeredAt() string {
	if o == nil || o.TriggeredAt == nil {
		var ret string
		return ret
	}
	return *o.TriggeredAt
}

// GetTriggeredAtOk returns a tuple with the TriggeredAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TriggersTriggerExecutionHistoryItem) GetTriggeredAtOk() (*string, bool) {
	if o == nil || o.TriggeredAt == nil {
		return nil, false
	}
	return o.TriggeredAt, true
}

// HasTriggeredAt returns a boolean if a field has been set.
func (o *TriggersTriggerExecutionHistoryItem) HasTriggeredAt() bool {
	if o != nil && o.TriggeredAt != nil {
		return true
	}

	return false
}

// SetTriggeredAt gets a reference to the given string and assigns it to the TriggeredAt field.
func (o *TriggersTriggerExecutionHistoryItem) SetTriggeredAt(v string) {
	o.TriggeredAt = &v
}

func (o TriggersTriggerExecutionHistoryItem) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Manual != nil {
		toSerialize["manual"] = o.Manual
	}
	if o.TriggeredAt != nil {
		toSerialize["triggeredAt"] = o.TriggeredAt
	}
	return json.Marshal(toSerialize)
}

type NullableTriggersTriggerExecutionHistoryItem struct {
	value *TriggersTriggerExecutionHistoryItem
	isSet bool
}

func (v NullableTriggersTriggerExecutionHistoryItem) Get() *TriggersTriggerExecutionHistoryItem {
	return v.value
}

func (v *NullableTriggersTriggerExecutionHistoryItem) Set(val *TriggersTriggerExecutionHistoryItem) {
	v.value = val
	v.isSet = true
}

func (v NullableTriggersTriggerExecutionHistoryItem) IsSet() bool {
	return v.isSet
}

func (v *NullableTriggersTriggerExecutionHistoryItem) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTriggersTriggerExecutionHistoryItem(val *TriggersTriggerExecutionHistoryItem) *NullableTriggersTriggerExecutionHistoryItem {
	return &NullableTriggersTriggerExecutionHistoryItem{value: val, isSet: true}
}

func (v NullableTriggersTriggerExecutionHistoryItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTriggersTriggerExecutionHistoryItem) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


