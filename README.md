# Filter and Transform CRDs

The `Filter` and `Transform` CRDs are prototypes that build on the Knative
extensibility model to implement filter / transformation functions for use with
`knative/eventing` at a higher level!

### Problem Statement

As is, the low-level primitives that `knative/eventing` provides are extremely
powerful, supporting: delivery, fanout, and composition.  Composition is
extremely cool because the first function can transform or filter the message
that's queued into the result channel.

However, to express even simple filters and transformations, you currently have
to write a `Service` that performs the filtering (by returning an empty
response), and preserves the appropriate cloud event headers.  You then stitch
it in via:

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


## The Filter CRD

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

### Supported Filters

#### Event Type

As shown in the above example, you can filter on event types via:

```yaml
apiVersion: kfilter.mattmoor.io/v1alpha1
kind: Filter
metadata:
  name: im-a-filter
spec:
  # Keep pull request events only.
  eventType: dev.knative.source.github.issues
```

#### Body Patterns

You can match the body of the event via pattern matches.  The default mode
of matching is partial matches, so for example the following pattern:

```yaml
apiVersion: kfilter.mattmoor.io/v1alpha1
kind: Filter
metadata:
  name: im-a-filter
spec:
  body: {
    "foo": "bar"
  }
```

Will match any body that has a key `"foo"` with value `"bar"`.  For arrays
the partial match will accept a matching array prefix.

To turn a match from a partial match into an exact match you can wrap it as
follows:

```yaml
apiVersion: kfilter.mattmoor.io/v1alpha1
kind: Filter
metadata:
  name: im-a-filter
spec:
  body: {
    "[exact]": {
      "foo": "bar"
    }
  }
```

To match an arbitrary element, you can use the special string `[anything]`,
so for example if we want to match messages that have the key `"foo"` with
a slice containing at least three elements we can write:

```yaml
apiVersion: kfilter.mattmoor.io/v1alpha1
kind: Filter
metadata:
  name: im-a-filter
spec:
  body: {
    "foo": [
      "[anything]",
      "[anything]",
      "[anything]"
    ]
  }
```

To match one of several values, you can use the `[oneof]` keyword:

```yaml
apiVersion: kfilter.mattmoor.io/v1alpha1
kind: Filter
metadata:
  name: im-a-filter
spec:
  body: {
    "foo": {"[oneof]": ["bar", "baz"]}
  }
```

This will match any message with a top level key `"foo"` iff it has one of the
values: `"bar"` or `"baz"`.

## The Transform CRD

The Transform CRD is an abstraction that builds on `knative/serving` to provide a
scale-to-zero "transform" that can compose with `knative/eventing` and be entirely
described via a CRD (vs. code that needs to be containerized).

You would describe your `Transform` via a resource like:

```yaml
apiVersion: kfilter.mattmoor.io/v1alpha1
kind: Transform
metadata:
  name: im-a-transform
spec:
  template: |
    # Any Go template acting on the expected body.
    # If the result is empty, this will filter the result.
    foo: {{ .bar }}
```

This can then be linked as a transformation via:

```yaml
apiVersion: eventing.knative.dev/v1alpha1
kind: Subscription
metadata:
  name: transform-subscription
spec:
  channel:
    apiVersion: eventing.knative.dev/v1alpha1
    kind: Channel
    name: githubchannel
  subscriber:
    ref:
      apiVersion: kfilter.mattmoor.io/v1alpha1
      kind: Transform
      name: im-a-transform
  reply:
    channel:
      apiVersion: eventing.knative.dev/v1alpha1
      kind: Channel
      name: transformed-githubchannel

```



## Try it.

The filters this supports today are not particularly exciting, but if you'd like
to try it out `release.yaml` contains all of the components needed to install it
on your cluster.
