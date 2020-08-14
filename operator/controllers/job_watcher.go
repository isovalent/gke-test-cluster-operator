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
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/github"
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

	owner, err := w.GetOwner(ctx, req.NamespacedName, instance.GetOwnerReferences())
	if err != nil {
		w.MetricTracker.Errors.Inc()
		return ctrl.Result{}, err
	}
	if owner == nil {
		log.V(1).Info("object not owned by the operator")
		return ctrl.Result{}, nil
	}

	ghs := github.NewStatusUpdater(w.Log.WithValues("GitHubStatus", req.NamespacedName), owner.Object.ObjectMeta)

	if instance.GetDeletionTimestamp() != nil {
		log.V(1).Info("object is being deleted")
		return ctrl.Result{}, nil
	}

	if IsJobDone(*instance) {
		cluster := owner.Object
		key, err := client.ObjectKeyFromObject(cluster)
		if err != nil {
			w.MetricTracker.Errors.Inc()
			return ctrl.Result{}, err
		}

		if IsJobCompleted(*instance) {
			ghs.Update(ctx, github.StateSuccess, "test job completed")
		} else {
			ghs.Update(ctx, github.StateFailure, "test job failed")
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
			log.V(1).Info("failed to disown job")
			return ctrl.Result{}, err
		}

		log.V(1).Info("job has completed, deleting owner")
		err = w.Client.Delete(ctx, cluster)
		if err != nil {
			log.V(1).Info("deletion failed")
			w.MetricTracker.Errors.Inc()
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func IsJobCompleted(job batchv1.Job) bool {
	return job.Status.CompletionTime != nil
}

func IsJobDone(job batchv1.Job) bool {
	hasCondition := func(t, s, r string) bool {
		for _, condition := range job.Status.Conditions {
			if string(condition.Type) == t && string(condition.Status) == s && condition.Reason == r {
				return true
			}
		}
		return false
	}
	return IsJobCompleted(job) || hasCondition("Failed", "True", "BackoffLimitExceeded")
}
