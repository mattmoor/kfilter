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

package kfilter

import (
	"testing"

	"github.com/knative/pkg/controller"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/knative/serving/pkg/reconciler"
	. "github.com/knative/serving/pkg/reconciler/testing"
	v1alpha1testing "github.com/knative/serving/pkg/reconciler/v1alpha1/testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	kfv1alpha1 "github.com/mattmoor/kfilter/pkg/apis/kfilter/v1alpha1"
	clientset "github.com/mattmoor/kfilter/pkg/client/clientset/versioned"
	"github.com/mattmoor/kfilter/pkg/reconciler/kfilter/resources"
	. "github.com/mattmoor/kfilter/pkg/reconciler/testing"
)

const filterImage = "filter:image"

// This is heavily based on the way the OpenShift Ingress controller tests its reconciliation method.
func TestReconcile(t *testing.T) {
	table := TableTest{{
		Name: "bad workqueue key",
		Key:  "too/many/parts",
	}, {
		Name: "key not found",
		Key:  "foo/not-found",
	}, {
		Name: "create knative service",
		Key:  "foo/bar",
		Objects: []runtime.Object{
			kf("bar", "foo"),
		},
		WantCreates: []metav1.Object{
			svc(kf("bar", "foo")),
		},
	}}

	// TODO(mattmoor): Correct the Knative Service
	// TODO(mattmoor): Propagate Knative Service Status

	table.Test(t, MakeFactory(func(listers *Listers, opt reconciler.Options,
		kfClient clientset.Interface) controller.Reconciler {
		return &Reconciler{
			Base:             reconciler.NewBase(opt, controllerAgentName),
			kfilterclientset: kfClient,
			serviceLister:    listers.GetServiceLister(),
			filterLister:     listers.GetFilterLister(),
			filterImage:      filterImage,
		}
	}))
}

func kf(name, namespace string, opts ...FilterOption) *kfv1alpha1.Filter {
	kf := &kfv1alpha1.Filter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	for _, opt := range opts {
		opt(kf)
	}

	return kf
}

func svc(kf *kfv1alpha1.Filter, opts ...v1alpha1testing.ServiceOption) *v1alpha1.Service {
	svc := resources.MakeKService(kf, filterImage)

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}
