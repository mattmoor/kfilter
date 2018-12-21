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

func compileLiteral(pattern interface{}, exact bool) (Matcher, error) {
	switch obj := pattern.(type) {
	case map[string]interface{}:
		if len(obj) == 1 {
			for k, v := range obj {
				switch k {
				case "[oneof]":
					return compileOneOf(v)
				case "[exact]":
					return compileLiteral(v, true)
				default: // Not a keyword
				}
			}
		}
		matchers := make(map[string]Matcher)
		for k, v := range obj {
			m, err := Compile(v)
			if err != nil {
				return nil, err
			}
			matchers[k] = m
		}
		return &mapLiteral{
			matchers: matchers,
			exact:    exact,
		}, nil

	case []interface{}:
		matchers := make([]Matcher, 0, len(obj))
		for _, v := range obj {
			m, err := Compile(v)
			if err != nil {
				return nil, err
			}
			matchers = append(matchers, m)
		}
		return &sliceLiteral{
			matchers: matchers,
			exact:    exact,
		}, nil

	case string:
		// Check for keywords
		switch obj {
		case "[anything]":
			return &anything{}, nil
		default:
			return (*stringLiteral)(&obj), nil
		}
	case float64:
		return (*floatLiteral)(&obj), nil
	case bool:
		return (*boolLiteral)(&obj), nil

	default:
		return nil, fmt.Errorf("Unrecognized type: %T", pattern)
	}
}

type mapLiteral struct {
	matchers map[string]Matcher
	exact    bool
}

// mapLiteral implement Matcher
var _ Matcher = (*mapLiteral)(nil)

func (ml *mapLiteral) Match(elt interface{}) bool {
	obj, ok := elt.(map[string]interface{})
	if !ok {
		return false
	}
	if got, want := len(obj), len(ml.matchers); got != want {
		if got < want {
			// We always expect the "want" keys.
			return false
		}
		// got > want
		if ml.exact {
			return false
		}
	}

	seen := 0
	for key, match := range ml.matchers {
		value, ok := obj[key]
		if !ok {
			if ml.exact {
				// Bail early, not allowed.
				return false
			}
			continue
		}
		if !match.Match(value) {
			// The value for this key does not match.
			return false
		}
		seen++
	}

	// If we saw all of the desired keys, then this is a match.
	return seen == len(ml.matchers)
}

type sliceLiteral struct {
	matchers []Matcher
	exact    bool
}

// sliceLiteral implement Matcher
var _ Matcher = (*sliceLiteral)(nil)

func (ml *sliceLiteral) Match(elt interface{}) bool {
	obj, ok := elt.([]interface{})
	if !ok {
		return false
	}
	if got, want := len(obj), len(ml.matchers); got != want {
		if got < want {
			// We always expect the "want" keys.
			return false
		}
		// got > want
		if ml.exact {
			return false
		}
	}

	seen := 0
	for idx, match := range ml.matchers {
		value := obj[idx]
		if !match.Match(value) {
			// The value for this idx does not match.
			return false
		}
		seen++
	}

	return seen == len(ml.matchers)
}

type stringLiteral string

// stringLiteral implement Matcher
var _ Matcher = (*stringLiteral)(nil)

func (sl *stringLiteral) Match(elt interface{}) bool {
	s, ok := elt.(string)
	if !ok {
		return false
	}
	return string(*sl) == s
}

type boolLiteral bool

// boolLiteral implement Matcher
var _ Matcher = (*boolLiteral)(nil)

func (sl *boolLiteral) Match(elt interface{}) bool {
	s, ok := elt.(bool)
	if !ok {
		return false
	}
	return bool(*sl) == s
}

type floatLiteral float64

// floatLiteral implement Matcher
var _ Matcher = (*floatLiteral)(nil)

func (sl *floatLiteral) Match(elt interface{}) bool {
	s, ok := elt.(float64)
	if !ok {
		return false
	}
	return float64(*sl) == s
}
