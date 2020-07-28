// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers_test

import (
	"context"
	"flag"
	"math/rand"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/handler"
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

	resourcePrefix = flag.String("resource-prefix", "test-"+utilrand.String(5), "resource prefix")
	pollInterval   = flag.Duration("poll-interval", 5*time.Second, "polling interval")
	pollTimeout    = flag.Duration("poll-timeout", 120*time.Second, "polling timeout")
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
		UseExistingCluster: new(bool),
	}
	*env.UseExistingCluster = true

	cfg, err := env.Start()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(cfg).ToNot(BeNil())

	scheme := runtime.NewScheme()

	g.Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
	g.Expect(clustersv1alpha1.AddToScheme(scheme)).To(Succeed())
	g.Expect(cnrm.AddToScheme(scheme)).To(Succeed())

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
	g.Expect(configRenderer.Load()).To(Succeed())
	g.Expect(configRenderer.ApplyDefaults("basic", basic.NewDefaults())).To(Succeed())
	g.Expect(configRenderer.ApplyDefaultsForClusterAccessResources(basic.NewDefaults())).To(Succeed())
	g.Expect(configRenderer.ApplyDefaultsForTestRunnerJobResources(basic.NewDefaults())).To(Succeed())

	g.Expect((&controllers.TestClusterGKEReconciler{
		ClientLogger:   controllers.NewClientLogger(mgr, ctrl.Log, "TestClusterGKE"),
		Scheme:         mgr.GetScheme(),
		ConfigRenderer: configRenderer,
	}).SetupWithManager(mgr)).To(Succeed())

	g.Expect((&controllers.CNRMContainerClusterWatcher{
		ClientLogger:   controllers.NewClientLogger(mgr, ctrl.Log, "CNRMWatcher"),
		Scheme:         mgr.GetScheme(),
		ConfigRenderer: configRenderer,
	}).SetupWithManager(mgr)).To(Succeed())

	objChan := make(chan *unstructured.Unstructured)
	g.Expect((&TestCNRMContainerClusterWatcher{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		objChan: objChan,
	}).SetupWithManager(mgr)).To(Succeed())

	stop := make(chan struct{})
	go func() { g.Expect(mgr.Start(stop)).To(Succeed()) }()

	teardown := func() {
		logf.SetLogger(logf.NullLogger{}) // attempt to avoid panic "Log in goroutine after TestControllers has completed"

		close(stop)
		g.Expect(env.Stop()).To(Succeed())
	}

	return NewControllerSubTestManager(kubeClient, *resourcePrefix, objChan), teardown
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

type KeyAndObjs struct {
	Key  types.NamespacedName
	Objs *unstructured.UnstructuredList
}

func newEmptyClusterCoreObjs(namespace, name string) *KeyAndObjs {
	key := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	objs := &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			*cnrm.NewContainerCluster(),
			*cnrm.NewContainerNodePool(),
			*cnrm.NewComputeNetwork(),
			*cnrm.NewComputeSubnetwork(),
		},
	}
	return &KeyAndObjs{Key: key, Objs: objs}
}

func newEmptyClusterAccessObjs(namespace, name string) []*KeyAndObjs {
	list := []*KeyAndObjs{
		{
			Key: types.NamespacedName{
				Name:      name + "-admin",
				Namespace: namespace,
			},
			Objs: &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					*cnrm.NewIAMServiceAccount(),
				},
			},
		},
		{
			Key: types.NamespacedName{
				Name:      name + "-workload-identity",
				Namespace: namespace,
			},
			Objs: &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					*cnrm.NewIAMPolicyMember(),
				},
			},
		},
		{
			Key: types.NamespacedName{
				Name:      name + "-cluster-admin",
				Namespace: namespace,
			},
			Objs: &unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{
					*cnrm.NewIAMPolicyMember(),
				},
			},
		},
	}

	return list
}

type TestCNRMContainerClusterWatcher struct {
	client.Client
	Scheme  *runtime.Scheme
	objChan chan *unstructured.Unstructured
}

func (w *TestCNRMContainerClusterWatcher) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("test-cnrm-containercluster-watcher", mgr, controller.Options{
		Reconciler: w,
	})
	if err != nil {
		return err
	}
	return c.Watch(cnrm.NewContainerClusterSource(), &handler.EnqueueRequestForObject{})
}

func (w *TestCNRMContainerClusterWatcher) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	instance := cnrm.NewContainerCluster()
	if err := w.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	w.objChan <- instance
	return ctrl.Result{}, nil
}

type ControllerSubTestManager struct {
	client          client.Client
	namespacePrefix string
	objChan         chan *unstructured.Unstructured
}

type ControllerSubTest struct {
	Client  client.Client
	ObjChan chan *unstructured.Unstructured

	t                          *testing.T
	testLabel, namespacePrefix string
	namespaces                 []*corev1.Namespace
}

func NewControllerSubTestManager(client client.Client, namespacePrefix string, objChan chan *unstructured.Unstructured) *ControllerSubTestManager {
	return &ControllerSubTestManager{
		client:          client,
		namespacePrefix: namespacePrefix,
		objChan:         objChan,
	}
}
func (cstm *ControllerSubTestManager) NewControllerSubTest(t *testing.T) *ControllerSubTest {
	t.Helper()

	return &ControllerSubTest{
		t:               t,
		Client:          cstm.client,
		namespacePrefix: cstm.namespacePrefix,
		ObjChan:         cstm.objChan,
	}
}

func (cst *ControllerSubTest) NextNamespace() string {
	name := cst.namespacePrefix + "-" + utilrand.String(5)
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

func (cst *ControllerSubTest) Run(name string, testFunc func(*WithT, *ControllerSubTest)) {
	cst.t.Run(name, func(t *testing.T) {
		// t.Parallel()
		testFunc(NewGomegaWithT(t), cst)
		cst.cleanup()
	})
}
