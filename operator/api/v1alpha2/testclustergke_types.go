// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TestClusterGKESpec defines the desired state of TestClusterGKE
type TestClusterGKESpec struct {
	// Important: Run "make misc.generate" to regenerate code after modifying this file

	// Project is the name of GCP project
	Project *string `json:"project,omitempty"`
	// ConfigTemplate is the name of configuration template to use
	ConfigTemplate *string `json:"configTemplate,omitempty"`
	// Location is a GCP zone or region
	Location *string `json:"location,omitempty"`
	// Location is a GCP region (derived from location)
	// TODO: not user-settable, read-only
	Region *string `json:"region,omitempty"`
	// KubernetesVersion is the version of Kubernetes to use
	KubernetesVersion *string `json:"kubernetesVersion,omitempty"`
	// JobSpec is the specification of test job
	JobSpec *TestClusterGKEJobSpec `json:"jobSpec,omitempty"`
	// MachineType is the GCP machine type
	MachineType *string `json:"machineType,omitempty"`
	// Nodes is the number of nodes
	Nodes *int `json:"nodes,omitempty"`
}

// JobSpec is the specification of test job
type TestClusterGKEJobSpec struct {
	// Runner specifies a container that will run control process that drives the tests
	Runner *TestClusterGKEJobRunnerSpec `json:"runner,omitempty"`
	// ImagesToTest is a set of application images that will be tested
	ImagesToTest *map[string]string `json:"imagesToTest,omitempty"`
}

// TestClusterGKEJobRunnerSpec is the specification of test job controll process container
type TestClusterGKEJobRunnerSpec struct {
	// Image that will drive the tests
	Image *string `json:"image,omitempty"`
	// Command that will be used
	Command []string `json:"command,omitempty"`
	// InitImage specifies the image used in init container
	InitImage *string `json:"initImage,omitempty"`
	// Env speficies environment variables for the runner
	Env []corev1.EnvVar `json:"env,omitempty"`
	// ConfigMap is a name of configmap of the runner
	ConfigMap *string `json:"configMap,omitempty"`
}

type TestClusterGKEConditions []TestClusterGKECondition

// TestClusterGKEStatus defines the observed state of TestClusterGKE
// +kubebuilder:subresource:status
type TestClusterGKEStatus struct {
	// Important: Run "make misc.generate" to regenerate code after modifying this file
	Conditions   TestClusterGKEConditions            `json:"conditions,omitempty"`
	Dependencies map[string]TestClusterGKEConditions `json:"dependencyConditions,omitempty"`
	ClusterName  *string                             `json:"clusterName,omitempty"`
}

type TestClusterGKECondition struct {
	Type               string      `json:"type"`
	Status             string      `json:"status"`
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
}

func (c *TestClusterGKEStatus) AllDependeciesReady() bool {
	readyDependecies := 0
	for _, dependencyConditions := range c.Dependencies {
		isReady := false
		for _, condition := range dependencyConditions {
			if condition.Type == "Ready" && condition.Status == "True" {
				isReady = true
			}
		}
		if isReady {
			readyDependecies++
		}
	}
	return len(c.Dependencies) == readyDependecies
}

func (c *TestClusterGKEStatus) HasReadyCondition() bool {
	if c == nil {
		return false
	}
	for _, condition := range c.Conditions {
		if condition.Type == "Ready" && condition.Status == "True" {
			return true
		}
	}
	return false
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// TestClusterGKE is the Schema for the testclustergkes API
type TestClusterGKE struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestClusterGKESpec   `json:"spec,omitempty"`
	Status TestClusterGKEStatus `json:"status,omitempty"`
}

type TestClusterGKE_WithoutTypeMeta struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestClusterGKESpec   `json:"spec,omitempty"`
	Status TestClusterGKEStatus `json:"status,omitempty"`
}

// WithoutTypeMeta returns a copy of t without TypeMeta, this is a workaround
// an issue with CUE (see https://github.com/cuelang/cue/discussions/439)
func (t *TestClusterGKE) WithoutTypeMeta() *TestClusterGKE_WithoutTypeMeta {
	return &TestClusterGKE_WithoutTypeMeta{
		ObjectMeta: t.ObjectMeta,
		Spec:       t.Spec,
		Status:     t.Status,
	}
}

// +kubebuilder:object:root=true

// TestClusterGKEList contains a list of TestClusterGKE
type TestClusterGKEList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TestClusterGKE `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TestClusterGKE{}, &TestClusterGKEList{})
}
