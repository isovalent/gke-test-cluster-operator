// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var log = logf.Log.WithName("testclustergke-resource")

func (r *TestClusterGKE) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-clusters-ci-cilium-io-v1alpha1-testclustergke,mutating=true,failurePolicy=fail,groups=clusters.ci.cilium.io,resources=testclustergkes,verbs=create;update,versions=v1alpha1,name=mtestclustergke.kb.io

var _ webhook.Defaulter = &TestClusterGKE{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *TestClusterGKE) Default() {
	log.Info("default", "name", r.Name)

	if r.Spec.Nodes == nil {
		r.Spec.Nodes = new(int)
		*r.Spec.Nodes = 2
	}

	if r.Spec.MachineType == nil {
		r.Spec.MachineType = new(string)
		*r.Spec.MachineType = "n1-standard-4"
	}

	if r.Spec.Project == nil {
		r.Spec.Project = new(string)
		*r.Spec.Project = "cilium-ci"
	}

	if r.Spec.Location == nil {
		r.Spec.Location = new(string)
		*r.Spec.Location = "europe-west2-b"
	}

	if r.Spec.Region == nil {
		r.Spec.Region = new(string)
		*r.Spec.Region = "europe-west2"
	}

	if r.Spec.JobSpec == nil {
		r.Spec.JobSpec = &TestClusterGKEJobSpec{}
	}

	if r.Spec.JobSpec.Runner == nil {
		r.Spec.JobSpec.Runner = &TestClusterGKEJobRunnerSpec{
			Image:     new(string),
			InitImage: new(string),
		}
	}

	if r.Spec.JobSpec.Runner.Image == nil {
		*r.Spec.JobSpec.Runner.Image = "docker.io/google/cloud-sdk:slim@sha256:a2bade78228faad59a16c36d440f10cfef58a6055cd997d19e258c59c78a409a"
	}

	if r.Spec.JobSpec.Runner.InitImage == nil {
		*r.Spec.JobSpec.Runner.InitImage = "quay.io/isovalent/gke-test-cluster-job-runner-init:28c3b8e6218d145398f78e1343d95b16012fc179"
	}
}

// +kubebuilder:webhook:verbs=create;update;delete,path=/validate-clusters-ci-cilium-io-v1alpha1-testclustergke,mutating=false,failurePolicy=fail,groups=clusters.ci.cilium.io,resources=testclustergkes,versions=v1alpha1,name=vtestclustergke.kb.io

var _ webhook.Validator = &TestClusterGKE{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *TestClusterGKE) ValidateCreate() error {
	log.Info("validate create", "namespace", r.Namespace, "name", r.Name)
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *TestClusterGKE) ValidateUpdate(old runtime.Object) error {
	log.Info("validate update", "namespace", r.Namespace, "name", r.Name)
	return errors.New("updates are not supported")
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *TestClusterGKE) ValidateDelete() error {
	log.Info("validate delete", "namespace", r.Namespace, "name", r.Name)
	return nil
}
