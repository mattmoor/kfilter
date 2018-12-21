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

package ktransform

import (
	"context"
	"fmt"
	"reflect"

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
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"

	kfv1alpha1 "github.com/mattmoor/kfilter/pkg/apis/kfilter/v1alpha1"
	clientset "github.com/mattmoor/kfilter/pkg/client/clientset/versioned"
	kfilterscheme "github.com/mattmoor/kfilter/pkg/client/clientset/versioned/scheme"
	informers "github.com/mattmoor/kfilter/pkg/client/informers/externalversions/kfilter/v1alpha1"
	listers "github.com/mattmoor/kfilter/pkg/client/listers/kfilter/v1alpha1"
	"github.com/mattmoor/kfilter/pkg/reconciler/ktransform/resources"
	"github.com/mattmoor/kfilter/pkg/reconciler/ktransform/resources/names"
)

const controllerAgentName = "ktransform-controller"

// Reconciler is the controller implementation for Filter resources
type Reconciler struct {
	*reconciler.Base

	kfilterclientset clientset.Interface

	transformImage string

	serviceLister   servinglisters.ServiceLister
	transformLister listers.TransformLister
}

// Check that we implement the controller.Reconciler interface.
var _ controller.Reconciler = (*Reconciler)(nil)

func init() {
	// Add ktransform-controller types to the default Kubernetes Scheme so Events can be
	// logged for ktransform-controller types.
	kfilterscheme.AddToScheme(scheme.Scheme)
}

// NewController returns a new kfilter controller
func NewController(
	opt reconciler.Options,
	kfilterclientset clientset.Interface,
	serviceInformer servinginformers.ServiceInformer,
	transformInformer informers.TransformInformer,
	transformImage string,
) *controller.Impl {
	r := &Reconciler{
		Base:             reconciler.NewBase(opt, controllerAgentName),
		kfilterclientset: kfilterclientset,
		serviceLister:    serviceInformer.Lister(),
		transformLister:  transformInformer.Lister(),
		transformImage:   transformImage,
	}
	impl := controller.NewImpl(r, r.Logger, "Transforms",
		reconciler.MustNewStatsReporter("Transforms", r.Logger))

	r.Logger.Info("Setting up event handlers")

	// Set up an event handler for when Transform resources change.
	transformInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    impl.Enqueue,
		UpdateFunc: controller.PassNew(impl.Enqueue),
		DeleteFunc: impl.Enqueue,
	})

	// Set up an event handler for when Knative Service resources that we own change.
	serviceInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.Filter(kfv1alpha1.SchemeGroupVersion.WithKind("Transform")),
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

	// Get the transform resource with this namespace/name
	original, err := c.transformLister.Transforms(namespace).Get(name)
	if errors.IsNotFound(err) {
		// The Transform resource may no longer exist, in which case we stop processing.
		runtime.HandleError(fmt.Errorf("transform %q in work queue no longer exists", key))
		return nil
	} else if err != nil {
		return err
	}
	kf := original.DeepCopy()

	// Reconcile this copy of the transform and then write back any status
	// updates regardless of whether the reconciliation errored out.
	err = c.reconcile(ctx, kf)
	if equality.Semantic.DeepEqual(original.Status, kf.Status) {
		// If we didn't change anything then don't call updateStatus.
		// This is important because the copy we loaded from the informer's
		// cache may be stale and we don't want to overwrite a prior update
		// to status with this stale state.
	} else if _, err := c.updateStatus(kf); err != nil {
		c.Logger.Warn("Failed to update ktransform status", zap.Error(err))
		c.Recorder.Eventf(kf, corev1.EventTypeWarning, "UpdateFailed",
			"Failed to update status for Transform %q: %v", kf.Name, err)
		return err
	}

	return nil
}

func (c *Reconciler) reconcile(ctx context.Context, kf *kfv1alpha1.Transform) error {
	if err := c.reconcileService(ctx, kf); err != nil {
		return err
	}
	return nil
}

func (c *Reconciler) reconcileService(ctx context.Context, kf *kfv1alpha1.Transform) error {
	svcName := names.KService(kf)
	service, err := c.serviceLister.Services(kf.Namespace).Get(svcName)
	if apierrs.IsNotFound(err) {
		desiredSvc := resources.MakeKService(kf, c.transformImage)
		service, err = c.ServingClientSet.ServingV1alpha1().Services(kf.Namespace).Create(desiredSvc)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		desiredSvc := resources.MakeKService(kf, c.transformImage)
		if !equality.Semantic.DeepEqual(service.Spec, desiredSvc.Spec) {
			service = service.DeepCopy()
			service.Spec = desiredSvc.Spec
			service, err = c.ServingClientSet.ServingV1alpha1().Services(kf.Namespace).Update(service)
			if err != nil {
				return err
			}
		}
	}

	// Make our status reflect the Service's status.
	kf.Status.PropagateServiceStatus(service.Status)

	return nil
}

func (c *Reconciler) updateStatus(desired *kfv1alpha1.Transform) (*kfv1alpha1.Transform, error) {
	config, err := c.transformLister.Transforms(desired.Namespace).Get(desired.Name)
	if err != nil {
		return nil, err
	}
	// If there's nothing to update, just return.
	if reflect.DeepEqual(config.Status, desired.Status) {
		return config, nil
	}
	// Don't modify the informers copy
	existing := config.DeepCopy()
	existing.Status = desired.Status
	// TODO: for CRD there's no updatestatus, so use normal update
	return c.kfilterclientset.KfilterV1alpha1().Transforms(desired.Namespace).Update(existing)
}
