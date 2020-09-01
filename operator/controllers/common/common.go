// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/go-logr/logr"
	clustersv1alpha2 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha2"
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

func (c *ClientLogger) GetOwner(ctx context.Context, objKey types.NamespacedName, ownerRefs []metav1.OwnerReference) (*clustersv1alpha2.TestClusterGKE, error) {
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
	ownerObj := &clustersv1alpha2.TestClusterGKE{}
	err := c.Get(ctx, key, ownerObj)
	if err != nil {
		return nil, err
	}

	return ownerObj, nil
}

func (c *ClientLogger) UpdateOwnerStatus(ctx context.Context, dependencyKind string, dependencyKey types.NamespacedName, status *clustersv1alpha2.TestClusterGKEStatus, owner *clustersv1alpha2.TestClusterGKE) error {
	if owner.Status.Dependencies == nil {
		owner.Status.Dependencies = map[string]clustersv1alpha2.TestClusterGKEConditions{}
	}
	key := fmt.Sprintf("%s:%s", dependencyKind, dependencyKey.String())
	owner.Status.Dependencies[key] = status.Conditions

	readinessStatus := "False"
	readinessReason := "DependenciesNotReady"
	readinessMessage := "Some dependencies are not ready yet"

	if owner.Status.AllDependeciesReady() {
		readinessStatus = "True"
		readinessReason = "AllDependenciesReady"
		readinessMessage = fmt.Sprintf("All %d dependencies are ready", len(owner.Status.Dependencies))
	}

	owner.Status.Conditions = clustersv1alpha2.TestClusterGKEConditions{{
		Type:               "Ready",
		Status:             readinessStatus,
		LastTransitionTime: metav1.Time{Time: time.Now()},
		Reason:             readinessReason,
		Message:            readinessMessage,
	}}

	c.Log.V(1).Info("updating owner status", "owner", owner)

	if err := c.Status().Update(ctx, owner); err != nil {
		return err
	}

	return nil
}

type LogviewService struct {
	Domain string
}

func (s *LogviewService) AccessURL(ctx context.Context, cl *ClientLogger, job *batchv1.Job) string {
	if s == nil {
		return ""
	}

	logviewCM := &corev1.ConfigMap{}

	key := types.NamespacedName{
		Name:      "gke-test-cluster-logview",
		Namespace: job.Namespace,
	}

	if err := cl.Get(ctx, key, logviewCM); err != nil {
		cl.Log.Error(err, "unable to retrieve configmap", "key", key)
		return ""
	}

	prefix, ok := logviewCM.Data["ingressRoutePrefix"]
	if !ok {
		cl.Log.Info("ingressRoutePrefix not set in configmap")
		return ""
	}

	pods := &corev1.PodList{
		Items: []corev1.Pod{
			corev1.Pod{},
		},
	}

	req, err := labels.NewRequirement("job-name", selection.Equals, []string{job.Name})
	if err != nil {
		cl.Log.Error(err, "unable to create pod selector")
		return ""
	}

	selector := labels.NewSelector().Add(*req)

	listOptions := &client.ListOptions{
		Namespace:     job.Namespace,
		LabelSelector: selector,
	}

	if err := cl.List(ctx, pods, listOptions); err != nil {
		cl.Log.Error(err, "unable to retrieve pods for job")
	}

	if len(pods.Items) == 0 {
		return ""
	}

	pod := pods.Items[0]

	return fmt.Sprintf("https://%s/%s/logs/%s", s.Domain, prefix, pod.Name)
}
