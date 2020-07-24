// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

	"github.com/isovalent/gke-test-cluster-management/operator/pkg/config"
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

// +kubebuilder:rbac:groups=iam.cnrm.cloud.google.com,resources=iampolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=iam.cnrm.cloud.google.com,resources=iampolicies/status,verbs=get;update;patch

// TestClusterGKEReconciler reconciles a TestClusterGKE object
type TestClusterGKEReconciler struct {
	ClientLogger
	Scheme *runtime.Scheme

	ConfigRenderer *config.Config
}

const (
	Finalizer = "finalizer.clusters.ci.cilium.io"
)

func (r *TestClusterGKEReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("Reconcile", req.NamespacedName)

	log.V(1).Info("request")

	instance := &clustersv1alpha1.TestClusterGKE{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
	// }

	// it's safe to re-generate object, as same name will be used
	objs, err := r.RenderObjects(instance)
	if err != nil {
		log.Error(err, "unable render config template")
		return ctrl.Result{}, err
	}

	log.Info("generated config", "items", objs.Items)

	// update status to store generated cluster name and prevent
	// more clusters being generated
	// (NB: r.RenderObjects sets instance.Status.ClusterName)
	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	// TODO (post-mvp)
	// - capture events
	// - ensure all dependency error are captured and we don't try to re-create
	// - ensure validation and defaulting webhook works, deploy cert-manager
	// - deploy Promethues and for monitoring the operator and configure alerts
	//   (try doing it with stackdriver)
	// - find way to deploy things into the test clusters, maybe use init container
	//   in runner pod for this, probably check standard config as part of cluster
	//   template
	// - move to own namespace, restrict access to configmaps within the namespaces
	// - review RBAC, investigate if resource name prefix can be used for core resources
	// - add RBAC role for CI to submit cluster requests
	// - add developer RBAC role bound to a namespace
	// - register runner job as GitHub Actions runner
	// - consider using a function proxy for request submissions from CI
	// - build image in GitHub Actions
	// - deploy using Flux
	// - ensure updates are handles as intended, i.e. errror
	// - implement pool object
	// - implement GCP project annotation
	// - implement job runner pod (use sonoboy as PoC)

	if err := objs.EachListItem(r.createOrSkip); err != nil {
		// TODO (post-mvp): send these errors as events
		log.Error(err, "unable reconcile object")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
func (r *TestClusterGKEReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clustersv1alpha1.TestClusterGKE{}).
		Complete(r)
}

func (r *TestClusterGKEReconciler) RenderObjects(instance *clustersv1alpha1.TestClusterGKE) (*unstructured.UnstructuredList, error) {
	objs, err := r.ConfigRenderer.RenderObjects(instance, true)
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
