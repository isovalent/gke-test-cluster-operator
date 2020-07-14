/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers_test

import (
	"context"
	"testing"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	. "github.com/onsi/gomega"
)

func TestControllers(t *testing.T) {
	cstm, teardown := setup(t)

	cstm.NewControllerSubTest(t).
		Run("simple test - create and delete objects", simpleCreateDeleteObjects)

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

	cnrmKey, cnrmObjs := newContainerClusterObjs(g, ns, "test-1")

	err = cnrmObjs.EachListItem(func(obj runtime.Object) error {
		return cst.Client.Get(ctx, cnrmKey, obj)
	})
	g.Expect(err).To(HaveOccurred())
	g.Expect(apierrors.IsNotFound(err)).To(BeTrue())

	g.Eventually(func() error {
		return cnrmObjs.EachListItem(func(obj runtime.Object) error {
			return cst.Client.Get(ctx, cnrmKey, obj)
		})
	}, *pollTimeout, *pollInterval).Should(Succeed())

	err = cst.Client.Delete(ctx, remoteObj)
	g.Expect(err).ToNot(HaveOccurred())

	err = cst.Client.Get(ctx, key, remoteObj)
	g.Expect(err).To(HaveOccurred())
	g.Expect(apierrors.IsNotFound(err)).To(BeTrue())

	g.Eventually(func() error {
		return cnrmObjs.EachListItem(func(obj runtime.Object) error {
			return cst.Client.Get(ctx, cnrmKey, obj)
		})
	}, *pollTimeout, *pollInterval).ShouldNot(Succeed())
}
