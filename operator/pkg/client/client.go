// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

// based on https://gist.github.com/ahmetb/548059cdbf12fb571e4e2f1e29c48997
package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/isovalent/gke-test-cluster-management/operator/api/cnrm"
	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
)

func googleScopes() []string {
	return []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email",
	}
}

const (
	googleAuthPlugin = "google" // so that this is different than "gcp" that's already in client-go tree.
)

type ClientSetBuilder interface {
	NewClientSet(cluster *cnrm.PartialContainerCluster) (kubernetes.Interface, error)
}

type GKEClientSetBuilder struct{}

func NewClientSetBuilder() (*GKEClientSetBuilder, error) {
	if err := registerAuthProviderPlugin(); err != nil {
		return nil, err
	}
	return &GKEClientSetBuilder{}, nil
}

func (GKEClientSetBuilder) NewClientSet(cluster *cnrm.PartialContainerCluster) (kubernetes.Interface, error) {
	config, err := newConfig(cluster.Status.Endpoint, cluster.Spec.MasterAuth.ClusterCACertificate)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func NewExternalClient(ctx context.Context, project, clusterName string) (kubernetes.Interface, client.Client, error) {
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}
	if err := cnrm.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}
	if err := v1alpha1.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}

	ts, err := google.DefaultTokenSource(ctx, compute.CloudPlatformScope)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get a token source: %v", err)
	}

	httpClient := oauth2.NewClient(ctx, ts)
	gke, err := container.New(httpClient)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create GKE API client: %v", err)
	}

	clusters, err := gke.Projects.Zones.Clusters.List(project, "-").Context(ctx).Do()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot list cluster in project %q: %v", project, err)
	}

	matches := 0

	var cluster *container.Cluster

	for _, c := range clusters.Clusters {
		if c.Name == clusterName {
			matches++
			cluster = c
		}
	}

	if matches == 0 {
		return nil, nil, fmt.Errorf("cluster %q could not be found", clusterName)
	}

	if matches > 1 {
		return nil, nil, fmt.Errorf("too many clusters using the same name %q (found %d, expected 1)", clusterName, matches)
	}

	if err := registerAuthProviderPlugin(); err != nil {
		return nil, nil, err
	}

	config, err := newConfig(cluster.Endpoint, cluster.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return nil, nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	restClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, nil, err
	}
	return clientSet, restClient, nil
}

var _ rest.AuthProvider = &googleAuthProvider{}

type googleAuthProvider struct {
	tokenSource oauth2.TokenSource
}

func (g *googleAuthProvider) WrapTransport(rt http.RoundTripper) http.RoundTripper {
	return &oauth2.Transport{
		Base:   rt,
		Source: g.tokenSource,
	}
}
func (g *googleAuthProvider) Login() error { return nil }

func newConfig(endpoint, ecodedCACert string) (*rest.Config, error) {
	caCert, err := base64.StdEncoding.DecodeString(ecodedCACert)
	if err != nil {
		return nil, fmt.Errorf("error decoding CA certificate: %v", err)
	}

	config := &rest.Config{
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(caCert),
		},
		Host:         endpoint,
		AuthProvider: &clientcmdapi.AuthProviderConfig{Name: googleAuthPlugin},
	}

	return config, nil
}

func registerAuthProviderPlugin() error {
	newGoogleAuthProvider := func(_ string, _ map[string]string, _ rest.AuthProviderConfigPersister) (rest.AuthProvider, error) {
		ts, err := google.DefaultTokenSource(context.Background(), googleScopes()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create google token source: %+v", err)
		}
		return &googleAuthProvider{tokenSource: ts}, nil
	}
	if err := rest.RegisterAuthProviderPlugin(googleAuthPlugin, newGoogleAuthProvider); err != nil {
		return fmt.Errorf("failed to register %s auth plugin: %v", googleAuthPlugin, err)
	}
	return nil
}
