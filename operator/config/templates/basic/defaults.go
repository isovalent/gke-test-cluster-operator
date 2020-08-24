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
			JobSpec: &v1alpha1.TestClusterGKEJobSpec{},
		},
	}

	defaults.Default()

	return defaults
}
