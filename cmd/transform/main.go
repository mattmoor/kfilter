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

	"github.com/mattmoor/kfilter/pkg/transform"
)

var (
	encodedTransform = flag.String("transform", "", "The base64 encoded transform expression.")
)

type Transform struct {
	m transform.Mutator
}

func (f *Transform) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Failed to unmarshal request body: %s", err)
		// TODO: Actually fail this request?
		w.WriteHeader(http.StatusOK)
		return
	}

	// Apply the compiled transformation to the event body.
	result, err := f.m.Mutate(payload)
	if err != nil {
		log.Printf("Failed to unmarshal payload: %s", err)
		// TODO: Actually fail this request?
		w.WriteHeader(http.StatusOK)
		return
	}
	// If the transformation returns null, then filter this message.
	if string(result) == "null" {
		log.Printf("Skipping: %q", ctx.EventType)
		w.WriteHeader(http.StatusOK)
		return
	}

	setHeaders(ctx, w.Header())
	w.Write(result)
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

	template, err := base64.StdEncoding.DecodeString(*encodedTransform)
	if err != nil {
		log.Fatalf("Unable to decode transform expression: %v", err)
	}
	log.Printf("Got transform expression: %v", string(template))

	mutator, err := transform.Compile(string(template))
	if err != nil {
		log.Fatalf("Unable to compile transform expression: %v", err)
	}

	f := &Transform{
		m: mutator,
	}

	http.ListenAndServe(":8080", f)
}
