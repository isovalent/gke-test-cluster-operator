// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	"github.com/isovalent/gke-test-cluster-management/operator/api/cnrm"
	"github.com/isovalent/gke-test-cluster-management/operator/controllers/common"
	gkeclient "github.com/isovalent/gke-test-cluster-management/operator/pkg/client"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/github"

	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/config"
)

// setup watchers only for functional depednencies of the cluster,
// no watchers are needed for IAM resources

// watch for object, check ownership separately
var cnrmEventHandler = &handler.EnqueueRequestForObject{}

type CNRMContainerClusterWatcher struct {
	common.ClientLogger
	gkeclient.ClientSetBuilder
	Scheme         *runtime.Scheme
	ConfigRenderer *config.Config
}

type CNRMContainerNodePoolSourceWatcher struct {
	common.ClientLogger
	Scheme *runtime.Scheme
}

type CNRMComputeNetworkWatcher struct {
	common.ClientLogger
	Scheme *runtime.Scheme
}

type CNRMComputeSubnetworkWatcher struct {
	common.ClientLogger
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
		if client.IgnoreNotFound(err) != nil {
			w.MetricTracker.Errors.Inc()
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	owner, err := w.GetOwner(ctx, req.NamespacedName, instance.GetOwnerReferences())
	if err != nil {
		return ctrl.Result{}, err
	}
	if owner == nil {
		log.V(1).Info("object not owned by the opertor")
		return ctrl.Result{}, nil
	}

	ghs := github.NewStatusUpdater(w.Log.WithValues("GitHubStatus", req.NamespacedName), owner.ObjectMeta)

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

	if err := w.UpdateOwnerStatus(ctx, status, owner); err != nil {
		log.Error(err, "failed to update owner status")
		w.MetricTracker.Errors.Inc()
		return ctrl.Result{}, err
	}

	if status.HasReadyCondition() {
		if err := w.EnsureTestRunnerJobClusterRoleBindingExists(ctx, instance); err != nil {
			w.MetricTracker.Errors.Inc()
			return ctrl.Result{}, err
		}

		if owner.Spec.JobSpec != nil {
			objs, err := w.RenderObjects(owner)
			if err != nil {
				log.Error(err, "failed to render job objects")
				w.MetricTracker.Errors.Inc()
				return ctrl.Result{}, err
			}
			log.Info("generated job", "items", objs.Items)
			ifCreated := func() {
				w.MetricTracker.JobsCreated.Inc()
				ghs.Update(ctx, github.StatePending, "test job launched")
			}
			if err := w.MaybeCreate(objs, ifCreated); err != nil {
				log.Error(err, "unable reconcile object")
				w.MetricTracker.Errors.Inc()
				return ctrl.Result{}, err
			}

		}
	}

	return ctrl.Result{}, nil
}

const TestRunnerJobClusterRoleBindingName = "test-job-runner"

func (w *CNRMContainerClusterWatcher) EnsureTestRunnerJobClusterRoleBindingExists(ctx context.Context, instance *unstructured.Unstructured) error {

	cluster, err := cnrm.ParsePartialContainerCluster(instance)
	if err != nil {
		return err
	}

	clusterClient, err := w.NewClientSet(cluster)
	if err != nil {
		return err
	}

	project, ok := cluster.Annotations["cnrm.cloud.google.com/project-id"]
	if !ok {
		return fmt.Errorf("unable to get project ID")
	}

	serviceAccountEmail := fmt.Sprintf("%s-admin@%s.iam.gserviceaccount.com", cluster.Name, project)

	crbClient := clusterClient.RbacV1().ClusterRoleBindings()

	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: TestRunnerJobClusterRoleBindingName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
		Subjects: []rbacv1.Subject{{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "User",
			Name:     serviceAccountEmail,
		}},
	}

	_, err = crbClient.Get(ctx, TestRunnerJobClusterRoleBindingName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		if _, err := crbClient.Create(ctx, crb, metav1.CreateOptions{}); err != nil {
			return err
		}
	}

	return nil
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

func (r *CNRMContainerClusterWatcher) RenderObjects(ownerObj *clustersv1alpha1.TestClusterGKE) (*unstructured.UnstructuredList, error) {
	objs, err := r.ConfigRenderer.RenderTestRunnerJobResources(ownerObj)
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
