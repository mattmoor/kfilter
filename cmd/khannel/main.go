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
	"context"
	"flag"
	"log"

	cloudevents "github.com/cloudevents/sdk-go"
	"golang.org/x/sync/errgroup"

	"github.com/mattmoor/kfilter/pkg/reconciler/khannel/resources"
)

var (
	subscriptions = flag.String("subscriptions", "", "The base64 encoded JSON list of subscriptions.")
)

func main() {
	flag.Parse()

	subs := resources.DecodeSubscriptions(*subscriptions)
	log.Printf("Got subscriptions: %v", subs)

	client, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("OMG VOLCANO: %v", err)
	}

	doshit := func(ctx context.Context, ce cloudevents.Event, cer *cloudevents.EventResponse) error {
		grp := errgroup.Group{}
		for _, sub := range subs {
			sub := sub
			grp.Go(func() error {
				var err error
				resp := &ce
				if sub.Subscriber != "" {
					subCtx := cloudevents.ContextWithTarget(ctx, sub.Subscriber)
					resp, err = client.Send(subCtx, ce)
					if err != nil {
						return err
					}
				}

				if sub.ReplyTo != "" && resp != nil {
					subCtx := cloudevents.ContextWithTarget(ctx, sub.ReplyTo)
					if _, err := client.Send(subCtx, ce); err != nil {
						return err
					}
				}
				return nil
			})
		}
		return grp.Wait()
	}

	client.StartReceiver(context.Background(), doshit)
}
