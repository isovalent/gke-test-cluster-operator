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

	. "github.com/onsi/gomega"

	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

func TestControllers(t *testing.T) {
	g := NewGomegaWithT(t)

	client, env := setup(t)
	ctx := context.Background()

	key := types.NamespacedName{
		Namespace: "default",
		Name:      "test-1",
	}

	obj := &clustersv1alpha1.TestClusterGKE{}
	err := client.Get(ctx, key, obj)
	g.Expect(err).ToNot(HaveOccurred())

	teardown(t, env)
}
