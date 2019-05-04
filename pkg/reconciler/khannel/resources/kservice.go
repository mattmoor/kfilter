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
	"encoding/json"

	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	eventingv1alpha1 "github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	"github.com/mattmoor/kfilter/pkg/reconciler/khannel/resources/names"
)

type Subscription struct {
	Subscriber string `json:"subscriber"`
	ReplyTo    string `json:"replyTo,omitempty"`
}

func EncodeSubscriptions(s []Subscription) string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err.Error())
	}
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeSubscriptions(raw string) (s []Subscription) {
	jason, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		panic(err.Error())
	}
	if err := json.Unmarshal(jason, &s); err != nil {
		panic(err.Error())
	}
	return
}

func MakeKService(channel *eventingv1alpha1.Channel, image string) *v1alpha1.Service {
	s := []Subscription{}
	if channel.Spec.Subscribable != nil {
		for _, sub := range channel.Spec.Subscribable.Subscribers {
			s = append(s, Subscription{
				Subscriber: sub.SubscriberURI,
				ReplyTo:    sub.ReplyURI,
			})
		}
	}

	encodedSubscriptions := EncodeSubscriptions(s)

	boolTrue := true
	return &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      names.KService(channel),
			Namespace: channel.Namespace,
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         "eventing.knative.dev/v1alpha1",
				BlockOwnerDeletion: &boolTrue,
				Controller:         &boolTrue,
				Kind:               "Channel",
				Name:               channel.Name,
				UID:                channel.UID,
			}},
			Annotations: channel.ObjectMeta.Annotations,
		},
		Spec: v1alpha1.ServiceSpec{
			RunLatest: &v1alpha1.RunLatestType{
				Configuration: v1alpha1.ConfigurationSpec{
					RevisionTemplate: v1alpha1.RevisionTemplateSpec{
						Spec: v1alpha1.RevisionSpec{
							Container: corev1.Container{
								Image: image,
								Args: []string{
									"-subscriptions", encodedSubscriptions,
								},
							},
						},
					},
				},
			},
		},
	}
}
