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
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/knative/pkg/cloudevents"
)

var (
	filterType = flag.String("type", "", "The event type to keep.")
)

type Filter struct{}

type event struct {
	Data string `json:"data,omitEmpty"`
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
	log.Printf("Received body as: %q", string(body))
	e := event{}
	err = json.Unmarshal(body, &e)
	if err != nil {
		log.Printf("Failed to unmarshal event data: %s", err)
		// TODO: Actually fail this request?
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only let pull request events through.
	if ctx.EventType != *filterType {
		log.Printf("Skipping: %q", ctx.EventType)
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Printf("Received event as: %+v", e)
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

	http.ListenAndServe(":8080", &Filter{})
}
