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

package filter

import (
	"fmt"
)

func compileOneOf(pattern interface{}) (Matcher, error) {
	switch obj := pattern.(type) {
	case []interface{}:
		if len(obj) < 2 {
			return nil, fmt.Errorf("[oneof] should be given multiple elements, got: %d", len(obj))
		}
		var matchers []Matcher
		for _, p := range obj {
			m, err := Compile(p)
			if err != nil {
				return nil, err
			}
			matchers = append(matchers, m)
		}
		return &oneOf{
			matchers: matchers,
		}, nil
	default:
		return nil, fmt.Errorf("[oneof] must be given a list, got: %T", pattern)
	}
}

type oneOf struct {
	matchers []Matcher
}

// oneOf implement Matcher
var _ Matcher = (*oneOf)(nil)

func (oo *oneOf) Match(elt interface{}) bool {
	for _, m := range oo.matchers {
		if m.Match(elt) {
			return true
		}
	}
	return false
}
