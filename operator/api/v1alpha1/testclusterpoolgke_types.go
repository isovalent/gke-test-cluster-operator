// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TestClusterPoolGKESpec defines the desired state of TestClusterPoolGKE
type TestClusterPoolGKESpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of TestClusterPoolGKE. Edit TestClusterPoolGKE_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// TestClusterPoolGKEStatus defines the observed state of TestClusterPoolGKE
type TestClusterPoolGKEStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// TestClusterPoolGKE is the Schema for the testclusterpoolgkes API
type TestClusterPoolGKE struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestClusterPoolGKESpec   `json:"spec,omitempty"`
	Status TestClusterPoolGKEStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TestClusterPoolGKEList contains a list of TestClusterPoolGKE
type TestClusterPoolGKEList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TestClusterPoolGKE `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TestClusterPoolGKE{}, &TestClusterPoolGKEList{})
}
