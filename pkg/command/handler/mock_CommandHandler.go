// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package handler

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockCommandHandler is an autogenerated mock type for the CommandHandler type
type MockCommandHandler struct {
	mock.Mock
}

// Handle provides a mock function with given fields: ctx, arguments
func (_m *MockCommandHandler) Handle(ctx context.Context, arguments interface{}) error {
	ret := _m.Called(ctx, arguments)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, arguments)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}