// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var testclustergkelog = logf.Log.WithName("testclustergke-resource")

func (r *TestClusterGKE) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-clusters-ci-cilium-io-v1alpha1-testclustergke,mutating=true,failurePolicy=fail,groups=clusters.ci.cilium.io,resources=testclustergkes,verbs=create;update,versions=v1alpha1,name=mtestclustergke.kb.io

var _ webhook.Defaulter = &TestClusterGKE{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *TestClusterGKE) Default() {
	testclustergkelog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-clusters-ci-cilium-io-v1alpha1-testclustergke,mutating=false,failurePolicy=fail,groups=clusters.ci.cilium.io,resources=testclustergkes,versions=v1alpha1,name=vtestclustergke.kb.io

var _ webhook.Validator = &TestClusterGKE{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *TestClusterGKE) ValidateCreate() error {
	testclustergkelog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *TestClusterGKE) ValidateUpdate(old runtime.Object) error {
	testclustergkelog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *TestClusterGKE) ValidateDelete() error {
	testclustergkelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
