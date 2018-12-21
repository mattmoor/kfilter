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

package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/knative/pkg/cloudevents"

	"github.com/mattmoor/kfilter/pkg/filter"
)

var (
	filterType    = flag.String("type", "", "The event type to keep.")
	encodedFilter = flag.String("filter", "", "The base64 encoded filter expression.")
)

type Filter struct {
	m filter.Matcher
}

func (f *Filter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var body json.RawMessage
	ctx, err := cloudevents.Binary.FromRequest(&body, r)
	if err != nil {
		log.Printf("Failed to parse events from the request: %s", err)
		// TODO: Actually fail this request?
		w.WriteHeader(http.StatusOK)
		return
	}
	log.Printf("Received Context: %+v", ctx)
	log.Printf("Received body as: %#v", string(body))

	var unstructured map[string]interface{}
	err = json.Unmarshal(body, &unstructured)
	if err != nil {
		log.Printf("Failed to unmarshal payload: %s", err)
		// TODO: Actually fail this request?
		w.WriteHeader(http.StatusOK)
		return
	}

	// If specified, only let pull request events through.
	if *filterType != "" {
		if ctx.EventType != *filterType {
			log.Printf("Skipping: %q", ctx.EventType)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	// Check to see if the compiled filter matches the body.
	if !f.m.Match(unstructured) {
		log.Printf("Skipping: %q", ctx.EventType)
		w.WriteHeader(http.StatusOK)
		return
	}

	setHeaders(ctx, w.Header())
	w.Write(body)
}

func setHeaders(context *cloudevents.EventContext, header http.Header) {
	// These are required ones.
	header.Add(cloudevents.HeaderCloudEventsVersion, cloudevents.CloudEventsVersion)
	header.Add(cloudevents.HeaderEventID, context.EventID)
	header.Add(cloudevents.HeaderEventType, context.EventType)
	header.Add(cloudevents.HeaderSource, context.Source)

	header.Add(cloudevents.HeaderEventTime, context.EventTime.Format(time.RFC3339Nano))
	if context.EventTypeVersion != "" {
		header.Add(cloudevents.HeaderEventTypeVersion, context.EventTypeVersion)
	}
	if context.SchemaURL != "" {
		header.Add(cloudevents.HeaderSchemaURL, context.SchemaURL)
	}
	header.Add(cloudevents.HeaderContentType, context.ContentType)
}

func main() {
	flag.Parse()

	expression, err := base64.StdEncoding.DecodeString(*encodedFilter)
	if err != nil {
		log.Fatalf("Unable to decode filter expression: %v", err)
	}
	log.Printf("Got filter expression: %v", string(expression))

	var unstructured map[string]interface{}
	if err := json.Unmarshal([]byte(expression), &unstructured); err != nil {
		log.Fatalf("Unable to unmarshal filter expression: %v", err)
	}

	matcher, err := filter.Compile(unstructured)
	if err != nil {
		log.Fatalf("Unable to compile filter expression: %v", err)
	}

	f := &Filter{
		m: matcher,
	}

	http.ListenAndServe(":8080", f)
}
