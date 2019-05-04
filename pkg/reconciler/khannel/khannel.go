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

package khannel

import (
	"context"
	"fmt"
	"reflect"

	eventingv1alpha1 "github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	eventingclientset "github.com/knative/eventing/pkg/client/clientset/versioned"
	eventinginformers "github.com/knative/eventing/pkg/client/informers/externalversions/eventing/v1alpha1"
	eventinglisters "github.com/knative/eventing/pkg/client/listers/eventing/v1alpha1"
	"github.com/knative/pkg/controller"
	servinginformers "github.com/knative/serving/pkg/client/informers/externalversions/serving/v1alpha1"
	servinglisters "github.com/knative/serving/pkg/client/listers/serving/v1alpha1"
	"github.com/knative/serving/pkg/reconciler"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"

	"github.com/mattmoor/kfilter/pkg/reconciler/khannel/resources"
	"github.com/mattmoor/kfilter/pkg/reconciler/khannel/resources/names"
)

const controllerAgentName = "khannel-controller"

// Reconciler is the controller implementation for Channel resources
type Reconciler struct {
	*reconciler.Base

	eventingclientset eventingclientset.Interface

	khannelImage string

	channelLister eventinglisters.ChannelLister
	serviceLister servinglisters.ServiceLister
}

// Check that we implement the controller.Reconciler interface.
var _ controller.Reconciler = (*Reconciler)(nil)

// NewController returns a new khannel controller
func NewController(
	opt reconciler.Options,
	eventingclientset eventingclientset.Interface,
	serviceInformer servinginformers.ServiceInformer,
	channelInformer eventinginformers.ChannelInformer,
	khannelImage string,
) *controller.Impl {
	r := &Reconciler{
		Base:              reconciler.NewBase(opt, controllerAgentName),
		eventingclientset: eventingclientset,
		serviceLister:     serviceInformer.Lister(),
		channelLister:     channelInformer.Lister(),
		khannelImage:      khannelImage,
	}
	impl := controller.NewImpl(r, r.Logger, "Khannels",
		reconciler.MustNewStatsReporter("Khannels", r.Logger))

	r.Logger.Info("Setting up event handlers")

	// Set up an event handler for when Khannel resources change.
	channelInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    impl.Enqueue,
		UpdateFunc: controller.PassNew(impl.Enqueue),
		DeleteFunc: impl.Enqueue,
	})

	// Set up an event handler for when Knative Service resources that we own change.
	serviceInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.Filter(eventingv1alpha1.SchemeGroupVersion.WithKind("Channel")),
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc:    impl.EnqueueControllerOf,
			UpdateFunc: controller.PassNew(impl.EnqueueControllerOf),
			DeleteFunc: impl.EnqueueControllerOf,
		},
	})

	return impl
}

// Reconcile implements controller.Reconciler
func (c *Reconciler) Reconcile(ctx context.Context, key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the channel resource with this namespace/name
	original, err := c.channelLister.Channels(namespace).Get(name)
	if errors.IsNotFound(err) {
		// The Channel resource may no longer exist, in which case we stop processing.
		runtime.HandleError(fmt.Errorf("channel %q in work queue no longer exists", key))
		return nil
	} else if err != nil {
		return err
	}

	prov := original.Spec.Provisioner
	if prov.APIVersion != "eventing.knative.dev/v1alpha1" ||
		prov.Kind != "ClusterChannelProvisioner" ||
		prov.Name != "khannel" {
		c.Logger.Infof("Skipping non-khannel channel: %q", original.Name)
		return nil
	}
	channel := original.DeepCopy()

	// Reconcile this copy of the channel and then write back any status
	// updates regardless of whether the reconciliation errored out.
	err = c.reconcile(ctx, channel)
	if equality.Semantic.DeepEqual(original.Status, channel.Status) {
		// If we didn't change anything then don't call updateStatus.
		// This is important because the copy we loaded from the informer's
		// cache may be stale and we don't want to overwrite a prior update
		// to status with this stale state.
		c.Logger.Infof("SAME STATUS")
	} else if _, err := c.updateStatus(channel); err != nil {
		c.Logger.Warn("Failed to update channel status", zap.Error(err))
		c.Recorder.Eventf(channel, corev1.EventTypeWarning, "UpdateFailed",
			"Failed to update status for Channel %q: %v", channel.Name, err)
		return err
	}

	return nil
}

func (c *Reconciler) reconcile(ctx context.Context, channel *eventingv1alpha1.Channel) error {
	channel.Status.InitializeConditions()
	if err := c.reconcileService(ctx, channel); err != nil {
		return err
	}
	return nil
}

func (c *Reconciler) reconcileService(ctx context.Context, channel *eventingv1alpha1.Channel) error {
	c.Logger.Infof("Going to reconcile channel: %q", channel.Name)
	svcName := names.KService(channel)
	service, err := c.serviceLister.Services(channel.Namespace).Get(svcName)
	if apierrs.IsNotFound(err) {
		desiredSvc := resources.MakeKService(channel, c.khannelImage)
		service, err = c.ServingClientSet.ServingV1alpha1().Services(channel.Namespace).Create(desiredSvc)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		desiredSvc := resources.MakeKService(channel, c.khannelImage)
		if !equality.Semantic.DeepEqual(service.Spec, desiredSvc.Spec) {
			service = service.DeepCopy()
			service.Spec = desiredSvc.Spec
			service, err = c.ServingClientSet.ServingV1alpha1().Services(channel.Namespace).Update(service)
			if err != nil {
				return err
			}
		}
	}

	if service.Generation == service.Status.ObservedGeneration {
		if service.Status.IsReady() {
			channel.Status.MarkProvisioned()
			channel.Status.SetAddress(service.Status.Address.Hostname)
			// We update ObservedGeneration here because it is literally the only indication that
			// is surfaced to Channel consumers that the underlying Service has been reprogrammed.
			channel.Status.ObservedGeneration = channel.Generation

			c.Logger.Infof("Ready: %v", channel.Status)
			return nil
		}
	}
	return nil
}

func (c *Reconciler) updateStatus(desired *eventingv1alpha1.Channel) (*eventingv1alpha1.Channel, error) {
	config, err := c.channelLister.Channels(desired.Namespace).Get(desired.Name)
	if err != nil {
		return nil, err
	}
	// If there's nothing to update, just return.
	if reflect.DeepEqual(config.Status, desired.Status) {
		c.Logger.Infof("IT IS EQUAL DAMNIT")
		return config, nil
	}
	c.Logger.Infof("CALLING UPDATE")
	// Don't modify the informers copy
	existing := config.DeepCopy()
	existing.Status = desired.Status
	return c.eventingclientset.EventingV1alpha1().Channels(desired.Namespace).UpdateStatus(existing)
}
