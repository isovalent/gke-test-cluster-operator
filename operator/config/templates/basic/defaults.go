// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package basic

import (
	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewDefaults() *v1alpha2.TestClusterGKE {
	defaults := &v1alpha2.TestClusterGKE{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
		},
		Spec: v1alpha2.TestClusterGKESpec{
			JobSpec: &v1alpha2.TestClusterGKEJobSpec{},
		},
	}

	defaults.Default()

	return defaults
}
