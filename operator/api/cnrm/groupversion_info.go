/*


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

// Package  API Schema definitions for the CNRM API

package cnrm

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func AddToScheme(s *runtime.Scheme) error {
	if err := (&scheme.Builder{
		GroupVersion: schema.GroupVersion{
			Group:   "container.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	}).AddToScheme(s); err != nil {
		return err
	}
	if err := (&scheme.Builder{
		GroupVersion: schema.GroupVersion{
			Group:   "compute.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	}).AddToScheme(s); err != nil {
		return err
	}
	return nil
}

func NewContainerCluster() *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(
		schema.GroupVersionKind{
			Kind:    "ContainerCluster",
			Group:   "container.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	)
	return obj
}

func NewContainerClusterSource() source.Source {
	return &source.Kind{Type: NewContainerCluster()}
}

func NewContainerNodePool() *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(
		schema.GroupVersionKind{
			Kind:    "ContainerNodePool",
			Group:   "container.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	)
	return obj
}

func NewContainerNodePoolSource() source.Source {
	return &source.Kind{Type: NewContainerNodePool()}
}

func NewComputeNetwork() *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(
		schema.GroupVersionKind{
			Kind:    "ComputeNetwork",
			Group:   "compute.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	)
	return obj
}

func NewComputeNetworkSource() source.Source {
	return &source.Kind{Type: NewComputeNetwork()}
}

func NewComputeSubnetwork() *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(
		schema.GroupVersionKind{
			Kind:    "ComputeSubnetwork",
			Group:   "compute.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	)
	return obj
}

func NewComputeSubnetworkSource() source.Source {
	return &source.Kind{Type: NewComputeSubnetwork()}
}
