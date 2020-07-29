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

	if err := (&controllers.TestClusterGKEReconciler{
		ClientLogger:   controllers.NewClientLogger(mgr, ctrl.Log, "TestClusterGKE"),
		Scheme:         mgr.GetScheme(),
		ConfigRenderer: configRenderer,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "TestClusterGKE")
		os.Exit(1)
	}
	if err := (&controllers.TestClusterPoolGKEReconciler{
		ClientLogger: controllers.NewClientLogger(mgr, ctrl.Log, "TestClusterPoolGKE"),
		Scheme:       mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "TestClusterPoolGKE")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	clientSetBuilder, err := gkeclient.NewClientSetBuilder()
	if err != nil {
		setupLog.Error(err, "unable to construct clientset builder")
		os.Exit(1)
	}

	if err := (&controllers.CNRMContainerClusterWatcher{
		ClientLogger:     controllers.NewClientLogger(mgr, ctrl.Log, "CNRMContainerClusterWatcher"),
		Scheme:           mgr.GetScheme(),
		ConfigRenderer:   configRenderer,
		ClientSetBuilder: *clientSetBuilder,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CNRMContainerClusterWatcher")
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
