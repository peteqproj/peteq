// Code generated by mockery v2.7.5. DO NOT EDIT.

package logger

import mock "github.com/stretchr/testify/mock"

// MockLogger is an autogenerated mock type for the Logger type
type MockLogger struct {
	mock.Mock
}

// Fork provides a mock function with given fields: keysAndValues
func (_m *MockLogger) Fork(keysAndValues ...interface{}) Logger {
	var _ca []interface{}
	_ca = append(_ca, keysAndValues...)
	ret := _m.Called(_ca...)

	var r0 Logger
	if rf, ok := ret.Get(0).(func(...interface{}) Logger); ok {
		r0 = rf(keysAndValues...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Logger)
		}
	}

	return r0
}

// Info provides a mock function with given fields: msg, keysAndValues
func (_m *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, keysAndValues...)
	_m.Called(_ca...)
}

// V provides a mock function with given fields: level
func (_m *MockLogger) V(level int) Logger {
	ret := _m.Called(level)

	var r0 Logger
	if rf, ok := ret.Get(0).(func(int) Logger); ok {
		r0 = rf(level)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Logger)
		}
	}

	return r0
}
