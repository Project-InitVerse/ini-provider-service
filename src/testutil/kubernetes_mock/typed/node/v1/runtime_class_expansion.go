// Code generated by mockery 2.12.1. DO NOT EDIT.

package kubernetes_mocks

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// RuntimeClassExpansion is an autogenerated mock type for the RuntimeClassExpansion type
type RuntimeClassExpansion struct {
	mock.Mock
}

// NewRuntimeClassExpansion creates a new instance of RuntimeClassExpansion. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewRuntimeClassExpansion(t testing.TB) *RuntimeClassExpansion {
	mock := &RuntimeClassExpansion{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
