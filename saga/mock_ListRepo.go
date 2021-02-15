// Code generated by mockery v2.5.1. DO NOT EDIT.

package saga

import (
	list "github.com/peteqproj/peteq/domain/list"
	mock "github.com/stretchr/testify/mock"
)

// MockListRepo is an autogenerated mock type for the ListRepo type
type MockListRepo struct {
	mock.Mock
}

// GetListByName provides a mock function with given fields: userID, name
func (_m *MockListRepo) GetListByName(userID string, name string) (list.List, error) {
	ret := _m.Called(userID, name)

	var r0 list.List
	if rf, ok := ret.Get(0).(func(string, string) list.List); ok {
		r0 = rf(userID, name)
	} else {
		r0 = ret.Get(0).(list.List)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(userID, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
