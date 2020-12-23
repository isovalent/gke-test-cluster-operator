// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-operator/api/v1alpha1"

	"github.com/isovalent/gke-test-cluster-operator/controllers/common"
)

// TestClusterPoolGKEReconciler reconciles a TestClusterPoolGKE object
type TestClusterPoolGKEReconciler struct {
	common.ClientLogger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=clusters.ci.cilium.io,resources=testclusterpoolgkes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=clusters.ci.cilium.io,resources=testclusterpoolgkes/status,verbs=get;update;patch

func (r *TestClusterPoolGKEReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("testclusterpoolgke", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *TestClusterPoolGKEReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clustersv1alpha1.TestClusterPoolGKE{}).
		Complete(r)
}
