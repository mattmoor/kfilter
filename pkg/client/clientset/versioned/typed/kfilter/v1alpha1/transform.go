/*
Copyright 2018 Matt Moore

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
	v1alpha1 "github.com/mattmoor/kfilter/pkg/apis/kfilter/v1alpha1"
	scheme "github.com/mattmoor/kfilter/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// TransformsGetter has a method to return a TransformInterface.
// A group's client should implement this interface.
type TransformsGetter interface {
	Transforms(namespace string) TransformInterface
}

// TransformInterface has methods to work with Transform resources.
type TransformInterface interface {
	Create(*v1alpha1.Transform) (*v1alpha1.Transform, error)
	Update(*v1alpha1.Transform) (*v1alpha1.Transform, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Transform, error)
	List(opts v1.ListOptions) (*v1alpha1.TransformList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Transform, err error)
	TransformExpansion
}

// transforms implements TransformInterface
type transforms struct {
	client rest.Interface
	ns     string
}

// newTransforms returns a Transforms
func newTransforms(c *KfilterV1alpha1Client, namespace string) *transforms {
	return &transforms{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the transform, and returns the corresponding transform object, and an error if there is any.
func (c *transforms) Get(name string, options v1.GetOptions) (result *v1alpha1.Transform, err error) {
	result = &v1alpha1.Transform{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("transforms").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Transforms that match those selectors.
func (c *transforms) List(opts v1.ListOptions) (result *v1alpha1.TransformList, err error) {
	result = &v1alpha1.TransformList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("transforms").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested transforms.
func (c *transforms) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("transforms").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a transform and creates it.  Returns the server's representation of the transform, and an error, if there is any.
func (c *transforms) Create(transform *v1alpha1.Transform) (result *v1alpha1.Transform, err error) {
	result = &v1alpha1.Transform{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("transforms").
		Body(transform).
		Do().
		Into(result)
	return
}

// Update takes the representation of a transform and updates it. Returns the server's representation of the transform, and an error, if there is any.
func (c *transforms) Update(transform *v1alpha1.Transform) (result *v1alpha1.Transform, err error) {
	result = &v1alpha1.Transform{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("transforms").
		Name(transform.Name).
		Body(transform).
		Do().
		Into(result)
	return
}

// Delete takes name of the transform and deletes it. Returns an error if one occurs.
func (c *transforms) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("transforms").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *transforms) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("transforms").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched transform.
func (c *transforms) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Transform, err error) {
	result = &v1alpha1.Transform{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("transforms").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}