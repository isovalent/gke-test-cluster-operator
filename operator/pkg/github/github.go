// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

// based on https://gist.github.com/ahmetb/548059cdbf12fb571e4e2f1e29c48997
package github

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"

	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
)

func GetClient(ctx context.Context) (*github.Client, error) {
	token := os.Getenv("GH_TOKEN")

	if token == "" {
		return nil, fmt.Errorf("GH_TOKEN not set")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}

func UpdateClusterStatus(ctx context.Context, cluster *clustersv1alpha1.TestClusterGKE) error {
	repo, ok := cluster.Annotations["ci.cilium.io/repository-name"]
	if !ok {
		return fmt.Errorf("repository name not present in cluster %s annotations", cluster.Name)
	}
	owner, ok := cluster.Annotations["ci.cilium.io/repository-owner"]
	if !ok {
		return fmt.Errorf("repository owner not present in cluster %s annotations", cluster.Name)
	}
	hash, ok := cluster.Annotations["ci.cilium.io/commit-hash"]
	if !ok {
		return fmt.Errorf("commit hash not present in cluster %s annotations", cluster.Name)
	}

	client, err := GetClient(ctx)
	if err != nil {
		return err
	}

	_, _, err = client.Repositories.CreateStatus(ctx, owner, repo, hash, cluster.GetGithubStatus())
	return err
}
