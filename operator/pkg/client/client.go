// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

// based on https://gist.github.com/ahmetb/548059cdbf12fb571e4e2f1e29c48997
package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/isovalent/gke-test-cluster-management/operator/api/cnrm"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
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
	newGoogleAuthProvider := func(_ string, _ map[string]string, _ rest.AuthProviderConfigPersister) (rest.AuthProvider, error) {
		ts, err := google.DefaultTokenSource(context.Background(), googleScopes()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create google token source: %+v", err)
		}

		return &googleAuthProvider{tokenSource: ts}, nil
	}
	if err := rest.RegisterAuthProviderPlugin(googleAuthPlugin, newGoogleAuthProvider); err != nil {
		return nil, fmt.Errorf("failed to register %s auth plugin: %v", googleAuthPlugin, err)
	}
	return &GKEClientSetBuilder{}, nil
}

func (GKEClientSetBuilder) NewClientSet(cluster *cnrm.PartialContainerCluster) (kubernetes.Interface, error) {
	caCert, err := base64.StdEncoding.DecodeString(cluster.Spec.MasterAuth.ClusterCACertificate)
	if err != nil {
		return nil, fmt.Errorf("error decoding CA certificate: %v", err)
	}

	config := &rest.Config{
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(caCert),
		},
		Host:         cluster.Status.Endpoint,
		AuthProvider: &clientcmdapi.AuthProviderConfig{Name: googleAuthPlugin},
	}

	return kubernetes.NewForConfig(config)
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
