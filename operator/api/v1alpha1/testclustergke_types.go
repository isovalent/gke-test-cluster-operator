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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TestClusterGKESpec defines the desired state of TestClusterGKE
type TestClusterGKESpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ConfigTemplate is the name of configuration template to use
	ConfigTemplate *string `json:"configTemplate,omitempty"`
	// Location is a GCP zone or region
	Location *string `json:"location,omitempty"`
	// KubernetesVersion is the version of Kubernetes to use
	KubernetesVersion *string `json:"kubernetesVersion,omitempty"`
	// JobSpec is the specification of test job
	JobSpec *JobSpec `json:"jobSpec,omitempty"`
}

// JobSpec is the specification of test job
type JobSpec struct {
	// RunnerImage is the image that will drive the tests
	RunnerImage *string `json:"runnerImage,omitempty"`
	// ImagesToTest is a set of application images that will be tested
	ImagesToTest *map[string]string `json:"imagesToTest,omitempty"`
	// RunnerConfigMap is a name of configmap of the runner
	RunnerConfigMap *string `json:"runnerConfigMap,omitempty"`
}

// TestClusterGKEStatus defines the observed state of TestClusterGKE
type TestClusterGKEStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// TestClusterGKE is the Schema for the testclustergkes API
type TestClusterGKE struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestClusterGKESpec   `json:"spec,omitempty"`
	Status TestClusterGKEStatus `json:"status,omitempty"`
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
