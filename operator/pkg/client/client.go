// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

// based on https://gist.github.com/ahmetb/548059cdbf12fb571e4e2f1e29c48997
package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/isovalent/gke-test-cluster-management/operator/api/cnrm"
	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha2"
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

// NewExternalClient will return a ClientSet along with a REST client for the given management cluster,
// it expect that given cluster to be present in exactly one GCP location
func NewExternalClient(ctx context.Context, project, clusterName string) (kubernetes.Interface, client.Client, error) {
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}
	if err := cnrm.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}
	if err := v1alpha2.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}

	opts := []option.ClientOption{}

	creds, err := maybeGetCredenialsFromJSON(ctx)
	if err != nil {
		return nil, nil, err
	}

	if creds != nil {
		opts = append(opts, option.WithCredentials(creds))
	}

	gke, err := container.NewService(ctx, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create GKE API client: %v", err)
	}

	clusters, err := gke.Projects.Zones.Clusters.List(project, "-").Context(ctx).Do()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot list cluster in project %q: %v", project, err)
	}

	matches := 0

	var cluster *container.Cluster

	// gke.Projects.Zones.Clusters.List returns a list of clusters for all regions and zones,
	// which is handy, and prevents having to supply the exact location (and whether it is a
	// region or just a zone)
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

func maybeGetCredenialsFromJSON(ctx context.Context) (*google.Credentials, error) {
	if serviceAccountKey := os.Getenv("GCP_SERVICE_ACCOUNT_KEY"); serviceAccountKey != "" {
		credsData, err := base64.StdEncoding.DecodeString(serviceAccountKey)
		if err != nil {
			return nil, fmt.Errorf("error decoding GCP_SERVICE_ACCOUNT_KEY: %v", err)
		}

		creds, err := google.CredentialsFromJSON(ctx, []byte(credsData), compute.CloudPlatformScope)
		if err != nil {
			return nil, fmt.Errorf("error loading credentials: %v", err)
		}
		return creds, nil
	}
	return nil, nil
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
	if endpoint == "" || ecodedCACert == "" {
		return nil, fmt.Errorf("unexpected empty cluster credentials")
	}
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
		ctx := context.Background()
		creds, err := maybeGetCredenialsFromJSON(ctx)
		if err != nil {
			return nil, err
		}
		if creds != nil {
			return &googleAuthProvider{tokenSource: creds.TokenSource}, nil
		}
		ts, err := google.DefaultTokenSource(ctx, googleScopes()...)
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
