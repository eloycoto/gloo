package kube

import (
	"fmt"
	"time"

	"reflect"
	"sync"

	"github.com/pkg/errors"
	"github.com/solo-io/glue/internal/pkg/kube/controller"
	kubev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/listers/core/v1"
	v1beta1listers "k8s.io/client-go/listers/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

type ingressSyncer struct {
	errors chan error

	// name of the kubernetes service for the ingress (envoy)
	ingressService string

	// are we reading all ingress classes, or just glue's?
	globalIngress bool

	client kubernetes.Interface

	ingressLister v1beta1listers.IngressLister
	serviceLister corev1.ServiceLister

	// cache ingress name -> versions so we don't bother updating ingresses in a loop
	cachedStatuses map[string]kubev1.LoadBalancerStatus

	// mutex to protect the map
	mu sync.RWMutex
}

func (c *ingressSyncer) Error() <-chan error {
	return c.errors
}

func NewIngressSyncer(cfg *rest.Config, resyncDuration time.Duration, stopCh <-chan struct{}, globalIngress bool, ingressService string) (*ingressSyncer, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, resyncDuration)
	ingressInformer := kubeInformerFactory.Extensions().V1beta1().Ingresses()
	serviceInformer := kubeInformerFactory.Core().V1().Services()

	c := &ingressSyncer{
		errors:         make(chan error),
		ingressService: ingressService,
		globalIngress:  globalIngress,
		client:         kubeClient,
		ingressLister:  ingressInformer.Lister(),
		serviceLister:  serviceInformer.Lister(),
		cachedStatuses: make(map[string]kubev1.LoadBalancerStatus),
	}

	kubeController := controller.NewController("glue-ingress-syncer", kubeClient,
		c.syncIngressStatus,
		ingressInformer.Informer(),
		serviceInformer.Informer())

	go kubeInformerFactory.Start(stopCh)
	go func() {
		kubeController.Run(2, stopCh)
	}()

	return c, nil
}

func (c *ingressSyncer) syncIngressStatus(namespace, name string, v interface{}) {
	if err := c.sync(); err != nil {
		c.errors <- err
	}
}

func (c *ingressSyncer) sync() error {
	ingresses, err := c.ingressLister.List(labels.Everything())
	if err != nil {
		return errors.Wrap(err, "failed to list ingresses")
	}
	service, err := c.serviceLister.Services(kubev1.NamespaceAll).Get(c.ingressService)
	if err != nil {
		return errors.Wrapf(err, "failed to get ingress service %s", c.ingressService)
	}
	for _, ingress := range ingresses {
		if !isOurIngress(c.globalIngress, ingress) {
			continue
		}
		c.mu.RLock()
		// ignore this ingress if it has the same status from our last update
		if reflect.DeepEqual(c.cachedStatuses[ingress.Name], ingress.Status.LoadBalancer) {
			c.mu.RUnlock()
			continue
		}
		c.mu.RUnlock()
		ingress.Status.LoadBalancer = service.Status.LoadBalancer
		updated, err := c.client.ExtensionsV1beta1().Ingresses(ingress.Namespace).Update(ingress)
		if err != nil {
			runtime.HandleError(errors.Wrap(err, "failed to update ingress with load status"))
		}
		c.mu.Lock()
		c.cachedStatuses[ingress.Name] = updated.Status.LoadBalancer
		c.mu.Unlock()
	}
	return nil
}