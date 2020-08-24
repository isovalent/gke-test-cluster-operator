// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/isovalent/gke-test-cluster-management/operator/api/cnrm"
	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

	"github.com/isovalent/gke-test-cluster-management/operator/config/templates/basic"
	"github.com/isovalent/gke-test-cluster-management/operator/controllers"
	"github.com/isovalent/gke-test-cluster-management/operator/controllers/common"
	controllerscommon "github.com/isovalent/gke-test-cluster-management/operator/controllers/common"
	gkeclient "github.com/isovalent/gke-test-cluster-management/operator/pkg/client"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/config"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = clustersv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme

	// register Config Connector scheme, so that the client can access its objects
	_ = cnrm.AddToScheme(scheme)
}

func main() {
	port := flag.Int("port", 9443, "port to listen on")
	metricsAddr := flag.String("metrics-addr", ":8080", "address the metric endpoint binds to")
	enableLeaderElection := flag.Bool("enable-leader-election", false, "enable leader election")
	leaderElectionID := flag.String("leader-election-id", "gke-test-cluster-operator.ci.cilium.io", "identifier to use for leader election")
	logviewDomain := flag.String("logview-domain", "", "domain to use for generating logview url")

	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	configRenderer, err := initConfigRenderer()

	if err != nil {
		setupLog.Error(err, "unable to setup conig and job renderers")
		os.Exit(2)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: *metricsAddr,
		Port:               *port,
		LeaderElection:     *enableLeaderElection,
		LeaderElectionID:   *leaderElectionID,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	metricTracker := controllerscommon.NewMetricTracker()

	if err := (&controllers.TestClusterGKEReconciler{
		ClientLogger:   controllerscommon.NewClientLogger(mgr, ctrl.Log, metricTracker, "TestClusterGKE"),
		Scheme:         mgr.GetScheme(),
		ConfigRenderer: configRenderer,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "TestClusterGKE")
		os.Exit(1)
	}
	if err := (&controllers.TestClusterPoolGKEReconciler{
		ClientLogger: controllerscommon.NewClientLogger(mgr, ctrl.Log, metricTracker, "TestClusterPoolGKE"),
		Scheme:       mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "TestClusterPoolGKE")
		os.Exit(1)
	}
	if err = (&clustersv1alpha1.TestClusterGKE{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "TestClusterGKE")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	clientSetBuilder, err := gkeclient.NewClientSetBuilder()
	if err != nil {
		setupLog.Error(err, "unable to construct clientset builder")
		os.Exit(1)
	}

	if err := (&controllers.CNRMContainerClusterWatcher{
		ClientLogger:     controllerscommon.NewClientLogger(mgr, ctrl.Log, metricTracker, "CNRMContainerClusterWatcher"),
		Scheme:           mgr.GetScheme(),
		ConfigRenderer:   configRenderer,
		ClientSetBuilder: *clientSetBuilder,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CNRMContainerClusterWatcher")
		os.Exit(1)
	}

	if err := (&controllers.JobWatcher{
		ClientLogger: controllerscommon.NewClientLogger(mgr, ctrl.Log, metricTracker, "JobWatcher"),
		Logview:      &common.LogviewService{Domain: *logviewDomain},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "JobWatcher")
		os.Exit(1)
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func initConfigRenderer() (*config.Config, error) {
	cr := &config.Config{
		BaseDirectory: "./config/templates",
	}
	if err := cr.Load(); err != nil {
		return nil, err
	}

	if err := cr.ApplyDefaults("basic", basic.NewDefaults()); err != nil {
		return nil, err
	}

	if err := cr.ApplyDefaultsForClusterAccessResources(basic.NewDefaults()); err != nil {
		return nil, err
	}

	if err := cr.ApplyDefaultsForClusterAccessResources(basic.NewDefaults()); err != nil {
		return nil, err
	}

	return cr, nil
}
