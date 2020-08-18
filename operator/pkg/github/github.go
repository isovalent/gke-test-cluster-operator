// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

// based on https://gist.github.com/ahmetb/548059cdbf12fb571e4e2f1e29c48997
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"

	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	metadataKeyPrefix   = "ci.cilium.io/github-"
	labelCommitHash     = metadataKeyPrefix + "commit-hash"
	annotationRepoOwner = metadataKeyPrefix + "repo-owner"
	annotationRepoName  = metadataKeyPrefix + "repo-name"
	annotationContext   = metadataKeyPrefix + "context"

	StateError   State = "error"
	StateFailure State = "failure"
	StatePending State = "pending"
	StateSuccess State = "success"
)

type State string

type StatusUpdater struct {
	log  logr.Logger
	meta metav1.ObjectMeta

	commitHash, name, owner, context string
}

// NewStatusUpdater constructs an updater to be used in controller context,
// it only logs errors, so doesn't introduce reconciliation failures
func NewStatusUpdater(log logr.Logger, meta metav1.ObjectMeta) *StatusUpdater {
	commitHash, ok := meta.Labels[labelCommitHash]
	if !ok {
		log.Info("will not update GitHub status", "missingLabel", labelCommitHash)
		return nil
	}

	annotationErr := fmt.Errorf("missing annotation")

	owner, ok := meta.Annotations[annotationRepoOwner]
	if !ok {
		log.Error(annotationErr, "missingAnnotations", annotationRepoOwner)
		return nil
	}
	name, ok := meta.Annotations[annotationRepoName]
	if !ok {
		log.Error(annotationErr, "missingAnnotations", annotationRepoName)
		return nil
	}

	context, ok := meta.Annotations[annotationContext]
	if !ok {
		context = fmt.Sprintf("gke-test-cluster-operator:%s/%s", meta.Namespace, meta.Name)
	}

	return &StatusUpdater{
		log:        log,
		meta:       meta,
		commitHash: commitHash,
		name:       name,
		owner:      owner,
		context:    context,
	}
}

func (s *StatusUpdater) Update(ctx context.Context, state State, description, url string) {
	if s == nil {
		return
	}

	client, err := NewClient(ctx)
	if err != nil {
		s.log.Error(err, "unable to create GitHub client")
		return
	}

	status := &github.RepoStatus{
		State:       new(string),
		Description: &description,
		Context:     &s.context,
		URL:         &url,
	}
	*status.State = string(state)

	currentCombinedStatus, _, err := client.Repositories.GetCombinedStatus(ctx, s.owner, s.name, s.commitHash, nil)
	if err != nil {
		s.log.Error(err, "unable to get current GitHub status", "repo", fmt.Sprintf("%s/%s", s.owner, s.name), "ref", s.commitHash)
		return
	}
	needsUpdate := false
	for _, currentStatus := range currentCombinedStatus.Statuses {
		if *currentStatus.Context == s.context {
			if *currentStatus.State == string(StateError) || *currentStatus.State == string(StatePending) {
				needsUpdate = true
			}
		}
	}
	if !needsUpdate {
		s.log.V(1).Info("GitHub status already up-to-date")
		return
	}
	_, _, err = client.Repositories.CreateStatus(ctx, s.owner, s.name, s.commitHash, status)
	if err != nil {
		s.log.Error(err, "unable to update GitHub status", "repo", fmt.Sprintf("%s/%s", s.owner, s.name), "ref", s.commitHash)
		return
	}
	return
}

func SetMetadata(cluster *clustersv1alpha1.TestClusterGKE, commitHash, repoOwner, repoName, context string) {
	if cluster.Labels == nil {
		cluster.Labels = map[string]string{}
	}
	if cluster.Annotations == nil {
		cluster.Annotations = map[string]string{}
	}
	cluster.Labels[labelCommitHash] = commitHash
	cluster.Annotations[annotationRepoOwner] = repoOwner
	cluster.Annotations[annotationRepoName] = repoName
	if context != "" {
		cluster.Annotations[annotationContext] = context
	}
}

// InitalStatusUpdate assumes that SetMetdata was called and makes a direct update to GitHub status API, which
// may result in an error
func InitalStatusUpdate(ctx context.Context, client *github.Client, cluster *clustersv1alpha1.TestClusterGKE) error {
	commitHash := cluster.Labels[labelCommitHash]
	owner := cluster.Annotations[annotationRepoOwner]
	name := cluster.Annotations[annotationRepoName]

	context, ok := cluster.Annotations[annotationContext]
	if !ok {
		context = fmt.Sprintf("gke-test-cluster-operator:%s/%s", cluster.Namespace, cluster.Name)
	}

	status := &github.RepoStatus{
		State:       new(string),
		Description: new(string),
		Context:     &context,
	}
	*status.Description = "creating test cluster"
	*status.State = string(StatePending)

	if _, _, err := client.Repositories.CreateStatus(ctx, owner, name, commitHash, status); err != nil {
		return err
	}
	return nil
}

func ParsePushEvent() (*github.PushEvent, error) {
	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	if eventPath == "" {
		return nil, fmt.Errorf("GITHUB_EVENT_PATH must be set")
	}

	eventName := os.Getenv("GITHUB_EVENT_NAME")
	if eventName == "" {
		return nil, fmt.Errorf("GITHUB_EVENT_NAME must be set")
	}

	if eventName != "push" {
		return nil, nil
	}

	eventData, err := ioutil.ReadFile(eventPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read event data: %v", err)
	}

	event := &github.PushEvent{}

	err = json.Unmarshal(eventData, event)
	if err != nil {
		return nil, fmt.Errorf("cannot parse event data: %v", err)
	}
	return event, nil
}

func NewClient(ctx context.Context) (*github.Client, error) {
	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN must be set")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	client := github.NewClient(oauth2.NewClient(ctx, ts))

	return client, nil
}
