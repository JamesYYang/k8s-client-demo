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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "k8s-client-demo/crd-controller/apis/emptyapp/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// EmptyAppLister helps list EmptyApps.
// All objects returned here must be treated as read-only.
type EmptyAppLister interface {
	// List lists all EmptyApps in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.EmptyApp, err error)
	// EmptyApps returns an object that can list and get EmptyApps.
	EmptyApps(namespace string) EmptyAppNamespaceLister
	EmptyAppListerExpansion
}

// emptyAppLister implements the EmptyAppLister interface.
type emptyAppLister struct {
	indexer cache.Indexer
}

// NewEmptyAppLister returns a new EmptyAppLister.
func NewEmptyAppLister(indexer cache.Indexer) EmptyAppLister {
	return &emptyAppLister{indexer: indexer}
}

// List lists all EmptyApps in the indexer.
func (s *emptyAppLister) List(selector labels.Selector) (ret []*v1alpha1.EmptyApp, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.EmptyApp))
	})
	return ret, err
}

// EmptyApps returns an object that can list and get EmptyApps.
func (s *emptyAppLister) EmptyApps(namespace string) EmptyAppNamespaceLister {
	return emptyAppNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// EmptyAppNamespaceLister helps list and get EmptyApps.
// All objects returned here must be treated as read-only.
type EmptyAppNamespaceLister interface {
	// List lists all EmptyApps in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.EmptyApp, err error)
	// Get retrieves the EmptyApp from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.EmptyApp, error)
	EmptyAppNamespaceListerExpansion
}

// emptyAppNamespaceLister implements the EmptyAppNamespaceLister
// interface.
type emptyAppNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all EmptyApps in the indexer for a given namespace.
func (s emptyAppNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.EmptyApp, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.EmptyApp))
	})
	return ret, err
}

// Get retrieves the EmptyApp from the indexer for a given namespace and name.
func (s emptyAppNamespaceLister) Get(name string) (*v1alpha1.EmptyApp, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("emptyapp"), name)
	}
	return obj.(*v1alpha1.EmptyApp), nil
}
