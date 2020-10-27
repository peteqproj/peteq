// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package server

import (
	api "github.com/peteqproj/peteq/pkg/api"
	mock "github.com/stretchr/testify/mock"

	socketio "github.com/googollee/go-socket.io"
)

// MockServer is an autogenerated mock type for the Server type
type MockServer struct {
	mock.Mock
}

// AddResource provides a mock function with given fields: r
func (_m *MockServer) AddResource(r api.Resource) error {
	ret := _m.Called(r)

	var r0 error
	if rf, ok := ret.Get(0).(func(api.Resource) error); ok {
		r0 = rf(r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddWS provides a mock function with given fields: ws
func (_m *MockServer) AddWS(ws *socketio.Server) error {
	ret := _m.Called(ws)

	var r0 error
	if rf, ok := ret.Get(0).(func(*socketio.Server) error); ok {
		r0 = rf(ws)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields:
func (_m *MockServer) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}