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

package testing

import (
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	fakeservingclientset "github.com/knative/serving/pkg/client/clientset/versioned/fake"
	servinglisters "github.com/knative/serving/pkg/client/listers/serving/v1alpha1"
	"github.com/knative/serving/pkg/reconciler/testing"
	"k8s.io/apimachinery/pkg/runtime"
	fakekubeclientset "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"

	kfv1alpha1 "github.com/mattmoor/kfilter/pkg/apis/kfilter/v1alpha1"
	fakekfclientset "github.com/mattmoor/kfilter/pkg/client/clientset/versioned/fake"
	kflisters "github.com/mattmoor/kfilter/pkg/client/listers/kfilter/v1alpha1"
)

var clientSetSchemes = []func(*runtime.Scheme){
	fakekubeclientset.AddToScheme,
	fakeservingclientset.AddToScheme,
	fakekfclientset.AddToScheme,
}

type Listers struct {
	sorter testing.ObjectSorter
}

func NewListers(objs []runtime.Object) Listers {
	scheme := runtime.NewScheme()

	for _, addTo := range clientSetSchemes {
		addTo(scheme)
	}

	ls := Listers{
		sorter: testing.NewObjectSorter(scheme),
	}

	ls.sorter.AddObjects(objs...)

	return ls
}

func (l *Listers) indexerFor(obj runtime.Object) cache.Indexer {
	return l.sorter.IndexerForObjectType(obj)
}

func (l *Listers) GetKubeObjects() []runtime.Object {
	return l.sorter.ObjectsForSchemeFunc(fakekubeclientset.AddToScheme)
}

func (l *Listers) GetKFilterObjects() []runtime.Object {
	return l.sorter.ObjectsForSchemeFunc(fakekfclientset.AddToScheme)
}

func (l *Listers) GetServingObjects() []runtime.Object {
	return l.sorter.ObjectsForSchemeFunc(fakeservingclientset.AddToScheme)
}

func (l *Listers) GetServiceLister() servinglisters.ServiceLister {
	return servinglisters.NewServiceLister(l.indexerFor(&v1alpha1.Service{}))
}

func (l *Listers) GetFilterLister() kflisters.FilterLister {
	return kflisters.NewFilterLister(l.indexerFor(&kfv1alpha1.Filter{}))
}
