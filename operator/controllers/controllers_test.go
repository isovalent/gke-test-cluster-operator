// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers_test

import (
	"context"
	"testing"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

	listOpts := &client.ListOptions{Namespace: ns}
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

	cnrmKey, cnrmObjs := newEmptyContainerClusterObjs(ns, clusterName)

	g.Expect(cnrmObjs.EachListItem(func(obj runtime.Object) error {
		return cst.Client.Get(ctx, cnrmKey, obj)
	})).Should(Succeed())

	err = cst.Client.Delete(ctx, remoteObj)
	g.Expect(err).ToNot(HaveOccurred())

	err = cst.Client.Get(ctx, key, remoteObj)
	g.Expect(err).To(HaveOccurred())
	g.Expect(apierrors.IsNotFound(err)).To(BeTrue())

	g.Eventually(func() bool {
		deleted := 0
		for i := range cnrmObjs.Items {
			err := cst.Client.Get(ctx, cnrmKey, &cnrmObjs.Items[i])
			if apierrors.IsNotFound(err) {
				deleted++
			}
		}
		return deleted == 4
	}, *pollTimeout, *pollInterval).Should(BeTrue())
}

func createDeleteClusterWithStatusUpdates(g *WithT, cst *ControllerSubTest) {
	ctx := context.Background()
	ns := cst.NextNamespace()

	key, obj := newTestClusterGKE(ns, "test-2")
	remoteObj := obj.DeepCopy()

	obj.Spec.ConfigTemplate = new(string)
	*obj.Spec.ConfigTemplate = "basic"

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
		g.Expect(obj.Status).NotTo(Equal(cnrmCluster.Object["status"]))
		// make an update simulating what CNRM would do
		g.Expect(cst.Client.Update(ctx, cnrmCluster)).To(Succeed())
		// expect the status to be exactly the same soon enough
		g.Eventually(func() []v1alpha1.TestClusterGKEStatusCondition {
			_ = cst.Client.Get(ctx, key, obj)
			return obj.Status.Conditions
		}, *pollTimeout, *pollInterval).Should(Equal(conditions))
	}

	g.Expect(cst.Client.Delete(ctx, remoteObj)).To(Succeed())

	// since there are no real expectations around state transitions on deletion,
	// there is no need to simulate what happens there

	err = cst.Client.Get(ctx, key, remoteObj)
	g.Expect(err).To(HaveOccurred())
	g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
}
