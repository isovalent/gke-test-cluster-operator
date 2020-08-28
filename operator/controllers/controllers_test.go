// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/gomega"

	"github.com/isovalent/gke-test-cluster-management/operator/api/cnrm"
	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha2"
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

	g.Eventually(func() bool {
		systemConfigMapObj := &corev1.ConfigMap{}
		systemConfigMapKey := types.NamespacedName{
			Name:      clusterName + "-system",
			Namespace: ns,
		}
		err := cst.Client.Get(ctx, systemConfigMapKey, systemConfigMapObj)
		return err == nil
	}, *pollTimeout, *pollInterval).Should(BeTrue())

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

	for _, tc := range []struct {
		command    []string
		shouldFail bool
		name       string
	}{
		// the job needs to run long enought for cluster to not
		// get deleted too quickly
		{
			command:    []string{"sleep", "25"},
			shouldFail: false,
			name:       "test-2-good-job",
		},
		{
			command:    []string{"sh", "-c", "sleep 25 && exit 1"},
			shouldFail: true,
			name:       "test-2-failing-job",
		},
	} {

		ns := cst.NextNamespace()

		key, obj := newTestClusterGKE(ns, tc.name)
		remoteObj := obj.DeepCopy()

		obj.Spec.ConfigTemplate = new(string)
		*obj.Spec.ConfigTemplate = "basic"

		obj.Spec.JobSpec = &v1alpha2.TestClusterGKEJobSpec{
			Runner: &v1alpha2.TestClusterGKEJobRunnerSpec{
				InitImage: new(string),
				Image:     new(string),
				Command:   tc.command,
			},
		}

		*obj.Spec.JobSpec.Runner.InitImage = "tianon/true@sha256:009cce421096698832595ce039aa13fa44327d96beedb84282a69d3dbcf5a81b"
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
		g.Expect(clusterName).To(HavePrefix(tc.name + "-"))
		g.Expect(clusterName).To(HaveLen(len(tc.name) + 6))

		g.Eventually(func() bool {
			if cst.Client.Get(ctx, key, remoteObj) != nil {
				return false
			}

			return (remoteObj.Status.ClusterName != nil) &&
				(*remoteObj.Status.ClusterName == clusterName)
		}, *pollTimeout, *pollInterval).Should(BeTrue())

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

		cnrmClusterDependecyKey := fmt.Sprintf("ContainerCluster:%s/%s", cnrmCluster.GetNamespace(), clusterName)

		for _, conditions := range createConditionSequence {
			cnrmCluster.Object["status"] = ContainerClusterStatus{
				Conditions: conditions,
			}
			// check status is not the same initially
			g.Expect(cst.Client.Get(ctx, key, obj)).To(Succeed())
			g.Expect(obj.Status.Conditions).NotTo(ConsistOf(cnrmCluster.Object["status"].(v1alpha2.TestClusterGKEStatus).Conditions))
			// make an update simulating what CNRM would do
			// NB: CNRM resources don't have status subresource
			g.Expect(cst.Client.Update(ctx, cnrmCluster)).To(Succeed())
			// expect the depenencies status to be exactly the same soon enough
			g.Eventually(func() []v1alpha2.TestClusterGKECondition {
				if cst.Client.Get(ctx, key, obj) != nil {
					return nil
				}

				if len(obj.Status.Dependencies) == 0 {
					return nil
				}

				if _, ok := obj.Status.Dependencies[cnrmClusterDependecyKey]; !ok {
					return nil
				}

				return obj.Status.Dependencies[cnrmClusterDependecyKey]
			}, *pollTimeout, *pollInterval).Should(ConsistOf(conditions))

			g.Expect(obj.Status.ClusterName).ToNot(BeNil())
			g.Expect(*obj.Status.ClusterName).To(Equal(clusterName))

			if conditions[0].Status == "True" {
				g.Expect(obj.Status.AllDependeciesReady()).To(BeTrue())
				g.Expect(obj.Status.HasReadyCondition()).To(BeTrue())

			} else {
				g.Expect(obj.Status.AllDependeciesReady()).To(BeFalse())
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
			if IsJobDone(jobObj) {
				return nil
			}
			return fmt.Errorf("test job is not done yet")
		}, *pollTimeout, *pollInterval).Should(Succeed())

		if tc.shouldFail {
			g.Expect(jobObj.Status.CompletionTime).To(BeNil())
			g.Expect(string(jobObj.Status.Conditions[0].Type)).To(Equal("Failed"))
			g.Expect(string(jobObj.Status.Conditions[0].Status)).To(Equal("True"))
			g.Expect(string(jobObj.Status.Conditions[0].Reason)).To(Equal("BackoffLimitExceeded"))
		} else {
			g.Expect(jobObj.Status.CompletionTime).ToNot(BeNil())
		}

		g.Eventually(func() error {
			err := cst.Client.Get(ctx, key, remoteObj)
			if err == nil {
				return fmt.Errorf("test cluster not deleted yet")
			}
			if apierrors.IsNotFound(err) {
				return nil
			}
			return err
		}, *pollTimeout, *pollInterval).Should(Succeed())

	}
}
