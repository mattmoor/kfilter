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

package fake

import (
	v1alpha1 "github.com/mattmoor/kfilter/pkg/apis/kfilter/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeTransforms implements TransformInterface
type FakeTransforms struct {
	Fake *FakeKfilterV1alpha1
	ns   string
}

var transformsResource = schema.GroupVersionResource{Group: "kfilter.mattmoor.io", Version: "v1alpha1", Resource: "transforms"}

var transformsKind = schema.GroupVersionKind{Group: "kfilter.mattmoor.io", Version: "v1alpha1", Kind: "Transform"}

// Get takes name of the transform, and returns the corresponding transform object, and an error if there is any.
func (c *FakeTransforms) Get(name string, options v1.GetOptions) (result *v1alpha1.Transform, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(transformsResource, c.ns, name), &v1alpha1.Transform{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Transform), err
}

// List takes label and field selectors, and returns the list of Transforms that match those selectors.
func (c *FakeTransforms) List(opts v1.ListOptions) (result *v1alpha1.TransformList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(transformsResource, transformsKind, c.ns, opts), &v1alpha1.TransformList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.TransformList{ListMeta: obj.(*v1alpha1.TransformList).ListMeta}
	for _, item := range obj.(*v1alpha1.TransformList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested transforms.
func (c *FakeTransforms) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(transformsResource, c.ns, opts))

}

// Create takes the representation of a transform and creates it.  Returns the server's representation of the transform, and an error, if there is any.
func (c *FakeTransforms) Create(transform *v1alpha1.Transform) (result *v1alpha1.Transform, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(transformsResource, c.ns, transform), &v1alpha1.Transform{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Transform), err
}

// Update takes the representation of a transform and updates it. Returns the server's representation of the transform, and an error, if there is any.
func (c *FakeTransforms) Update(transform *v1alpha1.Transform) (result *v1alpha1.Transform, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(transformsResource, c.ns, transform), &v1alpha1.Transform{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Transform), err
}

// Delete takes name of the transform and deletes it. Returns an error if one occurs.
func (c *FakeTransforms) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(transformsResource, c.ns, name), &v1alpha1.Transform{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeTransforms) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(transformsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.TransformList{})
	return err
}

// Patch applies the patch and returns the patched transform.
func (c *FakeTransforms) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Transform, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(transformsResource, c.ns, name, data, subresources...), &v1alpha1.Transform{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Transform), err
}