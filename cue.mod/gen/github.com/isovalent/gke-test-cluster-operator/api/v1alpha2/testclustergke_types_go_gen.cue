// Code generated by cue get go. DO NOT EDIT.

//cue:generate cue get go github.com/isovalent/gke-test-cluster-operator/api/v1alpha2

package v1alpha2

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestClusterGKESpec defines the desired state of TestClusterGKE
#TestClusterGKESpec: {
	// Project is the name of GCP project
	project?: null | string @go(Project,*string)

	// ConfigTemplate is the name of configuration template to use
	configTemplate?: null | string @go(ConfigTemplate,*string)

	// Location is a GCP zone or region
	location?: null | string @go(Location,*string)

	// Location is a GCP region (derived from location)
	// TODO: not user-settable, read-only
	region?: null | string @go(Region,*string)

	// KubernetesVersion is the version of Kubernetes to use
	kubernetesVersion?: null | string @go(KubernetesVersion,*string)

	// JobSpec is the specification of test job
	jobSpec?: null | #TestClusterGKEJobSpec @go(JobSpec,*TestClusterGKEJobSpec)

	// MachineType is the GCP machine type
	machineType?: null | string @go(MachineType,*string)

	// Nodes is the number of nodes
	nodes?: null | int @go(Nodes,*int)
}

// TestClusterGKEJobSpec is the specification of test job
#TestClusterGKEJobSpec: {
	// Runner specifies a container that will run control process that drives the tests
	runner?: null | #TestClusterGKEJobRunnerSpec @go(Runner,*TestClusterGKEJobRunnerSpec)

	// ImagesToTest is a set of application images that will be tested
	imagesToTest?: null | {[string]: string} @go(ImagesToTest,*map[string]string)
}

// TestClusterGKEJobRunnerSpec is the specification of test job controll process container
#TestClusterGKEJobRunnerSpec: {
	// Image that will drive the tests
	image?: null | string @go(Image,*string)

	// Command that will be used
	command?: [...string] @go(Command,[]string)

	// InitImage specifies the image used in init container
	initImage?: null | string @go(InitImage,*string)

	// Env speficies environment variables for the runner
	env?: [...corev1.#EnvVar] @go(Env,[]corev1.EnvVar)

	// ConfigMap is a name of configmap of the runner
	configMap?: null | string @go(ConfigMap,*string)
}

// TestClusterGKEStatus defines the observed state of TestClusterGKE
#TestClusterGKEStatus: {
	conditions?: #CommonConditions @go(Conditions)

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:XPreserveUnknownFields
	dependencyConditions?: {[string]: #CommonConditions} @go(Dependencies,map[string]CommonConditions)
	clusterName?: null | string @go(ClusterName,*string)
}

#CommonCondition: {
	type:                string       @go(Type)
	status:              string       @go(Status)
	lastTransitionTime?: metav1.#Time @go(LastTransitionTime)
	reason?:             string       @go(Reason)
	message?:            string       @go(Message)
}

#CommonConditions: [...#CommonCondition]

// TestClusterGKE is the Schema for the TestClustersGKE API
#TestClusterGKE: {
	metav1.#TypeMeta
	metadata?: metav1.#ObjectMeta    @go(ObjectMeta)
	spec?:     #TestClusterGKESpec   @go(Spec)
	status?:   #TestClusterGKEStatus @go(Status)
}

#TestClusterGKE_WithoutTypeMeta: {
	metadata?: metav1.#ObjectMeta    @go(ObjectMeta)
	spec?:     #TestClusterGKESpec   @go(Spec)
	status?:   #TestClusterGKEStatus @go(Status)
}

// TestClusterGKEList contains a list of TestClusterGKE
#TestClusterGKEList: {
	metav1.#TypeMeta
	metadata?: metav1.#ListMeta @go(ListMeta)
	items: [...#TestClusterGKE] @go(Items,[]TestClusterGKE)
}
