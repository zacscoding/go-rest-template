// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// Provider is an autogenerated mock type for the Provider type
type Provider struct {
	mock.Mock
}

// RecordApiCount provides a mock function with given fields: code, method, path
func (_m *Provider) RecordApiCount(code int, method string, path string) {
	_m.Called(code, method, path)
}

// RecordApiLatency provides a mock function with given fields: code, method, path, elapsed
func (_m *Provider) RecordApiLatency(code int, method string, path string, elapsed time.Duration) {
	_m.Called(code, method, path, elapsed)
}

// RecordCache provides a mock function with given fields: key, hit
func (_m *Provider) RecordCache(key string, hit bool) {
	_m.Called(key, hit)
}

type mockConstructorTestingTNewProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewProvider creates a new instance of Provider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProvider(t mockConstructorTestingTNewProvider) *Provider {
	mock := &Provider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
