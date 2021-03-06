// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/prometheus/client_golang/prometheus"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	clustersv1alpha2 "github.com/isovalent/gke-test-cluster-operator/api/v1alpha2"

	"github.com/isovalent/gke-test-cluster-operator/controllers/common"
	"github.com/isovalent/gke-test-cluster-operator/pkg/config"
	"github.com/isovalent/gke-test-cluster-operator/pkg/github"
)

// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=serviceaccounts/status,verbs=get;update;patch

// +kubebuilder:rbac:groups="batch",resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="batch",resources=jobs/status,verbs=get;update;patch

// +kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;create

// +kubebuilder:rbac:groups="",resources=services,verbs=get;create

// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=pods/status,verbs=get

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=clusters.ci.cilium.io,resources=testclustersgke,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=clusters.ci.cilium.io,resources=testclustersgke/status,verbs=get;update;patch

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

	instance := &clustersv1alpha2.TestClusterGKE{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if client.IgnoreNotFound(err) != nil {
			r.MetricTracker.Errors.Inc()
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if instance.GetDeletionTimestamp() != nil {
		log.V(1).Info("object is being deleted")
		return ctrl.Result{}, nil
	}

	ghs := github.NewStatusUpdater(r.Log.WithValues("GitHubStatus", req.NamespacedName), instance.ObjectMeta)

	if instance.Status.ClusterName == nil {
		generatedName := instance.Name + "-" + utilrand.String(5)
		log.V(1).Info("generated new cluster name", "status.clusterName", generatedName)
		instance.Status.ClusterName = &generatedName
		if err := r.Status().Update(ctx, instance); err != nil {
			return ctrl.Result{}, err
		}
	}

	// it's safe to re-generate object, as same name will be used
	log.V(1).Info("regenerating config", "intance", instance)
	objs, err := r.RenderObjects(instance)
	if err != nil {
		errMsg := "unable render config template"
		log.Error(err, errMsg)
		ghs.Update(ctx, github.StateError, "controller error: "+errMsg, "")
		r.MetricTracker.Errors.Inc()
		return ctrl.Result{}, err
	}

	log.Info("generated config", "items", objs.Items)

	ifCreated := func() {
		r.MetricTracker.ClustersCreated.Inc()
		ghs.Update(ctx, github.StatePending, "cluster created", "")
	}
	if err := r.MaybeCreate(objs, ifCreated); err != nil {
		errMsg := "unable to reconcile objects"
		log.Error(err, errMsg)
		ghs.Update(ctx, github.StateError, "controller error: "+errMsg, "")
		r.MetricTracker.Errors.Inc()
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *TestClusterGKEReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clustersv1alpha2.TestClusterGKE{}).
		Complete(r)
}

func (r *TestClusterGKEReconciler) RenderObjects(instance *clustersv1alpha2.TestClusterGKE) (*unstructured.UnstructuredList, error) {
	objs, err := r.ConfigRenderer.RenderAllClusterResources(instance)
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
