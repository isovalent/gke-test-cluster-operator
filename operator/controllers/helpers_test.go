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
	"context"
	"flag"
	"fmt"
	"math/rand"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/isovalent/gke-test-cluster-management/operator/api/cnrm"
	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	"github.com/isovalent/gke-test-cluster-management/operator/config/templates/basic"
	"github.com/isovalent/gke-test-cluster-management/operator/controllers"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/config"
)

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	resourcePrefix = flag.String("resource-prefix", fmt.Sprintf("test-%d", rng.Uint64()), "resource prefix")
	crdPath        = flag.String("crd-path", filepath.Join("..", "config", "crd", "bases"), "path to CRDs")
	pollInterval   = flag.Duration("poll-interval", 10*time.Second, "polling interval")
	pollTimeout    = flag.Duration("poll-timeout", 2*time.Minute, "polling timeout")
)

type TLogger struct {
	t *testing.T
}

func (t *TLogger) Write(p []byte) (int, error) {
	t.t.Log(string(p))
	return len(p), nil
}

func setup(t *testing.T) (*ControllerSubTestManager, func()) {
	t.Helper()

	logf.SetLogger(zap.LoggerTo(&TLogger{t: t}, true))

	g := NewGomegaWithT(t)

	var err error

	env := &envtest.Environment{
		CRDDirectoryPaths:     []string{*crdPath},
		UseExistingCluster:    new(bool),
		ErrorIfCRDPathMissing: true,
	}
	*env.UseExistingCluster = true

	cfg, err := env.Start()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(cfg).ToNot(BeNil())

	scheme := runtime.NewScheme()

	err = clientgoscheme.AddToScheme(scheme)
	g.Expect(err).NotTo(HaveOccurred())

	err = clustersv1alpha1.AddToScheme(scheme)
	g.Expect(err).NotTo(HaveOccurred())

	err = cnrm.AddToScheme(scheme)
	g.Expect(err).NotTo(HaveOccurred())

	kubeClient, err := client.New(cfg, client.Options{Scheme: scheme})
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(kubeClient).ToNot(BeNil())

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	g.Expect(err).ToNot(HaveOccurred())

	configRenderer := &config.Config{
		BaseDirectory: "../config/templates",
	}
	err = configRenderer.Load()
	g.Expect(err).ToNot(HaveOccurred())
	configRenderer.ApplyDefaults("basic", basic.NewDefaults())
	g.Expect(err).ToNot(HaveOccurred())

	err = (&controllers.TestClusterGKEReconciler{
		Client:         mgr.GetClient(),
		Log:            ctrl.Log.WithName("controllers").WithName("TestClusterGKE"),
		Scheme:         mgr.GetScheme(),
		ConfigRenderer: configRenderer,
	}).SetupWithManager(mgr)
	g.Expect(err).ToNot(HaveOccurred())

	stop := make(chan struct{})
	go func() {
		err := mgr.Start(stop)
		g.Expect(err).ToNot(HaveOccurred())
	}()

	teardown := func() {
		close(stop)
		err := env.Stop()
		g.Expect(err).ToNot(HaveOccurred())
	}

	return NewControllerSubTestManager(kubeClient, *resourcePrefix), teardown
}

func newTestClusterGKE(namespace, name string) (types.NamespacedName, *clustersv1alpha1.TestClusterGKE) {
	key := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	obj := &clustersv1alpha1.TestClusterGKE{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: clustersv1alpha1.TestClusterGKESpec{},
	}
	return key, obj
}

func newContainerClusterObjs(g *gomega.WithT, namespace, name string) (types.NamespacedName, *unstructured.UnstructuredList) {
	key := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	objs := &unstructured.UnstructuredList{}

	apis := []struct{ kind, apiVersion string }{
		{kind: "ContainerCluster", apiVersion: "container.cnrm.cloud.google.com/v1beta1"},
		{kind: "ComputeNetwork", apiVersion: "compute.cnrm.cloud.google.com/v1beta1"},
		{kind: "ComputeSubnetwork", apiVersion: "compute.cnrm.cloud.google.com/v1beta1"},
	}

	const jsfmt = `
	{
		"kind": %q,
		"apiVersion": %q,
		"metadata": {
		  "namespace": %q,
		  "name": %q,
		  "labels": {
			"cluster": %q
		  }
		},
		"spec": {}
	}`

	for _, api := range apis {
		obj := &unstructured.Unstructured{}
		err := obj.UnmarshalJSON([]byte(fmt.Sprintf(jsfmt, api.kind, api.apiVersion, namespace, name, name)))
		g.Expect(err).ToNot(HaveOccurred())
		objs.Items = append(objs.Items, *obj)
	}

	return key, objs
}

type ControllerSubTestManager struct {
	client          client.Client
	nextObjectID    uint64
	namespacePrefix string
}

type ControllerSubTest struct {
	Client client.Client

	t               *testing.T
	nextObjectID    uint64
	namespacePrefix string
	namespaces      []*corev1.Namespace
}

func NewControllerSubTestManager(client client.Client, namespacePrefix string) *ControllerSubTestManager {
	return &ControllerSubTestManager{
		client:          client,
		namespacePrefix: namespacePrefix,
	}
}
func (cstm *ControllerSubTestManager) NewControllerSubTest(t *testing.T) *ControllerSubTest {
	t.Helper()

	objectID := atomic.AddUint64(&cstm.nextObjectID, 1)
	namespacePrefix := fmt.Sprintf("%s-%d", cstm.namespacePrefix, objectID)

	return &ControllerSubTest{
		t:               t,
		Client:          cstm.client,
		namespacePrefix: namespacePrefix,
	}
}

func (cst *ControllerSubTest) NextNamespace() string {
	objectID := atomic.AddUint64(&cst.nextObjectID, 1)
	name := fmt.Sprintf("%s-%d", cst.namespacePrefix, objectID)
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"test": cst.namespacePrefix,
			},
		},
	}
	ctx := context.Background()
	if err := cst.Client.Create(ctx, ns); err != nil {
		cst.t.Fatalf("failed to create new namespace %q: %s", name, err.Error())
	}
	cst.namespaces = append(cst.namespaces, ns)
	return name
}

func (cst *ControllerSubTest) cleanup() {
	ctx := context.Background()

	for _, ns := range cst.namespaces {
		if err := cst.Client.Delete(ctx, ns); err != nil {
			cst.t.Fatalf("failed to delete namespace %q: %s", ns.Name, err.Error())
		}
	}
}

func (cst *ControllerSubTest) Run(name string, testFunc func(*gomega.WithT, *ControllerSubTest)) {
	cst.t.Run(name, func(t *testing.T) {
		testFunc(NewGomegaWithT(t), cst)
		cst.cleanup()
	})
}
