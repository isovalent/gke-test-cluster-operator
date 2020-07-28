// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

// Package cnrm contains API Schema definitions for the CNRM API
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
	if err := (&scheme.Builder{
		GroupVersion: schema.GroupVersion{
			Group:   "iam.cnrm.cloud.google.com",
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

func NewContainerClusterList() *unstructured.UnstructuredList {
	objs := &unstructured.UnstructuredList{}
	objs.SetGroupVersionKind(
		schema.GroupVersionKind{
			Kind:    "ContainerClusterList",
			Group:   "container.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	)
	return objs
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

func NewContainerNodePoolList() *unstructured.UnstructuredList {
	objs := &unstructured.UnstructuredList{}
	objs.SetGroupVersionKind(
		schema.GroupVersionKind{
			Kind:    "ContainerNodePoolList",
			Group:   "container.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	)
	return objs
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

func NewComputeNetworkList() *unstructured.UnstructuredList {
	objs := &unstructured.UnstructuredList{}
	objs.SetGroupVersionKind(
		schema.GroupVersionKind{
			Kind:    "ComputeNetworkList",
			Group:   "compute.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	)
	return objs
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

func NewComputeSubnetworkList() *unstructured.UnstructuredList {
	objs := &unstructured.UnstructuredList{}
	objs.SetGroupVersionKind(
		schema.GroupVersionKind{
			Kind:    "ComputeSubnetworkList",
			Group:   "compute.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	)
	return objs
}

func NewIAMServiceAccount() *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(
		schema.GroupVersionKind{
			Kind:    "IAMServiceAccount",
			Group:   "iam.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	)
	return obj
}

func NewIAMServiceAccountSource() source.Source {
	return &source.Kind{Type: NewIAMServiceAccount()}
}

func NewIAMServiceAccountList() *unstructured.UnstructuredList {
	objs := &unstructured.UnstructuredList{}
	objs.SetGroupVersionKind(
		schema.GroupVersionKind{
			Kind:    "IAMServiceAccountList",
			Group:   "ima.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	)
	return objs
}

func NewIAMPolicyMember() *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(
		schema.GroupVersionKind{
			Kind:    "IAMPolicyMember",
			Group:   "iam.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	)
	return obj
}

func NewIAMPolicyMemberSource() source.Source {
	return &source.Kind{Type: NewIAMPolicyMember()}
}

func NewIAMPolicyMemberList() *unstructured.UnstructuredList {
	objs := &unstructured.UnstructuredList{}
	objs.SetGroupVersionKind(
		schema.GroupVersionKind{
			Kind:    "IAMPolicyMember",
			Group:   "ima.cnrm.cloud.google.com",
			Version: "v1beta1",
		},
	)
	return objs
}
