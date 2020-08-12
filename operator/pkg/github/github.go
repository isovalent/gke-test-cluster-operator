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

func NewClient(ctx context.Context) (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN not set")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}

const (
	metadataKeyPrefix   = "ci.cilium.io/github-"
	annotationRepoOwner = metadataKeyPrefix + "repo-owner"
	annotationRepoName  = metadataKeyPrefix + "repo-name"
	labelCommitHash     = metadataKeyPrefix + "commit-hash"
)

func UpdateClusterStatus(ctx context.Context, cluster *clustersv1alpha1.TestClusterGKE) error {
	commitHash, ok := cluster.Labels[labelCommitHash]
	if !ok {
		return nil
	}

	name, ok := cluster.Annotations[annotationRepoOwner]
	if !ok {
		return fmt.Errorf("annotations %q is not set", annotationRepoOwner)
	}
	owner, ok := cluster.Annotations[annotationRepoName]
	if !ok {
		return fmt.Errorf("annotations %q is not set", annotationRepoName)
	}

	client, err := NewClient(ctx)
	if err != nil {
		return err
	}

	_, _, err = client.Repositories.CreateStatus(ctx, owner, name, commitHash, cluster.GetGithubStatus())
	return err
}
