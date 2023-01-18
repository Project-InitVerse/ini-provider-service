// Code generated by mockery 2.12.1. DO NOT EDIT.

package kubernetes_mocks

import (
	context "context"

	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mock "github.com/stretchr/testify/mock"

	testing "testing"

	types "k8s.io/apimachinery/pkg/types"

	v1 "k8s.io/client-go/applyconfigurations/discovery/v1"

	watch "k8s.io/apimachinery/pkg/watch"
)

// EndpointSliceInterface is an autogenerated mock type for the EndpointSliceInterface type
type EndpointSliceInterface struct {
	mock.Mock
}

// Apply provides a mock function with given fields: ctx, endpointSlice, opts
func (_m *EndpointSliceInterface) Apply(ctx context.Context, endpointSlice *v1.EndpointSliceApplyConfiguration, opts metav1.ApplyOptions) (*discoveryv1.EndpointSlice, error) {
	ret := _m.Called(ctx, endpointSlice, opts)

	var r0 *discoveryv1.EndpointSlice
	if rf, ok := ret.Get(0).(func(context.Context, *v1.EndpointSliceApplyConfiguration, metav1.ApplyOptions) *discoveryv1.EndpointSlice); ok {
		r0 = rf(ctx, endpointSlice, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*discoveryv1.EndpointSlice)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *v1.EndpointSliceApplyConfiguration, metav1.ApplyOptions) error); ok {
		r1 = rf(ctx, endpointSlice, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, endpointSlice, opts
func (_m *EndpointSliceInterface) Create(ctx context.Context, endpointSlice *discoveryv1.EndpointSlice, opts metav1.CreateOptions) (*discoveryv1.EndpointSlice, error) {
	ret := _m.Called(ctx, endpointSlice, opts)

	var r0 *discoveryv1.EndpointSlice
	if rf, ok := ret.Get(0).(func(context.Context, *discoveryv1.EndpointSlice, metav1.CreateOptions) *discoveryv1.EndpointSlice); ok {
		r0 = rf(ctx, endpointSlice, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*discoveryv1.EndpointSlice)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *discoveryv1.EndpointSlice, metav1.CreateOptions) error); ok {
		r1 = rf(ctx, endpointSlice, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, name, opts
func (_m *EndpointSliceInterface) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	ret := _m.Called(ctx, name, opts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.DeleteOptions) error); ok {
		r0 = rf(ctx, name, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteCollection provides a mock function with given fields: ctx, opts, listOpts
func (_m *EndpointSliceInterface) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	ret := _m.Called(ctx, opts, listOpts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.DeleteOptions, metav1.ListOptions) error); ok {
		r0 = rf(ctx, opts, listOpts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, name, opts
func (_m *EndpointSliceInterface) Get(ctx context.Context, name string, opts metav1.GetOptions) (*discoveryv1.EndpointSlice, error) {
	ret := _m.Called(ctx, name, opts)

	var r0 *discoveryv1.EndpointSlice
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) *discoveryv1.EndpointSlice); ok {
		r0 = rf(ctx, name, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*discoveryv1.EndpointSlice)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, metav1.GetOptions) error); ok {
		r1 = rf(ctx, name, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, opts
func (_m *EndpointSliceInterface) List(ctx context.Context, opts metav1.ListOptions) (*discoveryv1.EndpointSliceList, error) {
	ret := _m.Called(ctx, opts)

	var r0 *discoveryv1.EndpointSliceList
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) *discoveryv1.EndpointSliceList); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*discoveryv1.EndpointSliceList)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, metav1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Patch provides a mock function with given fields: ctx, name, pt, data, opts, subresources
func (_m *EndpointSliceInterface) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*discoveryv1.EndpointSlice, error) {
	_va := make([]interface{}, len(subresources))
	for _i := range subresources {
		_va[_i] = subresources[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, name, pt, data, opts)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *discoveryv1.EndpointSlice
	if rf, ok := ret.Get(0).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) *discoveryv1.EndpointSlice); ok {
		r0 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*discoveryv1.EndpointSlice)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) error); ok {
		r1 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, endpointSlice, opts
func (_m *EndpointSliceInterface) Update(ctx context.Context, endpointSlice *discoveryv1.EndpointSlice, opts metav1.UpdateOptions) (*discoveryv1.EndpointSlice, error) {
	ret := _m.Called(ctx, endpointSlice, opts)

	var r0 *discoveryv1.EndpointSlice
	if rf, ok := ret.Get(0).(func(context.Context, *discoveryv1.EndpointSlice, metav1.UpdateOptions) *discoveryv1.EndpointSlice); ok {
		r0 = rf(ctx, endpointSlice, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*discoveryv1.EndpointSlice)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *discoveryv1.EndpointSlice, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, endpointSlice, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Watch provides a mock function with given fields: ctx, opts
func (_m *EndpointSliceInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	ret := _m.Called(ctx, opts)

	var r0 watch.Interface
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) watch.Interface); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(watch.Interface)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, metav1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewEndpointSliceInterface creates a new instance of EndpointSliceInterface. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewEndpointSliceInterface(t testing.TB) *EndpointSliceInterface {
	mock := &EndpointSliceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
