// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	"github.com/go-logr/logr"
	"github.com/isovalent/gke-test-cluster-management/operator/api/cnrm"
	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/job"
)

// setup watchers only for functional depednencies of the cluster,
// no watchers are needed for IAM resources

// watch for object, check ownership separately
var cnrmEventHandler = &handler.EnqueueRequestForObject{}

type CNRMContainerClusterWatcher struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	JobRenderer *job.Config
}

type CNRMContainerNodePoolSourceWatcher struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type CNRMComputeNetworkWatcher struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type CNRMComputeSubnetworkWatcher struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (w *CNRMContainerClusterWatcher) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("cnrm-containercluster-watcher", mgr, controller.Options{
		Reconciler: w,
	})
	if err != nil {
		return err
	}
	return c.Watch(cnrm.NewContainerClusterSource(), cnrmEventHandler)
}

func (w *CNRMContainerClusterWatcher) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := w.Log.WithValues("Reconcile", req.NamespacedName)

	log.V(1).Info("request")

	instance := cnrm.NewContainerCluster()
	if err := w.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	owner, err := GetOwner(req.NamespacedName, instance.GetOwnerReferences())
	if err != nil {
		return ctrl.Result{}, err
	}
	if owner == nil {
		log.V(1).Info("object not owned by the opertor")
		return ctrl.Result{}, nil
	}

	if instance.GetDeletionTimestamp() != nil {
		log.V(1).Info("object is being deleted")
		return ctrl.Result{}, nil
	}

	status, err := w.GetContainerClusterStatus(instance)
	if err != nil {
		log.Error(err, "failed to get status")
		return ctrl.Result{}, err
	}

	if status == nil {
		log.V(1).Info("empty status")
		return ctrl.Result{}, nil
	}

	log.V(1).Info("reconciling status", "status", status)

	if err := owner.UpdateStatus(w.Client, status, req.NamespacedName); err != nil {
		log.Error(err, "failed to update owner status")
		return ctrl.Result{}, err
	}

	if status.HasReadyCondition() && owner.Object.Spec.JobSpec != nil {
		objs, err := w.RenderObjects(owner.Object, req.Name)
		if err != nil {
			log.Error(err, "failed to render job object")
			return ctrl.Result{}, err
		}
		log.Info("generated job", "items", objs.Items)

		if err := objs.EachListItem(w.createOrSkip); err != nil {
			log.Error(err, "unable reconcile object")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

type ContainerClusterStatus = clustersv1alpha1.TestClusterGKEStatus
type ContainerClusterStatusCondition = clustersv1alpha1.TestClusterGKEStatusCondition

func (*CNRMContainerClusterWatcher) GetContainerClusterStatus(instance *unstructured.Unstructured) (*ContainerClusterStatus, error) {
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

func (r *CNRMContainerClusterWatcher) RenderObjects(ownerObj *clustersv1alpha1.TestClusterGKE, name string) (*unstructured.UnstructuredList, error) {
	objs, err := r.JobRenderer.RenderObjects(ownerObj, name)
	if err != nil {
		return nil, err
	}

	for i := range objs.Items {
		// not using objs.EachListItem here sicne it would require type conversion
		if err := controllerutil.SetControllerReference(ownerObj, &objs.Items[i], r.Scheme); err != nil {
			return nil, err
		}
	}

	return objs, nil
}

func (w *CNRMContainerClusterWatcher) createOrSkip(obj runtime.Object) error {
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		return err
	}

	ctx := context.Background()
	log := w.Log.WithValues("createOrSkip", key)

	remoteObj := obj.DeepCopyObject()
	getErr := w.Client.Get(ctx, key, remoteObj)
	if apierrors.IsNotFound(getErr) {
		log.Info("will create", "obj", obj)
		return w.Client.Create(ctx, obj)
	}
	if getErr == nil {
		log.Info("already exists", "remoteObj", remoteObj)
	}
	return getErr
}

func (w *CNRMContainerNodePoolSourceWatcher) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("cnrm-containernodepool-watcher", mgr, controller.Options{
		Reconciler: w,
	})
	if err != nil {
		return err
	}
	return c.Watch(cnrm.NewContainerNodePoolSource(), cnrmEventHandler)
}

func (w *CNRMContainerNodePoolSourceWatcher) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	log := w.Log.WithValues("Reconcile", req.NamespacedName)
	log.V(1).Info("request")
	return ctrl.Result{}, nil
}

func (w *CNRMComputeNetworkWatcher) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("cnrm-computenetwork-watcher", mgr, controller.Options{
		Reconciler: w,
	})
	if err != nil {
		return err
	}
	return c.Watch(cnrm.NewComputeNetworkSource(), cnrmEventHandler)
}

func (w *CNRMComputeNetworkWatcher) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	log := w.Log.WithValues("Reconcile", req.NamespacedName)
	log.V(1).Info("request")
	return ctrl.Result{}, nil
}

func (w *CNRMComputeSubnetworkWatcher) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("cnrm-computesubnetwork-watcher", mgr, controller.Options{
		Reconciler: w,
	})
	if err != nil {
		return err
	}
	return c.Watch(cnrm.NewComputeSubnetworkSource(), cnrmEventHandler)
}

func (w *CNRMComputeSubnetworkWatcher) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	log := w.Log.WithValues("Reconcile", req.NamespacedName)
	log.V(1).Info("request")
	return ctrl.Result{}, nil
}

type Owner struct {
	Object *clustersv1alpha1.TestClusterGKE
}

func GetOwner(objKey types.NamespacedName, ownerRefs []metav1.OwnerReference) (*Owner, error) {
	numOwners := len(ownerRefs)
	if numOwners == 0 {
		return nil, nil
	}

	if numOwners > 1 {
		return nil, fmt.Errorf("object %q unexpected number of owners (%d)", objKey, numOwners)
	}

	owner := ownerRefs[0]

	ownerObj := &clustersv1alpha1.TestClusterGKE{
		ObjectMeta: metav1.ObjectMeta{
			Name:      owner.Name,
			Namespace: objKey.Namespace,
			UID:       owner.UID, // maybe not needed really?
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: owner.APIVersion,
			Kind:       owner.Kind,
		},
	}

	return &Owner{Object: ownerObj}, nil
}

func (o *Owner) Key() types.NamespacedName {
	return types.NamespacedName{
		Name:      o.Object.Name,
		Namespace: o.Object.Namespace,
	}
}

func (o *Owner) UpdateStatus(client client.Client, status *ContainerClusterStatus, objKey types.NamespacedName) error {
	ctx := context.Background()

	if err := client.Get(ctx, o.Key(), o.Object); err != nil {
		return err
	}

	o.Object.Status.Conditions = status.Conditions
	o.Object.Status.Endpoint = status.Endpoint

	if err := client.Update(ctx, o.Object); err != nil {
		return err
	}

	return nil
}
