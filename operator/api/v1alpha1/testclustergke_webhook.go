// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"errors"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var log = logf.Log.WithName("testclustergke-resource")

func (c *TestClusterGKE) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(c).
		Complete()
}

var _ webhook.Defaulter = &TestClusterGKE{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (c *TestClusterGKE) Default() {
	if c.Name != "" { // avoid loging internal annonymous objects
		log.Info("applying defaults", "name", c.Name, "namespace", c.Namespace)
	}

	if c.Spec.Project == nil {
		c.Spec.Project = new(string)
		*c.Spec.Project = "cilium-ci"
	}

	if c.Spec.Location == nil {
		c.Spec.Location = new(string)
		*c.Spec.Location = "europe-west2-b"
	}

	if c.Spec.Region == nil {
		c.Spec.Region = new(string)
		*c.Spec.Region = "europe-west2"
	}

	if c.Spec.JobSpec != nil {
		if c.Spec.JobSpec.Runner == nil {
			c.Spec.JobSpec.Runner = &TestClusterGKEJobRunnerSpec{}
		}

		if c.Spec.JobSpec.Runner.Image == nil {
			c.Spec.JobSpec.Runner.Image = new(string)
			*c.Spec.JobSpec.Runner.Image = "docker.io/google/cloud-sdk:slim@sha256:a2bade78228faad59a16c36d440f10cfef58a6055cd997d19e258c59c78a409a"
		}

		if c.Spec.JobSpec.Runner.InitImage == nil {
			c.Spec.JobSpec.Runner.InitImage = new(string)
			*c.Spec.JobSpec.Runner.InitImage = "quay.io/isovalent/gke-test-cluster-job-runner-init:28c3b8e6218d145398f78e1343d95b16012fc179"
		}
	}

	if c.Spec.MachineType == nil {
		c.Spec.MachineType = new(string)
		*c.Spec.MachineType = "n1-standard-4"
	}

	if c.Spec.Nodes == nil {
		c.Spec.Nodes = new(int)
		*c.Spec.Nodes = 2
	}
}

var _ webhook.Validator = &TestClusterGKE{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (c *TestClusterGKE) ValidateCreate() error {
	log.Info("validate create", "namespace", c.Namespace, "name", c.Name)
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (c *TestClusterGKE) ValidateUpdate(old runtime.Object) error {
	log.V(1).Info("validate update", "namespace", c.Namespace, "name", c.Name, "new.Spec", c.Spec, "old.Spec", old.(*TestClusterGKE).Spec)
	if !equality.Semantic.DeepEqual(c.Spec, old.(*TestClusterGKE).Spec) {
		return errors.New("spec updates are not supported")
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (c *TestClusterGKE) ValidateDelete() error {
	log.Info("validate delete", "namespace", c.Namespace, "name", c.Name)
	return nil
}
