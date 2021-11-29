// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Important: Run "make misc.generate" to regenerate code after modifying this file

// TestClusterGKESpec defines the desired state of TestClusterGKE
type TestClusterGKESpec struct {
	// Project is the name of GCP project
	Project *string `json:"project,omitempty"`
	// ConfigTemplate is the name of configuration template to use
	ConfigTemplate *string `json:"configTemplate,omitempty"`
	// Region is a GCP region
	Region *string `json:"region,omitempty"`
	// Location is a GCP zone or region
	Location *string `json:"location,omitempty"`
	// MultiZone indicates whether the cluster is meant to span multiple zones
	MultiZone *bool `json:"multZone,omitempty"`
	// KubernetesVersion is the version of Kubernetes to use
	KubernetesVersion *string `json:"kubernetesVersion,omitempty"`
	// JobSpec is the specification of test job
	JobSpec *TestClusterGKEJobSpec `json:"jobSpec,omitempty"`
	// MachineType is the GCP machine type
	MachineType *string `json:"machineType,omitempty"`
	// Nodes is the number of nodes
	Nodes *int `json:"nodes,omitempty"`
}

// TestClusterGKEJobSpec is the specification of test job
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

// TestClusterGKEStatus defines the observed state of TestClusterGKE
type TestClusterGKEStatus struct {
	Conditions CommonConditions `json:"conditions,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:XPreserveUnknownFields
	Dependencies map[string]CommonConditions `json:"dependencyConditions,omitempty"`
	ClusterName  *string                     `json:"clusterName,omitempty"`
}

type (
	CommonCondition struct {
		Type               string      `json:"type"`
		Status             string      `json:"status"`
		LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
		Reason             string      `json:"reason,omitempty"`
		Message            string      `json:"message,omitempty"`
	}
	CommonConditions []CommonCondition
)

func (c *TestClusterGKEStatus) AllDependeciesReady() bool {
	readyDependecies := 0
	for _, dependencyConditions := range c.Dependencies {
		if dependencyConditions.HaveReadyCondition() {
			readyDependecies++
		}
	}
	return len(c.Dependencies) == readyDependecies
}

func (c *TestClusterGKEStatus) HasReadyCondition() bool {
	return c.Conditions.HaveReadyCondition()
}

func (c CommonConditions) HaveReadyCondition() bool {
	if c == nil {
		return false
	}
	for _, condition := range c {
		if condition.Type == "Ready" && condition.Status == "True" {
			return true
		}
	}
	return false
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:singular=testclustergke
// +kubebuilder:plural=testclustersgke
// +kubebuilder:resource:path=testclustersgke,shortName=tcg;tcgke;

// TestClusterGKE is the Schema for the TestClustersGKE API
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
