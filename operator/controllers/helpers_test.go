/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers_test

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
)

func setup(t *testing.T) (client.Client, *envtest.Environment) {
	g := NewGomegaWithT(t)

	var err error

	env := &envtest.Environment{
		CRDDirectoryPaths:  []string{filepath.Join("..", "config", "crd", "bases")},
		UseExistingCluster: new(bool),
	}
	*env.UseExistingCluster = true

	cfg, err := env.Start()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(cfg).ToNot(BeNil())

	err = clustersv1alpha1.AddToScheme(scheme.Scheme)
	g.Expect(err).NotTo(HaveOccurred())

	err = clustersv1alpha1.AddToScheme(scheme.Scheme)
	g.Expect(err).NotTo(HaveOccurred())

	kubeClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(kubeClient).ToNot(BeNil())

	return kubeClient, env
}

func teardown(t *testing.T, env *envtest.Environment) {
	g := NewGomegaWithT(t)
	err := env.Stop()
	g.Expect(err).ToNot(HaveOccurred())
}
