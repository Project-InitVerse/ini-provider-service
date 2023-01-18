/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	v1 "providerService/src/ubicpkg/api/ubicnet/v1"
	scheme "providerService/src/ubicpkg/client/clientset/versioned/scheme"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// InventoryRequestsGetter has a method to return a InventoryRequestInterface.
// A group's client should implement this interface.
type InventoryRequestsGetter interface {
	InventoryRequests() InventoryRequestInterface
}

// InventoryRequestInterface has methods to work with InventoryRequest resources.
type InventoryRequestInterface interface {
	Create(ctx context.Context, inventoryRequest *v1.InventoryRequest, opts metav1.CreateOptions) (*v1.InventoryRequest, error)
	Update(ctx context.Context, inventoryRequest *v1.InventoryRequest, opts metav1.UpdateOptions) (*v1.InventoryRequest, error)
	UpdateStatus(ctx context.Context, inventoryRequest *v1.InventoryRequest, opts metav1.UpdateOptions) (*v1.InventoryRequest, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.InventoryRequest, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.InventoryRequestList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.InventoryRequest, err error)
	InventoryRequestExpansion
}

// inventoryRequests implements InventoryRequestInterface
type inventoryRequests struct {
	client rest.Interface
}

// newInventoryRequests returns a InventoryRequests
func newInventoryRequests(c *UbicnetV1Client) *inventoryRequests {
	return &inventoryRequests{
		client: c.RESTClient(),
	}
}

// Get takes name of the inventoryRequest, and returns the corresponding inventoryRequest object, and an error if there is any.
func (c *inventoryRequests) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.InventoryRequest, err error) {
	result = &v1.InventoryRequest{}
	err = c.client.Get().
		Resource("inventoryrequests").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of InventoryRequests that match those selectors.
func (c *inventoryRequests) List(ctx context.Context, opts metav1.ListOptions) (result *v1.InventoryRequestList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.InventoryRequestList{}
	err = c.client.Get().
		Resource("inventoryrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested inventoryRequests.
func (c *inventoryRequests) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("inventoryrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a inventoryRequest and creates it.  Returns the server's representation of the inventoryRequest, and an error, if there is any.
func (c *inventoryRequests) Create(ctx context.Context, inventoryRequest *v1.InventoryRequest, opts metav1.CreateOptions) (result *v1.InventoryRequest, err error) {
	result = &v1.InventoryRequest{}
	err = c.client.Post().
		Resource("inventoryrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(inventoryRequest).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a inventoryRequest and updates it. Returns the server's representation of the inventoryRequest, and an error, if there is any.
func (c *inventoryRequests) Update(ctx context.Context, inventoryRequest *v1.InventoryRequest, opts metav1.UpdateOptions) (result *v1.InventoryRequest, err error) {
	result = &v1.InventoryRequest{}
	err = c.client.Put().
		Resource("inventoryrequests").
		Name(inventoryRequest.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(inventoryRequest).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *inventoryRequests) UpdateStatus(ctx context.Context, inventoryRequest *v1.InventoryRequest, opts metav1.UpdateOptions) (result *v1.InventoryRequest, err error) {
	result = &v1.InventoryRequest{}
	err = c.client.Put().
		Resource("inventoryrequests").
		Name(inventoryRequest.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(inventoryRequest).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the inventoryRequest and deletes it. Returns an error if one occurs.
func (c *inventoryRequests) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("inventoryrequests").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *inventoryRequests) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("inventoryrequests").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched inventoryRequest.
func (c *inventoryRequests) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.InventoryRequest, err error) {
	result = &v1.InventoryRequest{}
	err = c.client.Patch(pt).
		Resource("inventoryrequests").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
