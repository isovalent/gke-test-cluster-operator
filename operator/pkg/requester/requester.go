// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package requester

import (
	"context"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
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
	restClient      client.Client
	podClient       typedcorev1.PodInterface
	configMapClient typedcorev1.ConfigMapInterface
	key             types.NamespacedName
	project         string
	hasConfigMap    bool
}

func NewTestClusterRequest(ctx context.Context, project, managementCluster, namespace, name string) (*TestClusterRequest, error) {
	clientSet, restClient, err := gkeclient.NewExternalClient(ctx, project, managementCluster)
	if err != nil {
		return nil, err
	}

	tcr := &TestClusterRequest{
		project: project,
		key: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
		configMapClient: clientSet.CoreV1().ConfigMaps(namespace),
		podClient:       clientSet.CoreV1().Pods(namespace),
		restClient:      restClient,
	}
	return tcr, nil
}

func (tcr *TestClusterRequest) CreateRunnerConfigMap(ctx context.Context, initManifestPath string) error {
	initManifestData, err := ioutil.ReadFile(initManifestPath)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tcr.key.Name,
			Namespace: tcr.key.Namespace,
		},
		BinaryData: map[string][]byte{
			"init-manifest": initManifestData,
		},
	}

	_, err = tcr.configMapClient.Create(ctx, cm, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	tcr.hasConfigMap = true
	return nil
}

func (tcr *TestClusterRequest) CreateTestCluster(ctx context.Context, labels, annotations map[string]string, configTemplate, runnerImage string, runnerCommand ...string) error {

	err := tcr.restClient.Get(ctx, tcr.key, &v1alpha1.TestClusterGKE{})
	if !apierrors.IsNotFound(err) {
		return fmt.Errorf("cluster %q already exists in namespace %q", tcr.key.Name, tcr.key.Namespace)
	}

	cluster := &v1alpha1.TestClusterGKE{
		ObjectMeta: metav1.ObjectMeta{
			Name:        tcr.key.Name,
			Namespace:   tcr.key.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: v1alpha1.TestClusterGKESpec{
			Nodes:          new(int),
			ConfigTemplate: &configTemplate,
			Project:        &tcr.project,
			Location:       new(string),
			Region:         new(string),
			JobSpec: &v1alpha1.TestClusterGKEJobSpec{
				Runner: &v1alpha1.TestClusterGKEJobRunnerSpec{
					Image:     &runnerImage,
					Command:   runnerCommand,
					InitImage: new(string),
				},
			},
		},
	}
	*cluster.Spec.Nodes = 2
	*cluster.Spec.Location = "europe-west2-b"
	*cluster.Spec.Region = "europe-west2"
	*cluster.Spec.JobSpec.Runner.InitImage = "docker.io/errordeveloper/gke-test-cluster-job-runner-init:e8e34968c060a23cfbfb27012d38e5ccbd3e27fe"

	if tcr.hasConfigMap {
		cluster.Spec.JobSpec.Runner.ConfigMap = &tcr.key.Name
	}

	return tcr.restClient.Create(ctx, cluster)
}
