// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

	"github.com/isovalent/gke-test-cluster-management/operator/pkg/config"
)

// TestClusterGKEReconciler reconciles a TestClusterGKE object
type TestClusterGKEReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	ConfigRenderer *config.Config
}

const (
	Finalizer = "finalizer.clusters.ci.cilium.io"
)

// +kubebuilder:rbac:groups=clusters.ci.cilium.io,resources=testclustergkes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=clusters.ci.cilium.io,resources=testclustergkes/status,verbs=get;update;patch
func (r *TestClusterGKEReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("Reconcile", req.NamespacedName)

	instance := &clustersv1alpha1.TestClusterGKE{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
	// }

	objs, err := r.RenderObjects(instance)
	if err != nil {
		log.Error(err, "unable render config template")
		return ctrl.Result{}, err
	}
	log.Info("generated config", "items", objs.Items)

	// TODO (mvp)
	// - [ ] detect event type, error on updates
	// - [ ] handle deletion
	// - [x] write a few simple controller tests
	// - [ ] update RBAC configs
	// - [ ] de-kustomize configs
	// TODO (post-mvp)
	// - implement pool object
	// - implement GCP project annotation
	// - implement job runner pod (use sonoboy as PoC)
	// - use random cluster name, instead of same as test object
	// - wait for cluster to get created, update status

	if err := objs.EachListItem(r.createOrSkip); err != nil {
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
	objs, err := r.ConfigRenderer.RenderObjects(instance)
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

func (r *TestClusterGKEReconciler) createOrSkip(obj runtime.Object) error {
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		return err
	}

	ctx := context.Background()
	log := r.Log.WithValues("createOrSkip", key)

	// TODO (post-mvp) probably don't need to make a full copy,
	// should be able to copy just TypeMeta and ObjectMeta
	remoteObj := obj.DeepCopyObject()
	getErr := r.Client.Get(ctx, key, remoteObj)
	if apierrors.IsNotFound(getErr) {
		log.Info("will create", "obj", obj)
		return r.Client.Create(ctx, obj)
	}
	if getErr == nil {
		log.Info("already exists", "remoteObj", remoteObj)
	}
	return getErr
}
