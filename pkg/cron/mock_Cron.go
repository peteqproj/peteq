// Code generated by mockery v2.7.4. DO NOT EDIT.

package cron

import mock "github.com/stretchr/testify/mock"

// MockCron is an autogenerated mock type for the Cron type
type MockCron struct {
	mock.Mock
}

// AddFunc provides a mock function with given fields: trigger, cronExp
func (_m *MockCron) AddFunc(trigger string, cronExp string) {
	_m.Called(trigger, cronExp)
}

// Start provides a mock function with given fields:
func (_m *MockCron) Start() {
	_m.Called()
}

// Stop provides a mock function with given fields:
func (_m *MockCron) Stop() {
	_m.Called()
}
