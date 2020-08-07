// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/gomega"

	"github.com/isovalent/gke-test-cluster-management/operator/api/cnrm"
	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	. "github.com/isovalent/gke-test-cluster-management/operator/controllers"
)

func TestControllers(t *testing.T) {
	cstm, teardown := setup(t)

	cstm.NewControllerSubTest(t).
		Run("create and delete objects", simpleCreateDeleteObjects)

	cstm.NewControllerSubTest(t).
		Run("create and delete cluster with status updates", createDeleteClusterWithStatusUpdates)

	teardown()
}

func simpleCreateDeleteObjects(g *WithT, cst *ControllerSubTest) {
	initialClusterCount := getMetricIntValue(cst.MetricTracker.ClustersCreated)
	ctx := context.Background()
	ns := cst.NextNamespace()

	key, obj := newTestClusterGKE(ns, "test-1")
	remoteObj := obj.DeepCopy()

	obj.Spec.ConfigTemplate = new(string)
	*obj.Spec.ConfigTemplate = "basic"

	err := cst.Client.Get(ctx, key, remoteObj)
	g.Expect(err).To(HaveOccurred())
	g.Expect(apierrors.IsNotFound(err)).To(BeTrue())

	err = cst.Client.Create(ctx, obj)
	g.Expect(err).ToNot(HaveOccurred())

	err = cst.Client.Get(ctx, key, remoteObj)
	g.Expect(err).ToNot(HaveOccurred())

	listOpts := &client.ListOptions{
		Namespace: ns,
	}
	g.Eventually(func() bool {
		cnrmObjs := cnrm.NewContainerClusterList()
		err := cst.Client.List(ctx, cnrmObjs, listOpts)
		return err == nil && len(cnrmObjs.Items) > 0
	}, *pollTimeout, *pollInterval).Should(BeTrue())

	cnrmContainerClusterObjs := cnrm.NewContainerClusterList()
	g.Expect(cst.Client.List(ctx, cnrmContainerClusterObjs, listOpts)).To(Succeed())
	g.Expect(cnrmContainerClusterObjs.Items).To(HaveLen(1))

	clusterName := cnrmContainerClusterObjs.Items[0].GetName()
	g.Expect(clusterName).To(HavePrefix("test-1-"))
	g.Expect(clusterName).To(HaveLen(12))

	cnrmClusterCore := newEmptyClusterCoreObjs(ns, clusterName)

	g.Expect(cnrmClusterCore.Objs.EachListItem(func(obj runtime.Object) error {
		return cst.Client.Get(ctx, cnrmClusterCore.Key, obj)
	})).Should(Succeed())

	cnrmClusterAccessList := newEmptyClusterAccessObjs(ns, clusterName)

	for _, resources := range cnrmClusterAccessList {
		g.Expect(resources.Objs.EachListItem(func(obj runtime.Object) error {
			return cst.Client.Get(ctx, resources.Key, obj)
		})).Should(Succeed())
	}

	err = cst.Client.Delete(ctx, remoteObj)
	g.Expect(err).ToNot(HaveOccurred())

	err = cst.Client.Get(ctx, key, remoteObj)
	g.Expect(err).To(HaveOccurred())
	g.Expect(apierrors.IsNotFound(err)).To(BeTrue())

	g.Eventually(func() bool {
		deleted := 0
		for i := range cnrmClusterCore.Objs.Items {
			err := cst.Client.Get(ctx, cnrmClusterCore.Key, &cnrmClusterCore.Objs.Items[i])
			if apierrors.IsNotFound(err) {
				deleted++
			}
		}
		for _, resources := range cnrmClusterAccessList {
			for i := range resources.Objs.Items {
				err := cst.Client.Get(ctx, resources.Key, &resources.Objs.Items[i])
				if apierrors.IsNotFound(err) {
					deleted++
				}
			}
		}
		// we only check CNRM resrource, the service account is not
		// part of this and there is no doubt it would get deleted
		// just as well
		return deleted == 7
	}, *pollTimeout, *pollInterval).Should(BeTrue())

	g.Expect(getMetricIntValue(cst.MetricTracker.ClustersCreated)).To(Equal(initialClusterCount + 1))
}

