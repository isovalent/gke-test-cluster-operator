// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/isovalent/gke-test-cluster-management/operator/controllers/common"
)

// watch for object, check ownership separately
var jobEventHandler = &handler.EnqueueRequestForObject{}

type JobWatcher struct {
	common.ClientLogger
}

func (w *JobWatcher) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("gke-test-cluster-job-watcher", mgr, controller.Options{
		Reconciler: w,
	})
	if err != nil {
		return err
	}
	return c.Watch(&source.Kind{Type: &batchv1.Job{}}, jobEventHandler, predicate.ResourceVersionChangedPredicate{})
}

func (w *JobWatcher) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := w.Log.WithValues("Reconcile", req.NamespacedName)

	log.V(1).Info("request")

	instance := &batchv1.Job{}
	if err := w.Get(ctx, req.NamespacedName, instance); err != nil {
		if client.IgnoreNotFound(err) != nil {
			w.MetricTracker.Errors.Inc()
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	owner, err := GetOwner(req.NamespacedName, instance.GetOwnerReferences())
	if err != nil {
		w.MetricTracker.Errors.Inc()
		return ctrl.Result{}, err
	}
	if owner == nil {
		log.V(1).Info("object not owned by the operator")
		return ctrl.Result{}, nil
	}

	if instance.GetDeletionTimestamp() != nil {
		log.V(1).Info("object is being deleted")
		return ctrl.Result{}, nil
	}

	if IsJobCompleted(*instance) {
		cluster := owner.Object
		key, err := client.ObjectKeyFromObject(cluster)
		if err != nil {
			w.MetricTracker.Errors.Inc()
			return ctrl.Result{}, err
		}

		err = w.Client.Get(ctx, key, cluster)
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		if err != nil {
			w.MetricTracker.Errors.Inc()
			return ctrl.Result{}, err
		}
		if cluster.DeletionTimestamp != nil {
			return ctrl.Result{}, nil
		}

		instance.ObjectMeta.OwnerReferences = nil
		err = w.Client.Update(ctx, instance)
		if err != nil {
			log.V(1).Info("failed to remove job owner")
		}

		log.V(1).Info("job is completed, deleting owner")
		err = w.Client.Delete(ctx, cluster)
		if err != nil {
			log.V(1).Info("deletion failed")
			w.MetricTracker.Errors.Inc()
		}
	}

	return ctrl.Result{}, err
}

func IsJobCompleted(job batchv1.Job) bool {
	return job.Status.CompletionTime != nil
}
