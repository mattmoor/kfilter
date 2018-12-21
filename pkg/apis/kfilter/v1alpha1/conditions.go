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
	duckv1alpha1 "github.com/knative/pkg/apis/duck/v1alpha1"
)

const (
	// TransformConditionReady is set when the service is configured
	// and has available backends ready to receive traffic.
	ConditionReady = duckv1alpha1.ConditionReady

	// TransformConditionServiceReady is set to whether the underlying
	// Service has come up.
	ConditionServiceReady duckv1alpha1.ConditionType = "ServiceReady"
)

var condSet = duckv1alpha1.NewLivingConditionSet(ConditionServiceReady)
