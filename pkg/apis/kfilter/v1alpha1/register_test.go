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
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
)

func TestRegisterHelpers(t *testing.T) {
	if got, want := Kind("Revision"), "Revision.kfilter.mattmoor.io"; got.String() != want {
		t.Errorf("Kind(Filter) = %v, want %v", got.String(), want)
	}

	if got, want := Resource("Filter"), "Filter.kfilter.mattmoor.io"; got.String() != want {
		t.Errorf("Resource(Filter) = %v, want %v", got.String(), want)
	}

	if got, want := SchemeGroupVersion.String(), "kfilter.mattmoor.io/v1alpha1"; got != want {
		t.Errorf("SchemeGroupVersion() = %v, want %v", got, want)
	}

	scheme := runtime.NewScheme()
	if err := addKnownTypes(scheme); err != nil {
		t.Errorf("addKnownTypes() = %v", err)
	}
}