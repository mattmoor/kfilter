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
	"bytes"
	"encoding/json"
	"github.com/ghodss/yaml"
	"text/template"
)

type Mutator interface {
	Mutate(interface{}) ([]byte, error)
}

func Compile(tmpl string) (Mutator, error) {
	// Create a new template and parse the letter into it.
	t, err := template.New("compiled").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	return &blah{t: t}, nil
}

type blah struct {
	t *template.Template
}

// blah implements Mutator
var _ Mutator = (*blah)(nil)

func (b *blah) Mutate(body interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := b.t.Execute(buf, body); err != nil {
		return nil, err
	}

	var newBody map[string]interface{}
	if err := yaml.Unmarshal(buf.Bytes(), &newBody); err != nil {
		return nil, err
	}

	// Return the final result as serialized JSON
	return json.Marshal(newBody)
}
