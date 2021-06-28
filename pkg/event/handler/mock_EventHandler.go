// Code generated by mockery v2.7.4. DO NOT EDIT.

package handler

import (
	context "context"

	event "github.com/peteqproj/peteq/pkg/event"
	logger "github.com/peteqproj/peteq/pkg/logger"

	mock "github.com/stretchr/testify/mock"
)

// MockEventHandler is an autogenerated mock type for the EventHandler type
type MockEventHandler struct {
	mock.Mock
}

// Handle provides a mock function with given fields: ctx, ev, _a2
func (_m *MockEventHandler) Handle(ctx context.Context, ev event.Event, _a2 logger.Logger) error {
	ret := _m.Called(ctx, ev, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, event.Event, logger.Logger) error); ok {
		r0 = rf(ctx, ev, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *MockEventHandler) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
