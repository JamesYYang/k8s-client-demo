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

package v1alpha1

import (
	"context"
	v1alpha1 "k8s-client-demo/crd-controller/apis/emptyapp/v1alpha1"
	scheme "k8s-client-demo/crd-controller/client/clientset/versioned/scheme"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// EmptyAppsGetter has a method to return a EmptyAppInterface.
// A group's client should implement this interface.
type EmptyAppsGetter interface {
	EmptyApps(namespace string) EmptyAppInterface
}

// EmptyAppInterface has methods to work with EmptyApp resources.
type EmptyAppInterface interface {
	Create(ctx context.Context, emptyApp *v1alpha1.EmptyApp, opts v1.CreateOptions) (*v1alpha1.EmptyApp, error)
	Update(ctx context.Context, emptyApp *v1alpha1.EmptyApp, opts v1.UpdateOptions) (*v1alpha1.EmptyApp, error)
	UpdateStatus(ctx context.Context, emptyApp *v1alpha1.EmptyApp, opts v1.UpdateOptions) (*v1alpha1.EmptyApp, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.EmptyApp, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.EmptyAppList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.EmptyApp, err error)
	EmptyAppExpansion
}

// emptyApps implements EmptyAppInterface
type emptyApps struct {
	client rest.Interface
	ns     string
}

// newEmptyApps returns a EmptyApps
func newEmptyApps(c *CrdV1alpha1Client, namespace string) *emptyApps {
	return &emptyApps{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the emptyApp, and returns the corresponding emptyApp object, and an error if there is any.
func (c *emptyApps) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.EmptyApp, err error) {
	result = &v1alpha1.EmptyApp{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("emptyapps").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of EmptyApps that match those selectors.
func (c *emptyApps) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.EmptyAppList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.EmptyAppList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("emptyapps").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested emptyApps.
func (c *emptyApps) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("emptyapps").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a emptyApp and creates it.  Returns the server's representation of the emptyApp, and an error, if there is any.
func (c *emptyApps) Create(ctx context.Context, emptyApp *v1alpha1.EmptyApp, opts v1.CreateOptions) (result *v1alpha1.EmptyApp, err error) {
	result = &v1alpha1.EmptyApp{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("emptyapps").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(emptyApp).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a emptyApp and updates it. Returns the server's representation of the emptyApp, and an error, if there is any.
func (c *emptyApps) Update(ctx context.Context, emptyApp *v1alpha1.EmptyApp, opts v1.UpdateOptions) (result *v1alpha1.EmptyApp, err error) {
	result = &v1alpha1.EmptyApp{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("emptyapps").
		Name(emptyApp.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(emptyApp).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *emptyApps) UpdateStatus(ctx context.Context, emptyApp *v1alpha1.EmptyApp, opts v1.UpdateOptions) (result *v1alpha1.EmptyApp, err error) {
	result = &v1alpha1.EmptyApp{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("emptyapps").
		Name(emptyApp.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(emptyApp).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the emptyApp and deletes it. Returns an error if one occurs.
func (c *emptyApps) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("emptyapps").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *emptyApps) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("emptyapps").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched emptyApp.
func (c *emptyApps) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.EmptyApp, err error) {
	result = &v1alpha1.EmptyApp{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("emptyapps").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
