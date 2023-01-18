// Code generated by mockery 2.12.1. DO NOT EDIT.

package kubernetes_mocks

import (
	mock "github.com/stretchr/testify/mock"
	rest "k8s.io/client-go/rest"

	testing "testing"

	v1beta1 "k8s.io/client-go/kubernetes/typed/scheduling/v1beta1"
)

// SchedulingV1beta1Interface is an autogenerated mock type for the SchedulingV1beta1Interface type
type SchedulingV1beta1Interface struct {
	mock.Mock
}

// PriorityClasses provides a mock function with given fields:
func (_m *SchedulingV1beta1Interface) PriorityClasses() v1beta1.PriorityClassInterface {
	ret := _m.Called()

	var r0 v1beta1.PriorityClassInterface
	if rf, ok := ret.Get(0).(func() v1beta1.PriorityClassInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1beta1.PriorityClassInterface)
		}
	}

	return r0
}

// RESTClient provides a mock function with given fields:
func (_m *SchedulingV1beta1Interface) RESTClient() rest.Interface {
	ret := _m.Called()

	var r0 rest.Interface
	if rf, ok := ret.Get(0).(func() rest.Interface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(rest.Interface)
		}
	}

	return r0
}

// NewSchedulingV1beta1Interface creates a new instance of SchedulingV1beta1Interface. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewSchedulingV1beta1Interface(t testing.TB) *SchedulingV1beta1Interface {
	mock := &SchedulingV1beta1Interface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
