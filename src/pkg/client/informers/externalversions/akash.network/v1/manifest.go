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

package v1

import (
	"context"
	time "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	akashnetworkv1 "providerService/src/pkg/apis/akash.network/v1"
	versioned "providerService/src/pkg/client/clientset/versioned"
	internalinterfaces "providerService/src/pkg/client/informers/externalversions/internalinterfaces"
	v1 "providerService/src/pkg/client/listers/akash.network/v1"
)

// ManifestInformer provides access to a shared informer and lister for
// Manifests.
type ManifestInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ManifestLister
}

type manifestInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewManifestInformer constructs a new informer for Manifest type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewManifestInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredManifestInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredManifestInformer constructs a new informer for Manifest type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredManifestInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AkashV1().Manifests(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AkashV1().Manifests(namespace).Watch(context.TODO(), options)
			},
		},
		&akashnetworkv1.Manifest{},
		resyncPeriod,
		indexers,
	)
}

func (f *manifestInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredManifestInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *manifestInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&akashnetworkv1.Manifest{}, f.defaultInformer)
}

func (f *manifestInformer) Lister() v1.ManifestLister {
	return v1.NewManifestLister(f.Informer().GetIndexer())
}
