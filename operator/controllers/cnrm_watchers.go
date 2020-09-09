// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
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

	clustersv1alpha2 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha2"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/config"
)

const TestRunnerJobClusterRoleBindingName = "test-job-runner"

// setup watchers only for functional depednencies of the cluster,
// no watchers are needed for IAM resources

// watch for object, check ownership separately
var cnrmEventHandler = &handler.EnqueueRequestForObject{}

type CNRMContainerClusterWatcher struct {
	common.ClientLogger
	gkeclient.ClientSetBuilder
	Scheme *runtime.Scheme
}

type CNRMContainerNodePoolWatcher struct {
	common.ClientLogger
	ConfigRenderer *config.Config
	Scheme         *runtime.Scheme
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

	if instance.GetDeletionTimestamp() != nil {
		log.V(1).Info("object is being deleted")
		return ctrl.Result{}, nil
	}

	status, err := cnrm.ParsePartialStatus(instance)
	if err != nil {
		log.Error(err, "failed to get status")
		return ctrl.Result{}, err
	}

	if status == nil {
		log.V(1).Info("empty status")
		return ctrl.Result{}, nil
	}

	log.V(1).Info("reconciling status", "status", status)

	if err := w.UpdateOwnerStatus(ctx, "ContainerCluster", req.NamespacedName, status.Conditions, owner); err != nil {
		log.Error(err, "failed to update owner status")
		w.MetricTracker.Errors.Inc()
		return ctrl.Result{}, err
	}

	if status.HasReadyCondition() {
		if err := w.EnsureTestRunnerJobClusterRoleBindingExists(ctx, instance); err != nil {
			w.MetricTracker.Errors.Inc()
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

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
		w.Log.Info("created ClusterRoleBinding", "name", TestRunnerJobClusterRoleBindingName, "cluster.Name", cluster.Name, "cluster.Namespace", cluster.Namespace, "obj", crb)
		return nil
	}
	w.Log.Info("ClusterRoleBinding already exists", "name", TestRunnerJobClusterRoleBindingName, "cluster.Name", cluster.Name, "cluster.Namespace", cluster.Namespace)
	return nil
}

func (w *CNRMContainerNodePoolWatcher) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("cnrm-containernodepool-watcher", mgr, controller.Options{
		Reconciler: w,
	})
	if err != nil {
		return err
	}
	return c.Watch(cnrm.NewContainerNodePoolSource(), cnrmEventHandler)
}

func (w *CNRMContainerNodePoolWatcher) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := w.Log.WithValues("Reconcile", req.NamespacedName)

	log.V(1).Info("request")

	instance := cnrm.NewContainerNodePool()
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

	status, err := cnrm.ParsePartialStatus(instance)
	if err != nil {
		log.Error(err, "failed to get status")
		return ctrl.Result{}, err
	}

	if status == nil {
		log.V(1).Info("empty status")
		return ctrl.Result{}, nil
	}

	log.V(1).Info("reconciling status", "status", status)

	if err := w.UpdateOwnerStatus(ctx, "ContainerNodePool", req.NamespacedName, status.Conditions, owner); err != nil {
		log.Error(err, "failed to update owner status")
		w.MetricTracker.Errors.Inc()
		return ctrl.Result{}, err
	}

	if status.HasReadyCondition() && owner.Spec.JobSpec != nil {
		objs, err := w.RenderObjects(owner)
		if err != nil {
			log.Error(err, "failed to render job objects")
			w.MetricTracker.Errors.Inc()
			return ctrl.Result{}, err
		}
		log.Info("generated job", "items", objs.Items)
		ifCreated := func() {
			w.MetricTracker.JobsCreated.Inc()
			ghs.Update(ctx, github.StatePending, "test job launched", "")
		}
		if err := w.MaybeCreate(objs, ifCreated); err != nil {
			log.Error(err, "unable reconcile object")
			w.MetricTracker.Errors.Inc()
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *CNRMContainerNodePoolWatcher) RenderObjects(ownerObj *clustersv1alpha2.TestClusterGKE) (*unstructured.UnstructuredList, error) {
	objs, err := r.ConfigRenderer.RenderTestInfraWorkloads(ownerObj)
	if err != nil {
		return nil, err
	}

	// not using objs.EachListItem here since it would require type conversion
	for i := range objs.Items {
		// don't set controller reference for objects in different namespaces e.g. Grafana dashboards
		if objs.Items[i].GetNamespace() == ownerObj.Namespace {
			if err := controllerutil.SetControllerReference(ownerObj, &objs.Items[i], r.Scheme); err != nil {
				return nil, err
			}
		}
	}

	return objs, nil
}
