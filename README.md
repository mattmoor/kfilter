# Filter CRD

The `Filter` CRD is a prototype that builds on the Knative extensibility model
to implement filter functions for use with `knative/eventing` at a higher level!

### Problem Statement

As is, the low-level primitives that `knative/eventing` provides are extremely
powerful, supporting: delivery, fanout, and composition.  Composition is
extremely cool because the first function can transform or filter the message
that's queued into the result channel.

However, to express even simple filters, you currently have to write a `Service`
that performs the filtering (by returning an empty response), and preserves the
appropriate cloud event headers.  You then stitch it in via:

```yaml
apiVersion: eventing.knative.dev/v1alpha1
kind: Subscription
metadata:
  name: filter-subscription
spec:
  channel:
    apiVersion: eventing.knative.dev/v1alpha1
    kind: Channel
    name: githubchannel
  subscriber:
    ref:
      apiVersion: serving.knative.dev/v1alpha1
      kind: Service
      name: im-a-filter
  reply:
    channel:
      apiVersion: eventing.knative.dev/v1alpha1
      kind: Channel
      name: filtered-githubchannel

```

The purpose of this repository is to illustrate that the Knative extensibility
model is designed to enable folks to build even higher-level abstractions than
functions.


### The Filter CRD

The Filter CRD is an abstraction that builds on `knative/serving` to provide a
scale-to-zero "filter" that can compose with `knative/eventing` and be entirely
described via a CRD (vs. code that needs to be containerized).

You would describe your `Filter` via a resource like:

```yaml
apiVersion: kfilter.mattmoor.io/v1alpha1
kind: Filter
metadata:
  name: im-a-filter
spec:
  # Keep pull request events only.
  eventType: dev.knative.source.github.issues
```

This can then be linked as a filter via:

```yaml
apiVersion: eventing.knative.dev/v1alpha1
kind: Subscription
metadata:
  name: filter-subscription
spec:
  channel:
    apiVersion: eventing.knative.dev/v1alpha1
    kind: Channel
    name: githubchannel
  subscriber:
    ref:
      apiVersion: kfilter.mattmoor.io/v1alpha1
      kind: Filter
      name: im-a-filter
  reply:
    channel:
      apiVersion: eventing.knative.dev/v1alpha1
      kind: Channel
      name: filtered-githubchannel

```

Currently this only supports a trivial filter: cloud event type, but atop this
plumbing we can extend the filter power arbitrarily, e.g. matching elements of
the payload body.


### Try it.

The filters this supports today are not particularly exciting, but if you'd like
to try it out `release.yaml` contains all of the components needed to install it
on your cluster.
