// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1alpha2 "github.com/isovalent/gke-test-cluster-operator/api/v1alpha2"
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
	Conditions           v1alpha2.CommonConditions `json:"conditions,omitempty"`
	Endpoint             *string                   `json:"endpoint,omitempty"`
	Operation            *string                   `json:"operation,omitempty"`
	KubeconfigSecretName *string                   `json:"kubeconfigSecretName,omitempty"`
	ClusterName          *string                   `json:"clusterName,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
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
