// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/prometheus/client_golang/prometheus"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

	"github.com/isovalent/gke-test-cluster-management/operator/controllers/common"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/config"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/github"
)

// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=serviceaccounts/status,verbs=get;update;patch

// +kubebuilder:rbac:groups="batch",resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="batch",resources=jobs/status,verbs=get;update;patch

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=clusters.ci.cilium.io,resources=testclustergkes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=clusters.ci.cilium.io,resources=testclustergkes/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=container.cnrm.cloud.google.com,resources=containerclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=container.cnrm.cloud.google.com,resources=containerclusters/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=container.cnrm.cloud.google.com,resources=containernodepools,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=container.cnrm.cloud.google.com,resources=containernodepools/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=compute.cnrm.cloud.google.com,resources=computenetworks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=compute.cnrm.cloud.google.com,resources=computenetworks/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=compute.cnrm.cloud.google.com,resources=computesubnetworks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=compute.cnrm.cloud.google.com,resources=computesubnetworks/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=iam.cnrm.cloud.google.com,resources=iamserviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=iam.cnrm.cloud.google.com,resources=iamserviceaccounts/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=iam.cnrm.cloud.google.com,resources=iampolicymembers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=iam.cnrm.cloud.google.com,resources=iampolicymembers/status,verbs=get;update;patch

// TestClusterGKEReconciler reconciles a TestClusterGKE object
type TestClusterGKEReconciler struct {
	common.ClientLogger
	Scheme *runtime.Scheme

	ConfigRenderer *config.Config
	Metrics        TestClusterGKEReconcilerMetrics
}

// TestClusterGKEReconcilerMetrics contains metrics for TestClusterGKEReconciler
type TestClusterGKEReconcilerMetrics struct {
	ClustersCreatedMetric prometheus.Counter
	ClusterErrorMetric    prometheus.Counter
}

func (r *TestClusterGKEReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("Reconcile", req.NamespacedName)

	log.V(1).Info("request")

	instance := &clustersv1alpha1.TestClusterGKE{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if client.IgnoreNotFound(err) != nil {
			r.MetricTracker.Errors.Inc()
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// it's safe to re-generate object, as same name will be used
	objs, err := r.RenderObjects(instance)
	if err != nil {
		log.Error(err, "unable render config template")
		r.MetricTracker.Errors.Inc()
		return ctrl.Result{}, err
	}

	log.Info("generated config", "items", objs.Items)

	// update status to store generated cluster name and prevent
	// more clusters being generated
	// (NB: r.RenderObjects sets instance.Status.ClusterName)
	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	err, created := r.MaybeCreate(objs, r.MetricTracker.ClustersCreated)
	if err != nil {
		log.Error(err, "unable reconcile object")
		r.MetricTracker.Errors.Inc()
		return ctrl.Result{}, err
	}
	if created {
		err = github.UpdateClusterStatus(ctx, r.Log.WithValues("GitHubStatus", req.NamespacedName), instance)
		if err != nil {
			log.Error(err, "unable to update github status")
			r.MetricTracker.Errors.Inc()
		}
	}

	return ctrl.Result{}, nil
}

func (r *TestClusterGKEReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clustersv1alpha1.TestClusterGKE{}).
		Complete(r)
}

func (r *TestClusterGKEReconciler) RenderObjects(instance *clustersv1alpha1.TestClusterGKE) (*unstructured.UnstructuredList, error) {
	objs, err := r.ConfigRenderer.RenderAllClusterResources(instance, true)
	if err != nil {
		return nil, err
	}

	for i := range objs.Items {
		// not using objs.EachListItem here sicne it would require type conversion
		if err := controllerutil.SetControllerReference(instance, &objs.Items[i], r.Scheme); err != nil {
			return nil, err
		}
	}

	return objs, nil
}
