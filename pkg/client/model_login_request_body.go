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

// LoginRequestBody struct for LoginRequestBody
type LoginRequestBody struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

// NewLoginRequestBody instantiates a new LoginRequestBody object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLoginRequestBody(email string, password string) *LoginRequestBody {
	this := LoginRequestBody{}
	this.Email = email
	this.Password = password
	return &this
}

// NewLoginRequestBodyWithDefaults instantiates a new LoginRequestBody object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLoginRequestBodyWithDefaults() *LoginRequestBody {
	this := LoginRequestBody{}
	return &this
}

// GetEmail returns the Email field value
func (o *LoginRequestBody) GetEmail() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Email
}

// GetEmailOk returns a tuple with the Email field value
// and a boolean to check if the value has been set.
func (o *LoginRequestBody) GetEmailOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Email, true
}

// SetEmail sets field value
func (o *LoginRequestBody) SetEmail(v string) {
	o.Email = v
}

// GetPassword returns the Password field value
func (o *LoginRequestBody) GetPassword() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Password
}

// GetPasswordOk returns a tuple with the Password field value
// and a boolean to check if the value has been set.
func (o *LoginRequestBody) GetPasswordOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Password, true
}

// SetPassword sets field value
func (o *LoginRequestBody) SetPassword(v string) {
	o.Password = v
}

func (o LoginRequestBody) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["email"] = o.Email
	}
	if true {
		toSerialize["password"] = o.Password
	}
	return json.Marshal(toSerialize)
}

type NullableLoginRequestBody struct {
	value *LoginRequestBody
	isSet bool
}

func (v NullableLoginRequestBody) Get() *LoginRequestBody {
	return v.value
}

func (v *NullableLoginRequestBody) Set(val *LoginRequestBody) {
	v.value = val
	v.isSet = true
}

func (v NullableLoginRequestBody) IsSet() bool {
	return v.isSet
}

func (v *NullableLoginRequestBody) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLoginRequestBody(val *LoginRequestBody) *NullableLoginRequestBody {
	return &NullableLoginRequestBody{value: val, isSet: true}
}

func (v NullableLoginRequestBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLoginRequestBody) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


