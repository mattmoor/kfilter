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
	"encoding/base64"

	"github.com/knative/pkg/kmeta"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kfv1alpha1 "github.com/mattmoor/kfilter/pkg/apis/kfilter/v1alpha1"
	"github.com/mattmoor/kfilter/pkg/reconciler/kfilter/resources/names"
)

func MakeKService(kf *kfv1alpha1.Filter, image string) *v1alpha1.Service {
	encodedFilter := base64.StdEncoding.EncodeToString(kf.Spec.Body)

	return &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            names.KService(kf),
			Namespace:       kf.Namespace,
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(kf)},
			Annotations:     kf.ObjectMeta.Annotations,
		},
		Spec: v1alpha1.ServiceSpec{
			RunLatest: &v1alpha1.RunLatestType{
				Configuration: v1alpha1.ConfigurationSpec{
					RevisionTemplate: v1alpha1.RevisionTemplateSpec{
						Spec: v1alpha1.RevisionSpec{
							Container: corev1.Container{
								Image: image,
								Args: []string{
									"-type", kf.Spec.EventType,
									"-filter", encodedFilter,
								},
							},
						},
					},
				},
			},
		},
	}
}
