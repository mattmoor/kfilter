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

package v1alpha1

import (
	"encoding/json"

	duckv1alpha1 "github.com/knative/pkg/apis/duck/v1alpha1"
	"github.com/knative/pkg/kmeta"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Filter is a specification for a Filter resource
type Filter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FilterSpec   `json:"spec"`
	Status FilterStatus `json:"status"`
}

// Check that we can create OwnerReferences to a Filter.
var _ kmeta.OwnerRefable = (*Filter)(nil)

// FilterSpec is the spec for a Filter resource
type FilterSpec struct {
	// The cloud event type to keep
	// +optional
	EventType string `json:"eventType,omitempty"`

	// The filter to apply.
	// TODO(mattmoor): More detailed description.
	// +optional
	Body json.RawMessage `json:"body,omitempty"`
}

// FilterStatus is the status for a Filter resource
type FilterStatus struct {
	// Address holds the information needed for a Filter to be the target of an event.
	// +optional
	Address *duckv1alpha1.Addressable `json:"address,omitempty"`

	// Conditions communicates information about ongoing/complete
	// reconciliation processes that bring the "spec" inline with the observed
	// state of the world.
	// +optional
	Conditions duckv1alpha1.Conditions `json:"conditions,omitempty"`
}

func (r *Filter) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("Filter")
}

func (rs *FilterStatus) GetCondition(t duckv1alpha1.ConditionType) *duckv1alpha1.Condition {
	return condSet.Manage(rs).GetCondition(t)
}

func (rs *FilterStatus) InitializeConditions() {
	condSet.Manage(rs).InitializeConditions()
}

func (rs *FilterStatus) PropagateServiceStatus(ss v1alpha1.ServiceStatus) {
	rs.Address = ss.Address
	sr := ss.GetCondition(v1alpha1.ServiceConditionReady)
	if sr == nil {
		return
	}
	switch sr.Status {
	case corev1.ConditionTrue:
		condSet.Manage(rs).MarkTrue(ConditionServiceReady)
	case corev1.ConditionUnknown, corev1.ConditionFalse:
		condSet.Manage(rs).MarkFalse(ConditionServiceReady, sr.Reason, sr.Message)
	}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FilterList is a list of Filter resources
type FilterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Filter `json:"items"`
}
