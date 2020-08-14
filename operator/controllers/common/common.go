// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/go-logr/logr"
	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
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

func (c *ClientLogger) GetOwner(ctx context.Context, objKey types.NamespacedName, ownerRefs []metav1.OwnerReference) (*Owner, error) {
	numOwners := len(ownerRefs)
	if numOwners == 0 {
		return nil, nil
	}

	if numOwners > 1 {
		return nil, fmt.Errorf("object %q unexpected number of owners (%d)", objKey, numOwners)
	}

	owner := ownerRefs[0]

	key := types.NamespacedName{
		Name:      owner.Name,
		Namespace: objKey.Namespace,
	}
	ownerObj := &clustersv1alpha1.TestClusterGKE{}
	c.Get(ctx, key, ownerObj)

	return &Owner{Object: ownerObj}, nil
}

type Owner struct {
	Object *clustersv1alpha1.TestClusterGKE
}

func (o *Owner) Key() types.NamespacedName {
	return types.NamespacedName{
		Name:      o.Object.Name,
		Namespace: o.Object.Namespace,
	}
}

func (o *Owner) UpdateStatus(client client.Client, status *clustersv1alpha1.TestClusterGKEStatus, objKey types.NamespacedName) error {
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
