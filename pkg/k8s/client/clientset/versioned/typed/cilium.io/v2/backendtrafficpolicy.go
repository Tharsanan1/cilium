// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

// Code generated by client-gen. DO NOT EDIT.

package v2

import (
	"context"
	"time"

	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	scheme "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// BackendTrafficPoliciesGetter has a method to return a BackendTrafficPolicyInterface.
// A group's client should implement this interface.
type BackendTrafficPoliciesGetter interface {
	BackendTrafficPolicies(namespace string) BackendTrafficPolicyInterface
}

// BackendTrafficPolicyInterface has methods to work with BackendTrafficPolicy resources.
type BackendTrafficPolicyInterface interface {
	Create(ctx context.Context, backendTrafficPolicy *v2.BackendTrafficPolicy, opts v1.CreateOptions) (*v2.BackendTrafficPolicy, error)
	Update(ctx context.Context, backendTrafficPolicy *v2.BackendTrafficPolicy, opts v1.UpdateOptions) (*v2.BackendTrafficPolicy, error)
	UpdateStatus(ctx context.Context, backendTrafficPolicy *v2.BackendTrafficPolicy, opts v1.UpdateOptions) (*v2.BackendTrafficPolicy, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v2.BackendTrafficPolicy, error)
	List(ctx context.Context, opts v1.ListOptions) (*v2.BackendTrafficPolicyList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2.BackendTrafficPolicy, err error)
	BackendTrafficPolicyExpansion
}

// backendTrafficPolicies implements BackendTrafficPolicyInterface
type backendTrafficPolicies struct {
	client rest.Interface
	ns     string
}

// newBackendTrafficPolicies returns a BackendTrafficPolicies
func newBackendTrafficPolicies(c *CiliumV2Client, namespace string) *backendTrafficPolicies {
	return &backendTrafficPolicies{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the backendTrafficPolicy, and returns the corresponding backendTrafficPolicy object, and an error if there is any.
func (c *backendTrafficPolicies) Get(ctx context.Context, name string, options v1.GetOptions) (result *v2.BackendTrafficPolicy, err error) {
	result = &v2.BackendTrafficPolicy{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("backendtrafficpolicies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BackendTrafficPolicies that match those selectors.
func (c *backendTrafficPolicies) List(ctx context.Context, opts v1.ListOptions) (result *v2.BackendTrafficPolicyList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v2.BackendTrafficPolicyList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("backendtrafficpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested backendTrafficPolicies.
func (c *backendTrafficPolicies) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("backendtrafficpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a backendTrafficPolicy and creates it.  Returns the server's representation of the backendTrafficPolicy, and an error, if there is any.
func (c *backendTrafficPolicies) Create(ctx context.Context, backendTrafficPolicy *v2.BackendTrafficPolicy, opts v1.CreateOptions) (result *v2.BackendTrafficPolicy, err error) {
	result = &v2.BackendTrafficPolicy{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("backendtrafficpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(backendTrafficPolicy).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a backendTrafficPolicy and updates it. Returns the server's representation of the backendTrafficPolicy, and an error, if there is any.
func (c *backendTrafficPolicies) Update(ctx context.Context, backendTrafficPolicy *v2.BackendTrafficPolicy, opts v1.UpdateOptions) (result *v2.BackendTrafficPolicy, err error) {
	result = &v2.BackendTrafficPolicy{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("backendtrafficpolicies").
		Name(backendTrafficPolicy.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(backendTrafficPolicy).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *backendTrafficPolicies) UpdateStatus(ctx context.Context, backendTrafficPolicy *v2.BackendTrafficPolicy, opts v1.UpdateOptions) (result *v2.BackendTrafficPolicy, err error) {
	result = &v2.BackendTrafficPolicy{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("backendtrafficpolicies").
		Name(backendTrafficPolicy.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(backendTrafficPolicy).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the backendTrafficPolicy and deletes it. Returns an error if one occurs.
func (c *backendTrafficPolicies) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("backendtrafficpolicies").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *backendTrafficPolicies) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("backendtrafficpolicies").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched backendTrafficPolicy.
func (c *backendTrafficPolicies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2.BackendTrafficPolicy, err error) {
	result = &v2.BackendTrafficPolicy{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("backendtrafficpolicies").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
