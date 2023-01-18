// Code generated by mockery 2.12.1. DO NOT EDIT.

package kubernetes_mocks

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"

	v1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
)

// ClusterRoleBindingsGetter is an autogenerated mock type for the ClusterRoleBindingsGetter type
type ClusterRoleBindingsGetter struct {
	mock.Mock
}

// ClusterRoleBindings provides a mock function with given fields:
func (_m *ClusterRoleBindingsGetter) ClusterRoleBindings() v1beta1.ClusterRoleBindingInterface {
	ret := _m.Called()

	var r0 v1beta1.ClusterRoleBindingInterface
	if rf, ok := ret.Get(0).(func() v1beta1.ClusterRoleBindingInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1beta1.ClusterRoleBindingInterface)
		}
	}

	return r0
}

// NewClusterRoleBindingsGetter creates a new instance of ClusterRoleBindingsGetter. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewClusterRoleBindingsGetter(t testing.TB) *ClusterRoleBindingsGetter {
	mock := &ClusterRoleBindingsGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
