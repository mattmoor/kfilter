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

// Package filter implements a collections of matchers for accepting or
// rejecting messages based on some simple structural rules.
// 1. Partial-Literal matching
//   The main mode of operation is to match specified literals, so if
//   we get the filter expression like: {"foo": "bar"} then it will
//   match maps with that shape.  e.g. {"foo": "bar", "baz": true}
//   Note that this can only prefix-match arrays, see [anything].
//
// 2. Exact-Literal matching
//   This would kick in when we are passed a pattern with the shape:
//     {"[exact]": {"foo": "bar"}}
//   The result would be something very similar to #1, however, we
//   would not accept the unspecified "baz" key.
//   This applies to exactly one level of the object.
//
// 3. One Of matching
//   This would kick in when we are passed a pattern with the shape:
//     {"[oneof]": [{"foo": "bar"}, {"baz": "blah"}]}
//
// 4. Match Anything
//   This would kick in when we are passed a pattern with the shape:
//     [anything]
//   The idea is to allow folks to blank out values in partial literal
//   matching of arrays, or require a key without specifying the value.
//   e.g. {"foo": ["bar", "[anything]", "baz"]}
//   e.g. {"foo": "[anything]"}
//
// 5. TODO(mattmoor): Regexp Match
package filter