func createDeleteClusterWithStatusUpdates(g *WithT, cst *ControllerSubTest) {
	ctx := context.Background()
	ns := cst.NextNamespace()

	key, obj := newTestClusterGKE(ns, "test-2")
	remoteObj := obj.DeepCopy()

	obj.Spec.ConfigTemplate = new(string)
	*obj.Spec.ConfigTemplate = "basic"

	obj.Spec.JobSpec = &v1alpha1.TestClusterGKEJobSpec{
		Runner: &v1alpha1.TestClusterGKEJobRunnerSpec{
			Image:   new(string),
			Command: []string{"sleep", "10"},
		},
		SkipInit: true,
	}
	*obj.Spec.JobSpec.Runner.Image = "busybox:1.32"

	err := cst.Client.Get(ctx, key, remoteObj)
	g.Expect(err).To(HaveOccurred())
	g.Expect(apierrors.IsNotFound(err)).To(BeTrue())

	g.Expect(cst.Client.Create(ctx, obj)).To(Succeed())

	g.Expect(cst.Client.Get(ctx, key, remoteObj)).To(Succeed())

	listOpts := &client.ListOptions{Namespace: ns}
	g.Eventually(func() bool {
		cnrmObjs := cnrm.NewContainerClusterList()
		err := cst.Client.List(ctx, cnrmObjs, listOpts)
		return err == nil && len(cnrmObjs.Items) > 0
	}, *pollTimeout, *pollInterval).Should(BeTrue())

	cnrmObjs := cnrm.NewContainerClusterList()
	g.Expect(cst.Client.List(ctx, cnrmObjs, listOpts)).To(Succeed())
	g.Expect(cnrmObjs.Items).To(HaveLen(1))

	clusterName := cnrmObjs.Items[0].GetName()
	g.Expect(clusterName).To(HavePrefix("test-2-"))
	g.Expect(clusterName).To(HaveLen(12))

	// TODO: this is not the case for some reason, should check why that maybe...
	// g.Expect(obj.Status.ClusterName).ToNot(BeNil())
	// g.Expect(*obj.Status.ClusterName).To(Equal(clusterName))

	cnrmCluster := &cnrmObjs.Items[0]

	// creation sequence:
	// 1: {"conditions":[{"type":"Ready","status":"False","lastTransitionTime":"2020-07-14T17:49:03Z","reason":"DependencyNotReady","message":"reference ComputeNetwork test-clusters-dev/test-1 is not ready"}],"endpoint":"","operation":""}
	// 2: {"conditions":[{"type":"Ready","status":"False","lastTransitionTime":"2020-07-14T17:49:03Z","reason":"Updating","message":"Update in progress"}],"endpoint":"","operation":""}
	// 3: {"conditions":[{"type":"Ready","status":"True","lastTransitionTime":"2020-07-14T17:55:36Z","reason":"UpToDate","message":"The resource is up to date"}],"endpoint":"35.197.220.90","operation":""}
	t := v1.Time{Time: time.Now().Round(time.Second)}
	createConditionSequence := [][]ContainerClusterStatusCondition{
		{{
			Type:               "Ready",
			Status:             "False",
			Reason:             "DependencyNotReady",
			LastTransitionTime: t,
		}},
		{{
			Type:               "Ready",
			Status:             "False",
			Reason:             "Updating",
			LastTransitionTime: t,
		}},
		{{
			Type:               "Ready",
			Status:             "True",
			Reason:             "UpToDate",
			LastTransitionTime: t,
		}},
	}

	for _, conditions := range createConditionSequence {
		cnrmCluster.Object["status"] = ContainerClusterStatus{
			Conditions: conditions,
		}
		// check status is not the same initially
		g.Expect(cst.Client.Get(ctx, key, obj)).To(Succeed())
		g.Expect(obj.Status.Conditions).NotTo(ConsistOf(cnrmCluster.Object["status"].(v1alpha1.TestClusterGKEStatus).Conditions))
		// make an update simulating what CNRM would do
		g.Expect(cst.Client.Update(ctx, cnrmCluster)).To(Succeed())
		// expect the status to be exactly the same soon enough
		g.Eventually(func() []v1alpha1.TestClusterGKEStatusCondition {
			err := cst.Client.Get(ctx, key, obj)
			if err != nil {
				return nil
			}

			if len(obj.Status.Conditions) == 0 {
				return nil
			}
			return obj.Status.Conditions
		}, *pollTimeout, *pollInterval).Should(ConsistOf(conditions))

		g.Expect(obj.Status.ClusterName).ToNot(BeNil())
		g.Expect(*obj.Status.ClusterName).To(Equal(clusterName))

		if conditions[0].Status == "True" {
			g.Expect(obj.Status.HasReadyCondition()).To(BeTrue())
		} else {
			g.Expect(obj.Status.HasReadyCondition()).To(BeFalse())
		}
	}

	cnrmCheckLeakedObjs := cnrm.NewContainerClusterList()
	g.Expect(cst.Client.List(ctx, cnrmCheckLeakedObjs, listOpts))
	g.Expect(cnrmCheckLeakedObjs.Items).To(HaveLen(1))

	jobObj := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-runner-" + clusterName,
			Namespace: ns,
		},
	}

	jobKey := types.NamespacedName{
		Name:      jobObj.Name,
		Namespace: jobObj.Namespace,
	}

	g.Eventually(func() error {
		return cst.Client.Get(ctx, jobKey, &jobObj)
	}, *pollTimeout, *pollInterval).Should(Succeed())

	g.Eventually(func() error {
		err := cst.Client.Get(ctx, jobKey, &jobObj)
		if err != nil {
			return err
		}
		if IsJobCompleted(jobObj) {
			return nil
		}
		return fmt.Errorf("Test job is not complete yet")

	}, *pollTimeout, *pollInterval).Should(Succeed())

	g.Eventually(func() error {
		err := cst.Client.Get(ctx, key, remoteObj)
		if err == nil {
			return fmt.Errorf("Test cluster not deleted yet")
		}
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}, *pollTimeout, *pollInterval).Should(Succeed())
	return
}
