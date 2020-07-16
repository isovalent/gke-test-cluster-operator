// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	"github.com/go-logr/logr"
	"github.com/isovalent/gke-test-cluster-management/operator/api/cnrm"
	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
)

type CNRMWatcher struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (w *CNRMWatcher) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("cnrm-watcher", mgr, controller.Options{
		Reconciler: w,
	})
	if err != nil {
		return err
	}

	// watch for object, check ownership separately
	eventHandler := &handler.EnqueueRequestForObject{}

	if err := c.Watch(cnrm.NewContainerClusterSource(), eventHandler); err != nil {
		return err
	}
	if err := c.Watch(cnrm.NewContainerNodePoolSource(), eventHandler); err != nil {
		return err
	}
	if err := c.Watch(cnrm.NewComputeNetworkSource(), eventHandler); err != nil {
		return err
	}
	if err := c.Watch(cnrm.NewComputeSubnetworkSource(), eventHandler); err != nil {
		return err
	}

	return nil
}

func (w *CNRMWatcher) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := w.Log.WithValues("Reconcile", req.NamespacedName)

	log.V(1).Info("request")

	instance := cnrm.NewContainerCluster()
	if err := w.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	owners := instance.GetOwnerReferences()
	numOwners := len(owners)
	if numOwners == 0 {
		log.V(1).Info("object not owned by the opertor")
		return ctrl.Result{}, nil
	}

	if numOwners > 1 {
		return ctrl.Result{}, fmt.Errorf("object %q has unexpected number of owners (%d)", req.NamespacedName, numOwners)
	}

	ownerObj := &clustersv1alpha1.TestClusterGKE{
		ObjectMeta: metav1.ObjectMeta{
			Name:      owners[0].Name,
			Namespace: req.Namespace,
			UID:       owners[0].UID, // maybe not needed really?
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: owners[0].APIVersion,
			Kind:       owners[0].Kind,
		},
	}
	ownerKey := types.NamespacedName{
		Name:      owners[0].Name,
		Namespace: req.Namespace,
	}

	if instance.GetDeletionTimestamp() != nil {
		log.V(1).Info("deleting")
		return ctrl.Result{}, nil
	}

	status, err := GetContainerClusterStatus(instance)
	if err != nil {
		log.Error(err, "failed to get status")
		return ctrl.Result{}, err
	}

	if status == nil {
		log.V(1).Info("empty status")
		return ctrl.Result{}, nil
	}

	log.V(1).Info("reconciling status", "status", status)

	if err := w.Get(ctx, ownerKey, ownerObj); err != nil {
		return ctrl.Result{}, err
	}

	ownerObj.Status = *status
	log.Info("updating owner status", "owner", ownerKey)
	if err := w.Update(ctx, ownerObj); err != nil {
		return ctrl.Result{}, err
	}
	if status.HasReadyCondition() {
		// TODO: get creds
	}

	// TODO (mvp)
	// - [x] update status of the owner
	// - [ ] fetch credentials and create a secret
	// TODO (post-mvp)
	// - update status in a more sophisticated manner, the transition timestamps should corespond to the time of update
	// - review status structs, using the same struct is probably a naive idea
	// - inspect all of the 4 object and set parent condtion accordingly, so progress is fully trackable

	return ctrl.Result{}, nil
}

type ContainerClusterStatus = clustersv1alpha1.TestClusterGKEStatus
type ContainerClusterStatusCondition = clustersv1alpha1.TestClusterGKEStatusCondition

func GetContainerClusterStatus(instance *unstructured.Unstructured) (*ContainerClusterStatus, error) {
	// TestClusterGKEStatus is really based on CNRM's ContainerClusterStatus,
	// so the same type is used here while it's actually defined as part of
	// the main API
	statusObj, ok := instance.Object["status"]
	if !ok {
		// ignore objects without status,
		// presumably this just hasn't been populated yet
		return nil, nil
	}

	data, err := json.Marshal(statusObj)
	if err != nil {
		return nil, err
	}

	status := &ContainerClusterStatus{}
	if err := json.Unmarshal(data, status); err != nil {
		return nil, err
	}
	return status, nil
}
