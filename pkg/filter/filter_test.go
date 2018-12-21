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
	"testing"
)

func TestMatches(t *testing.T) {
	tests := []struct {
		name    string
		pattern interface{}
		input   interface{}
		want    bool
	}{{
		name: "partial-literal, exact match",
		pattern: map[string]interface{}{
			"foo": "bar",
		},
		input: map[string]interface{}{
			"foo": "bar",
		},
		want: true,
	}, {
		name: "partial-literal, no match",
		pattern: map[string]interface{}{
			"foo": "bar",
		},
		input: map[string]interface{}{
			"foo": "baz",
		},
		want: false,
	}, {
		name: "partial-literal, exact match, all primitive types",
		pattern: map[string]interface{}{
			"foo":  "bar",
			"baz":  true,
			"blah": 1234.3,
		},
		input: map[string]interface{}{
			"foo":  "bar",
			"baz":  true,
			"blah": 1234.3,
		},
		want: true,
	}, {
		name: "partial-literal, exact slice match, all primitive types",
		pattern: []interface{}{
			"bar",
			true,
			1234.3,
		},
		input: []interface{}{
			"bar",
			true,
			1234.3,
		},
		want: true,
	}, {
		name: "partial-literal, partial match",
		pattern: map[string]interface{}{
			"foo": "bar",
		},
		input: map[string]interface{}{
			"foo": "bar",
			"baz": true,
		},
		want: true,
	}, {
		name: "partial-literal, anything values",
		pattern: map[string]interface{}{
			"foo": "[anything]",
			"baz": "[anything]",
		},
		input: map[string]interface{}{
			"foo": "bar",
			"baz": true,
		},
		want: true,
	}, {
		name: "partial-literal, anything matches, all types",
		pattern: []interface{}{
			"[anything]",
			"[anything]",
			"[anything]",
		},
		input: []interface{}{
			"bar",
			true,
			1234.3,
			"extra element",
		},
		want: true,
	}, {
		name: "partial-literal, anything matches, not enough values",
		pattern: []interface{}{
			"[anything]",
			"[anything]",
			"[anything]",
		},
		input: []interface{}{
			"bar",
			true,
			// If we don't have three, then this will fail.
		},
		want: false,
	}, {
		name: "exact-literal, partial match fails",
		pattern: map[string]interface{}{
			"[exact]": map[string]interface{}{
				"foo": "bar",
			},
		},
		input: map[string]interface{}{
			"foo": "bar",
			"baz": true,
		},
		want: false,
	}, {
		name: "exact-literal, exact match",
		pattern: map[string]interface{}{
			"[exact]": map[string]interface{}{
				"foo": "bar",
			},
		},
		input: map[string]interface{}{
			"foo": "bar",
		},
		want: true,
	}, {
		name: "exact-literal, exact match",
		pattern: map[string]interface{}{
			"[exact]": []interface{}{
				"[anything]",
				"[anything]",
				"[anything]",
			},
		},
		input: []interface{}{
			"bar",
			true,
			1234.3,
		},
		want: true,
	}, {
		name: "exact-literal, partial match fails",
		pattern: map[string]interface{}{
			"[exact]": []interface{}{
				"[anything]",
				"[anything]",
				"[anything]",
			},
		},
		input: []interface{}{
			"bar",
			true,
			1234.3,
			"extra elements",
		},
		want: false,
	}, {
		name: "oneof, match first element",
		pattern: map[string]interface{}{
			"[oneof]": []interface{}{
				"foo",
				"bar",
			},
		},
		input: "foo",
		want:  true,
	}, {
		name: "oneof, match second element",
		pattern: map[string]interface{}{
			"[oneof]": []interface{}{
				"foo",
				"bar",
			},
		},
		input: "bar",
		want:  true,
	}, {
		name: "oneof, match no element",
		pattern: map[string]interface{}{
			"[oneof]": []interface{}{
				"foo",
				"bar",
			},
		},
		input: "baz",
		want:  false,
	}, {
		name: "map literal doesn't match non-map",
		pattern: map[string]interface{}{
			"foo": "bar",
		},
		input: []interface{}{
			"bar",
		},
		want: false,
	}, {
		name: "slice literal doesn't match non-slice",
		pattern: []interface{}{
			"bar",
		},
		input: map[string]interface{}{
			"foo": "bar",
		},
		want: false,
	}, {
		name: "partial-literal, too few keys",
		pattern: map[string]interface{}{
			"foo": "bar",
			"baz": "blah",
		},
		input: map[string]interface{}{
			"foo": "bar",
		},
		want: false,
	}, {
		name: "partial-literal, missing key",
		pattern: map[string]interface{}{
			"foo": "bar",
			"baz": "asdf",
		},
		input: map[string]interface{}{
			"foo":  "bar",
			"blah": true,
		},
		want: false,
	}, {
		name: "exact-literal, missing key",
		pattern: map[string]interface{}{
			"[exact]": map[string]interface{}{
				"foo": "bar",
				"baz": "asdf",
			},
		},
		input: map[string]interface{}{
			"foo":  "bar",
			"blah": true,
		},
		want: false,
	}, {
		name: "slice literal doesn't match element",
		pattern: []interface{}{
			"foo",
		},
		input: []interface{}{
			"bar",
		},
		want: false,
	}, {
		name:    "bool doesn't match non-bool",
		pattern: true,
		input:   "asdf",
		want:    false,
	}, {
		name:    "string doesn't match non-string",
		pattern: "foo",
		input:   true,
		want:    false,
	}, {
		name:    "float doesn't match non-float",
		pattern: 1234.5,
		input:   true,
		want:    false,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := Compile(test.pattern)
			if err != nil {
				t.Fatalf("Error compiling pattern %#v: %v", test.pattern, err)
			}
			if m == nil {
				t.Fatal("Compile = nil, wanted non-nil")
			}
			if got := m.Match(test.input); got != test.want {
				t.Errorf("m.Match(%#v) = %v, wanted %v", test.input, got, test.want)
			}
		})
	}
}

func TestCompileFailures(t *testing.T) {
	tests := []struct {
		name    string
		pattern interface{}
	}{{
		name:    "just an int",
		pattern: 36,
	}, {
		name: "empty oneof",
		pattern: map[string]interface{}{
			"[oneof]": []interface{}{},
		},
	}, {
		name: "single oneof",
		pattern: map[string]interface{}{
			"[oneof]": []interface{}{"asdf"},
		},
	}, {
		name: "oneof without list",
		pattern: map[string]interface{}{
			"[oneof]": "asfd",
		},
	}, {
		name: "oneof with nested error",
		pattern: map[string]interface{}{
			"[oneof]": []interface{}{
				"asdf",
				map[string]interface{}{
					"[oneof]": []interface{}{},
				},
			},
		},
	}, {
		name: "map with nested error",
		pattern: map[string]interface{}{
			"asdf": 123,
		},
	}, {
		name: "slice with nested error",
		pattern: []interface{}{
			123,
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := Compile(test.pattern)
			if err == nil {
				t.Fatalf("Compile() = %#v, wanted error", m)
			}
		})
	}
}
