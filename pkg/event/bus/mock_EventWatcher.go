// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package bus

import (
	handler "github.com/peteqproj/peteq/pkg/event/handler"
	mock "github.com/stretchr/testify/mock"
)

// MockEventWatcher is an autogenerated mock type for the EventWatcher type
type MockEventWatcher struct {
	mock.Mock
}

// Subscribe provides a mock function with given fields: name, _a1
func (_m *MockEventWatcher) Subscribe(name string, _a1 handler.EventHandler) {
	_m.Called(name, _a1)
}