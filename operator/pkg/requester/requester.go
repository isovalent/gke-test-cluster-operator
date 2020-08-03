// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package requester

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	gkeclient "github.com/isovalent/gke-test-cluster-management/operator/pkg/client"
)

const (
	DefaultProject           = "cilium-ci"
	DefaultManagementCluster = "management-cluster-0"
)

type TestClusterRequest struct {
	restClient client.Client
	podClient  typedcorev1.PodInterface
	key        types.NamespacedName
}

func NewTestClusterRequest(ctx context.Context, project, managementCluster, namespace, name string) (*TestClusterRequest, error) {
	clientSet, restClient, err := gkeclient.NewExternalClient(ctx, project, managementCluster)
	if err != nil {
		return nil, err
	}

	tcr := &TestClusterRequest{
		key: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
		podClient:  clientSet.CoreV1().Pods(namespace),
		restClient: restClient,
	}
	return tcr, nil
}

func (tcr *TestClusterRequest) CreateTestCluster(ctx context.Context, runnerJobImage string) error {

	err := tcr.restClient.Get(ctx, tcr.key, &v1alpha1.TestClusterGKE{})
	if !apierrors.IsNotFound(err) {
		return fmt.Errorf("cluster %q already exists in namespace %q", tcr.key.Name, tcr.key.Namespace)
	}

	cluster := &v1alpha1.TestClusterGKE{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tcr.key.Name,
			Namespace: tcr.key.Namespace,
		},
		Spec: v1alpha1.TestClusterGKESpec{
			JobSpec: &v1alpha1.TestClusterGKEJobSpec{
				RunnerImage: &runnerJobImage,
			},
		},
	}

	return tcr.restClient.Create(ctx, cluster)
}
