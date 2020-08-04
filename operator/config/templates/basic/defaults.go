// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package basic

import (
	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewDefaults() *v1alpha1.TestClusterGKE {
	defaults := &v1alpha1.TestClusterGKE{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
		},
		Spec: v1alpha1.TestClusterGKESpec{
			MachineType: new(string),
			Location:    new(string),
			Region:      new(string),
			JobSpec: &v1alpha1.TestClusterGKEJobSpec{
				Runner: &v1alpha1.TestClusterGKEJobRunnerSpec{
					Image: new(string),
				},
			},
		},
	}

	*defaults.Spec.MachineType = "n1-standard-4"
	*defaults.Spec.Location = "europe-west2-b"
	*defaults.Spec.Region = "europe-west2"
	*defaults.Spec.JobSpec.Runner.Image = "docker.io/google/cloud-sdk:slim@sha256:a2bade78228faad59a16c36d440f10cfef58a6055cd997d19e258c59c78a409a"

	return defaults
}
