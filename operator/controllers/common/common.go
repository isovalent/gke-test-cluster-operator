// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
)

type ClientLogger struct {
	client.Client
	Log           logr.Logger
	MetricTracker *MetricTracker
}

type MetricTracker struct {
	ClustersCreated prometheus.Counter
	JobsCreated     prometheus.Counter
	Errors          prometheus.Counter
}

func NewMetricTracker() *MetricTracker {
	t := MetricTracker{
		ClustersCreated: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "gke_test_cluster_operator_clusters_created",
			}),
		JobsCreated: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "gke_test_cluster_operator_jobs_created",
			}),
		Errors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "gke_test_cluster_operator_errors",
			}),
	}

	metrics.Registry.MustRegister(
		t.ClustersCreated,
		t.Errors,
	)

	return &t
}
func NewClientLogger(mgr manager.Manager, l logr.Logger, t *MetricTracker, name string) ClientLogger {
	return ClientLogger{
		Client:        mgr.GetClient(),
		Log:           l.WithName("controllers").WithName(name),
		MetricTracker: t,
	}
}

func (c *ClientLogger) MaybeCreate(list *unstructured.UnstructuredList, createdCallback func()) error {
	count := 0
	for _, item := range list.Items {
		created, err := c.createOrSkip(&item)
		if err != nil {
			return err
		}
		if created {
			count++
		}
	}
	if count == len(list.Items) {
		createdCallback()
	}
	return nil
}

func (c *ClientLogger) createOrSkip(obj runtime.Object) (bool, error) {
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	log := c.Log.WithValues("createOrSkip", key)
	client := c.Client

	remoteObj := obj.DeepCopyObject()
	getErr := client.Get(ctx, key, remoteObj)
	if apierrors.IsNotFound(getErr) {
		log.Info("will create", "obj", obj)
		err := client.Create(ctx, obj)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	if getErr == nil {
		log.Info("already exists", "remoteObj", remoteObj)
	}
	return false, getErr
}
