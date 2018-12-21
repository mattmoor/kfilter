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

package names

import (
	kfv1alpha1 "github.com/mattmoor/kfilter/pkg/apis/kfilter/v1alpha1"
)

// KService returns the name of the Knative Service child resource for given Transform
func KService(i *kfv1alpha1.Transform) string {
	return i.Name
}
