// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package requester

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha2"
	gkeclient "github.com/isovalent/gke-test-cluster-management/operator/pkg/client"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/github"
)

const (
	DefaultProject           = "cilium-ci"
	DefaultManagementCluster = "management-cluster-0"
)

type TestClusterRequest struct {
	restClient        client.Client
	podClient         typedcorev1.PodInterface
	configMapClient   typedcorev1.ConfigMapInterface
	key               types.NamespacedName
	project           string
	configMapName     *string
	fromGitHubActions bool
	cluster           *v1alpha2.TestClusterGKE
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
		configMapClient:   clientSet.CoreV1().ConfigMaps(namespace),
		podClient:         clientSet.CoreV1().Pods(namespace),
		restClient:        restClient,
		fromGitHubActions: os.Getenv("GITHUB_ACTIONS") == "true",
	}
	return tcr, nil
}

func (tcr *TestClusterRequest) CreateRunnerConfigMap(ctx context.Context, initManifestPath string) error {
	initManifestData, err := ioutil.ReadFile(initManifestPath)
	if err != nil {
		return err
	}

	configMapName := tcr.cluster.Name + "-user"
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
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
	tcr.configMapName = &configMapName
	return nil
}

func (tcr *TestClusterRequest) CreateTestCluster(ctx context.Context, configTemplate, description, runnerImage *string, runnerCommand ...string) error {
	err := tcr.restClient.Get(ctx, tcr.key, &v1alpha2.TestClusterGKE{})
	if !apierrors.IsNotFound(err) {
		return fmt.Errorf("cluster %q already exists in namespace %q", tcr.key.Name, tcr.key.Namespace)
	}

	cluster := &v1alpha2.TestClusterGKE{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tcr.key.Name,
			Namespace: tcr.key.Namespace,
		},
		Spec: v1alpha2.TestClusterGKESpec{
			Project: &tcr.project,
		},
	}
	if configTemplate != nil && *configTemplate != "" {
		cluster.Spec.ConfigTemplate = configTemplate
	}

	if runnerImage != nil && *runnerImage != "" {
		cluster.Spec.JobSpec = &v1alpha2.TestClusterGKEJobSpec{
			Runner: &v1alpha2.TestClusterGKEJobRunnerSpec{
				Image:   runnerImage,
				Command: runnerCommand,
			},
		}
	}
	if description != nil && *description != "" {
		cluster.Annotations = map[string]string{
			"ci.cilium.io/cluster-description": *description,
		}
	}

	if tcr.configMapName != nil {
		if cluster.Spec.JobSpec == nil {
			cluster.Spec.JobSpec = &v1alpha2.TestClusterGKEJobSpec{
				Runner: &v1alpha2.TestClusterGKEJobRunnerSpec{},
			}
		}
		cluster.Spec.JobSpec.Runner.ConfigMap = tcr.configMapName
	}

	if tcr.fromGitHubActions {
		event, err := github.ParsePushEvent()
		if err != nil {
			return err
		}
		if event != nil {
			github.SetMetadata(cluster, *event.HeadCommit.ID, *event.Repo.Organization, *event.Repo.Name, "")
		} else {
			// only push events are support, reset this to prevent
			// MaybeSendInitialGitHubStatusUpdate from being called
			tcr.fromGitHubActions = false
		}
	}
	tcr.cluster = cluster
	return tcr.restClient.Create(ctx, cluster)
}

func (tcr *TestClusterRequest) MaybeSendInitialGitHubStatusUpdate(ctx context.Context) error {
	if !tcr.fromGitHubActions {
		return nil
	}
	client, err := github.NewClient(ctx)
	if err != nil {
		return err
	}
	return github.InitalStatusUpdate(ctx, client, tcr.cluster)
}
