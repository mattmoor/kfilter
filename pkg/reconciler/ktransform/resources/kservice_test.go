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

package resources

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kfv1alpha1 "github.com/mattmoor/kfilter/pkg/apis/kfilter/v1alpha1"
)

func TestMakeVirtualService(t *testing.T) {
	boolTrue := true

	tests := []struct {
		name string
		kf   *kfv1alpha1.Transform
		img  string
		want *v1alpha1.Service
	}{{
		name: "test simple alias",
		kf: &kfv1alpha1.Transform{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "baz",
			},
			// TODO(mattmoor): Spec
		},
		img: "foo",
		want: &v1alpha1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "baz",
				OwnerReferences: []metav1.OwnerReference{{
					APIVersion:         "kfilter.mattmoor.io/v1alpha1",
					Kind:               "Transform",
					Name:               "foo",
					Controller:         &boolTrue,
					BlockOwnerDeletion: &boolTrue,
				}},
			},
			Spec: v1alpha1.ServiceSpec{
				RunLatest: &v1alpha1.RunLatestType{
					Configuration: v1alpha1.ConfigurationSpec{
						RevisionTemplate: v1alpha1.RevisionTemplateSpec{
							Spec: v1alpha1.RevisionSpec{
								Container: corev1.Container{
									Image: "foo",
								},
							},
						},
					},
				},
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := MakeKService(test.kf, test.img)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Unexpected KService (-want +got): %v", diff)
			}
		})
	}
}
