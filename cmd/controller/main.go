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
	"time"

	"github.com/knative/pkg/configmap"
	"github.com/knative/pkg/controller"
	"github.com/knative/pkg/logging"
	"github.com/knative/pkg/signals"
	servingclientset "github.com/knative/serving/pkg/client/clientset/versioned"
	servinginformers "github.com/knative/serving/pkg/client/informers/externalversions"
	"github.com/knative/serving/pkg/reconciler"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	clientset "github.com/mattmoor/kfilter/pkg/client/clientset/versioned"
	informers "github.com/mattmoor/kfilter/pkg/client/informers/externalversions"
	"github.com/mattmoor/kfilter/pkg/reconciler/kfilter"
	"github.com/mattmoor/kfilter/pkg/reconciler/ktransform"
)

const (
	threadsPerController = 2
)

var (
	masterURL  = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	kubeconfig = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")

	filterImage    = flag.String("filter", "", "The image that implements filtering.")
	transformImage = flag.String("transform", "", "The image that implements transformation.")
)

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	logger := logging.FromContext(context.TODO()).Named("controller")

	if *filterImage == "" {
		logger.Fatal("Error -filter is unspecified.")
	}

	if *transformImage == "" {
		logger.Fatal("Error -transform is unspecified.")
	}

	cfg, err := clientcmd.BuildConfigFromFlags(*masterURL, *kubeconfig)
	if err != nil {
		logger.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logger.Fatalf("Error building kubernetes clientset: %v", err)
	}

	servingClient, err := servingclientset.NewForConfig(cfg)
	if err != nil {
		logger.Fatalf("Error building serving clientset: %v", err)
	}

	kfClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		logger.Fatalf("Error building run clientset: %s", err.Error())
	}

	configMapWatcher := configmap.NewInformedWatcher(kubeClient, "kfilter-system")

	opt := reconciler.Options{
		KubeClientSet:    kubeClient,
		ServingClientSet: servingClient,
		ConfigMapWatcher: configMapWatcher,
		Logger:           logger,
		ResyncPeriod:     10 * time.Hour, // Based on controller-runtime default.
		StopChannel:      stopCh,
	}

	kfInformerFactory := informers.NewSharedInformerFactory(kfClient, opt.ResyncPeriod)
	servingInformerFactory := servinginformers.NewSharedInformerFactory(servingClient, opt.ResyncPeriod)

	// Our shared index informers.
	filterInformer := kfInformerFactory.Kfilter().V1alpha1().Filters()
	transformInformer := kfInformerFactory.Kfilter().V1alpha1().Transforms()
	serviceInformer := servingInformerFactory.Serving().V1alpha1().Services()

	// Add new controllers here.
	controllers := []*controller.Impl{
		kfilter.NewController(
			opt,
			kfClient,
			serviceInformer,
			filterInformer,
			*filterImage,
		),
		ktransform.NewController(
			opt,
			kfClient,
			serviceInformer,
			transformInformer,
			*transformImage,
		),
	}

	go kfInformerFactory.Start(stopCh)
	go servingInformerFactory.Start(stopCh)

	// Wait for the caches to be synced before starting controllers.
	logger.Info("Waiting for informer caches to sync")
	for i, synced := range []cache.InformerSynced{
		serviceInformer.Informer().HasSynced,
		filterInformer.Informer().HasSynced,
		transformInformer.Informer().HasSynced,
	} {
		if ok := cache.WaitForCacheSync(stopCh, synced); !ok {
			logger.Fatalf("failed to wait for cache at index %v to sync", i)
		}
	}

	// Start all of the controllers.
	for _, ctrlr := range controllers {
		go func(ctrlr *controller.Impl) {
			// We don't expect this to return until stop is called,
			// but if it does, propagate it back.
			if err := ctrlr.Run(threadsPerController, stopCh); err != nil {
				logger.Fatalf("Error running controller: %s", err.Error())
			}
		}(ctrlr)
	}

	<-stopCh
}
