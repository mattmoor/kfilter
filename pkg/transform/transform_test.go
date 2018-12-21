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

package transform

import (
	"testing"
)

func TestTransforms(t *testing.T) {
	tests := []struct {
		name     string
		template string
		input    interface{}
		want     string
	}{{
		name:     "fixed payload",
		template: `foo: bar`,
		input: map[string]interface{}{
			"asdf": "baz",
		},
		want: `{"foo":"bar"}`,
	}, {
		name: "simple substitution",
		template: `
foo:
  bar: {{ .asdf }}`,
		input: map[string]interface{}{
			"asdf": "baz",
		},
		want: `{"foo":{"bar":"baz"}}`,
	}, {
		name:     "conditional substitution (true)",
		template: `foo: {{ if eq .asdf "blah" }} bar {{ else }} baz {{ end }}`,
		input: map[string]interface{}{
			"asdf": "blah",
		},
		want: `{"foo":"bar"}`,
	}, {
		name:     "conditional substitution (false)",
		template: `foo: {{ if eq .asdf "blah" }} bar {{ else }} baz {{ end }}`,
		input: map[string]interface{}{
			"asdf": "not blah",
		},
		want: `{"foo":"baz"}`,
	}, {
		name: "non-trivial rewrite",
		template: `
{{ if eq .action "push" }}
commit: {{ .commit }}
comments:
{{ range $i, $v := .comments }}
- {{ $v.author }}: {{ $v.body }}
{{ end }}
{{ end }}
`,
		input: map[string]interface{}{
			"action": "push",
			"commit": "deadbeef",
			"comments": []interface{}{
				map[string]interface{}{
					"author": "mattmoor",
					"body":   "omg wtf?!?!",
				},
				map[string]interface{}{
					"author": "jonjohnsonjr",
					"body":   "NO U WTF!!!1!!",
				},
			},
		},
		want: `{"comments":[{"mattmoor":"omg wtf?!?!"},{"jonjohnsonjr":"NO U WTF!!!1!!"}],"commit":"deadbeef"}`,
	}, {
		name: "filter complex message via go template.",
		template: `
{{ if eq .action "push" }}
commit: {{ .commit }}
comments:
{{ range $i, $v := .comments }}
- {{ $v.author }}: {{ $v.body }}
{{ end }}
{{ end }}
`,
		input: map[string]interface{}{
			"action": "not a push",
			"comments": []interface{}{
				map[string]interface{}{
					"author": "mattmoor",
					"body":   "omg wtf?!?!",
				},
				map[string]interface{}{
					"author": "jonjohnsonjr",
					"body":   "NO U WTF!!!1!!",
				},
			},
		},
		want: `null`,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := Compile(test.template)
			if err != nil {
				t.Fatalf("Compile() = %v", err)
			}
			got, err := m.Mutate(test.input)
			if err != nil {
				t.Fatalf("m.Mutate() = %v", err)
			}
			if string(got) != test.want {
				t.Errorf("m.Mutate() = %v, wanted %v", string(got), test.want)
			}
		})
	}
}
