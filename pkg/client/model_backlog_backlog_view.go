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

// BacklogBacklogView struct for BacklogBacklogView
type BacklogBacklogView struct {
	Lists    *[]BacklogBacklogTaskList    `json:"lists,omitempty"`
	Projects *[]BacklogBacklogTaskProject `json:"projects,omitempty"`
	Tasks    *[]BacklogBacklogTask        `json:"tasks,omitempty"`
}

// NewBacklogBacklogView instantiates a new BacklogBacklogView object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBacklogBacklogView() *BacklogBacklogView {
	this := BacklogBacklogView{}
	return &this
}

// NewBacklogBacklogViewWithDefaults instantiates a new BacklogBacklogView object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBacklogBacklogViewWithDefaults() *BacklogBacklogView {
	this := BacklogBacklogView{}
	return &this
}

// GetLists returns the Lists field value if set, zero value otherwise.
func (o *BacklogBacklogView) GetLists() []BacklogBacklogTaskList {
	if o == nil || o.Lists == nil {
		var ret []BacklogBacklogTaskList
		return ret
	}
	return *o.Lists
}

// GetListsOk returns a tuple with the Lists field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BacklogBacklogView) GetListsOk() (*[]BacklogBacklogTaskList, bool) {
	if o == nil || o.Lists == nil {
		return nil, false
	}
	return o.Lists, true
}

// HasLists returns a boolean if a field has been set.
func (o *BacklogBacklogView) HasLists() bool {
	if o != nil && o.Lists != nil {
		return true
	}

	return false
}

// SetLists gets a reference to the given []BacklogBacklogTaskList and assigns it to the Lists field.
func (o *BacklogBacklogView) SetLists(v []BacklogBacklogTaskList) {
	o.Lists = &v
}

// GetProjects returns the Projects field value if set, zero value otherwise.
func (o *BacklogBacklogView) GetProjects() []BacklogBacklogTaskProject {
	if o == nil || o.Projects == nil {
		var ret []BacklogBacklogTaskProject
		return ret
	}
	return *o.Projects
}

// GetProjectsOk returns a tuple with the Projects field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BacklogBacklogView) GetProjectsOk() (*[]BacklogBacklogTaskProject, bool) {
	if o == nil || o.Projects == nil {
		return nil, false
	}
	return o.Projects, true
}

// HasProjects returns a boolean if a field has been set.
func (o *BacklogBacklogView) HasProjects() bool {
	if o != nil && o.Projects != nil {
		return true
	}

	return false
}

// SetProjects gets a reference to the given []BacklogBacklogTaskProject and assigns it to the Projects field.
func (o *BacklogBacklogView) SetProjects(v []BacklogBacklogTaskProject) {
	o.Projects = &v
}

// GetTasks returns the Tasks field value if set, zero value otherwise.
func (o *BacklogBacklogView) GetTasks() []BacklogBacklogTask {
	if o == nil || o.Tasks == nil {
		var ret []BacklogBacklogTask
		return ret
	}
	return *o.Tasks
}

// GetTasksOk returns a tuple with the Tasks field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BacklogBacklogView) GetTasksOk() (*[]BacklogBacklogTask, bool) {
	if o == nil || o.Tasks == nil {
		return nil, false
	}
	return o.Tasks, true
}

// HasTasks returns a boolean if a field has been set.
func (o *BacklogBacklogView) HasTasks() bool {
	if o != nil && o.Tasks != nil {
		return true
	}

	return false
}

// SetTasks gets a reference to the given []BacklogBacklogTask and assigns it to the Tasks field.
func (o *BacklogBacklogView) SetTasks(v []BacklogBacklogTask) {
	o.Tasks = &v
}

func (o BacklogBacklogView) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Lists != nil {
		toSerialize["lists"] = o.Lists
	}
	if o.Projects != nil {
		toSerialize["projects"] = o.Projects
	}
	if o.Tasks != nil {
		toSerialize["tasks"] = o.Tasks
	}
	return json.Marshal(toSerialize)
}

type NullableBacklogBacklogView struct {
	value *BacklogBacklogView
	isSet bool
}

func (v NullableBacklogBacklogView) Get() *BacklogBacklogView {
	return v.value
}

func (v *NullableBacklogBacklogView) Set(val *BacklogBacklogView) {
	v.value = val
	v.isSet = true
}

func (v NullableBacklogBacklogView) IsSet() bool {
	return v.isSet
}

func (v *NullableBacklogBacklogView) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBacklogBacklogView(val *BacklogBacklogView) *NullableBacklogBacklogView {
	return &NullableBacklogBacklogView{value: val, isSet: true}
}

func (v NullableBacklogBacklogView) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBacklogBacklogView) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
