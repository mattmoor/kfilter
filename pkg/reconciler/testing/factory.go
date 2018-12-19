/*
Copyright 2018 The Knative Authors.

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
	"testing"

	fakekubeclientset "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/record"

	"github.com/knative/pkg/controller"
	logtesting "github.com/knative/pkg/logging/testing"
	fakeclientset "github.com/knative/serving/pkg/client/clientset/versioned/fake"
	"github.com/knative/serving/pkg/reconciler"
	. "github.com/knative/serving/pkg/reconciler/testing"

	clientset "github.com/mattmoor/kfilter/pkg/client/clientset/versioned"
	fakekfclientset "github.com/mattmoor/kfilter/pkg/client/clientset/versioned/fake"
)

const (
	// maxEventBufferSize is the estimated max number of event notifications that
	// can be buffered during reconciliation.
	maxEventBufferSize = 10
)

// Ctor functions create a k8s controller with given params.
type Ctor func(*Listers, reconciler.Options, clientset.Interface) controller.Reconciler

// MakeFactory creates a reconciler factory with fake clients and controller created by `ctor`.
func MakeFactory(ctor Ctor) Factory {
	return func(t *testing.T, r *TableRow) (controller.Reconciler, ActionRecorderList, EventList) {
		ls := NewListers(r.Objects)

		kubeClient := fakekubeclientset.NewSimpleClientset(ls.GetKubeObjects()...)
		servingclient := fakeclientset.NewSimpleClientset(ls.GetServingObjects()...)
		kfClient := fakekfclientset.NewSimpleClientset(ls.GetKFilterObjects()...)
		eventRecorder := record.NewFakeRecorder(maxEventBufferSize)

		// Set up our Controller from the fakes.
		c := ctor(&ls, reconciler.Options{
			KubeClientSet:    kubeClient,
			ServingClientSet: servingclient,
			Recorder:         eventRecorder,
			Logger:           logtesting.TestLogger(t),
		}, kfClient)

		for _, reactor := range r.WithReactors {
			kubeClient.PrependReactor("*", "*", reactor)
			servingclient.PrependReactor("*", "*", reactor)
			kfClient.PrependReactor("*", "*", reactor)
		}

		// Validate all Create operations through the serving client.
		servingclient.PrependReactor("create", "*", ValidateCreates)
		servingclient.PrependReactor("update", "*", ValidateUpdates)

		actionRecorderList := ActionRecorderList{servingclient, kubeClient, kfClient}
		eventList := EventList{eventRecorder}

		return c, actionRecorderList, eventList
	}
}
