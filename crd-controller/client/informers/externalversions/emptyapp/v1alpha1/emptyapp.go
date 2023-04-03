/*
Copyright The Kubernetes Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	emptyappv1alpha1 "k8s-client-demo/crd-controller/apis/emptyapp/v1alpha1"
	versioned "k8s-client-demo/crd-controller/client/clientset/versioned"
	internalinterfaces "k8s-client-demo/crd-controller/client/informers/externalversions/internalinterfaces"
	v1alpha1 "k8s-client-demo/crd-controller/client/listers/emptyapp/v1alpha1"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// EmptyAppInformer provides access to a shared informer and lister for
// EmptyApps.
type EmptyAppInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.EmptyAppLister
}

type emptyAppInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewEmptyAppInformer constructs a new informer for EmptyApp type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewEmptyAppInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredEmptyAppInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredEmptyAppInformer constructs a new informer for EmptyApp type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredEmptyAppInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CrdV1alpha1().EmptyApps(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CrdV1alpha1().EmptyApps(namespace).Watch(context.TODO(), options)
			},
		},
		&emptyappv1alpha1.EmptyApp{},
		resyncPeriod,
		indexers,
	)
}

func (f *emptyAppInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredEmptyAppInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *emptyAppInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&emptyappv1alpha1.EmptyApp{}, f.defaultInformer)
}

func (f *emptyAppInformer) Lister() v1alpha1.EmptyAppLister {
	return v1alpha1.NewEmptyAppLister(f.Informer().GetIndexer())
}
