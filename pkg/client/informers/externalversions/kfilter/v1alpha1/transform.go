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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	kfilterv1alpha1 "github.com/mattmoor/kfilter/pkg/apis/kfilter/v1alpha1"
	versioned "github.com/mattmoor/kfilter/pkg/client/clientset/versioned"
	internalinterfaces "github.com/mattmoor/kfilter/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/mattmoor/kfilter/pkg/client/listers/kfilter/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// TransformInformer provides access to a shared informer and lister for
// Transforms.
type TransformInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.TransformLister
}

type transformInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewTransformInformer constructs a new informer for Transform type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewTransformInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredTransformInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredTransformInformer constructs a new informer for Transform type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredTransformInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KfilterV1alpha1().Transforms(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KfilterV1alpha1().Transforms(namespace).Watch(options)
			},
		},
		&kfilterv1alpha1.Transform{},
		resyncPeriod,
		indexers,
	)
}

func (f *transformInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredTransformInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *transformInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&kfilterv1alpha1.Transform{}, f.defaultInformer)
}

func (f *transformInformer) Lister() v1alpha1.TransformLister {
	return v1alpha1.NewTransformLister(f.Informer().GetIndexer())
}
