// Code generated by mockery 2.12.1. DO NOT EDIT.

package kubernetes_mocks

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// SelfSubjectAccessReviewExpansion is an autogenerated mock type for the SelfSubjectAccessReviewExpansion type
type SelfSubjectAccessReviewExpansion struct {
	mock.Mock
}

// NewSelfSubjectAccessReviewExpansion creates a new instance of SelfSubjectAccessReviewExpansion. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewSelfSubjectAccessReviewExpansion(t testing.TB) *SelfSubjectAccessReviewExpansion {
	mock := &SelfSubjectAccessReviewExpansion{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
