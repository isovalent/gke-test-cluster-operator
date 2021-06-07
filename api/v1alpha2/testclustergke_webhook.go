// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	"errors"
	"strings"

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
		log.V(1).Info("defaulting", "namespace", c.Namespace, "name", c.Name, "old.Spec", c.Spec)
	}

	if c.Spec.Project == nil {
		c.Spec.Project = new(string)
		*c.Spec.Project = "cilium-ci"
	}

	if c.Spec.ConfigTemplate == nil {
		c.Spec.ConfigTemplate = new(string)
		*c.Spec.ConfigTemplate = "basic"
	}

	if c.Spec.JobSpec != nil {
		if c.Spec.JobSpec.Runner == nil {
			c.Spec.JobSpec.Runner = &TestClusterGKEJobRunnerSpec{}
		}

		if c.Spec.JobSpec.Runner.Image == nil {
			c.Spec.JobSpec.Runner.Image = new(string)
			*c.Spec.JobSpec.Runner.Image = "quay.io/isovalent/gke-test-cluster-gcloud:803ff83d3786eb38ef05c95768060b0c7ae0fc4d"
		}

		if c.Spec.JobSpec.Runner.InitImage == nil {
			c.Spec.JobSpec.Runner.InitImage = new(string)
			*c.Spec.JobSpec.Runner.InitImage = "quay.io/isovalent/gke-test-cluster-initutil:854733411778d633350adfa1ae66bf11ba658a3f"
		}
	}

	if c.Spec.MultiZone == nil {
		c.Spec.MultiZone = new(bool)
		*c.Spec.MultiZone = false
	}

	if c.Spec.MachineType == nil {
		c.Spec.MachineType = new(string)
		*c.Spec.MachineType = "n1-standard-4"
	}

	if c.Spec.Nodes == nil {
		c.Spec.Nodes = new(int)
		*c.Spec.Nodes = 2
	}

	if c.Name != "" { // avoid loging internal annonymous objects
		log.V(1).Info("defaulting", "namespace", c.Namespace, "name", c.Name, "new.Spec", c.Spec)
	}
}

var _ webhook.Validator = &TestClusterGKE{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (c *TestClusterGKE) ValidateCreate() error {
	log.Info("validate create", "namespace", c.Namespace, "name", c.Name)
	if c.Spec.Region != nil {
		return errors.New("'spec.region' is not user-settable, use 'spec.location' instead")
	}
	if c.Spec.Location != nil {
		numLocationParts := len(strings.Split(*c.Spec.Location, "-"))
		if 2 > numLocationParts || numLocationParts > 3 {
			return errors.New("'spec.location' is invalid format")
		}
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (c *TestClusterGKE) ValidateUpdate(old runtime.Object) error {
	o := old.(*TestClusterGKE)
	log.V(1).Info("validate update", "namespace", c.Namespace, "name", c.Name, "new.Spec", c.Spec, "new.Status", c.Status, "old.Spec", o.Spec, "old.Status", o.Status)
	if !equality.Semantic.DeepEqual(c.Spec, o.Spec) {
		return errors.New("spec updates are not supported")
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (c *TestClusterGKE) ValidateDelete() error {
	log.Info("validate delete", "namespace", c.Namespace, "name", c.Name)
	return nil
}
