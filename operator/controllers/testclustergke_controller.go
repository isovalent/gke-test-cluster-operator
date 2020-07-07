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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
)

// TestClusterGKEReconciler reconciles a TestClusterGKE object
type TestClusterGKEReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=clusters.ci.cilium.io,resources=testclustergkes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=clusters.ci.cilium.io,resources=testclustergkes/status,verbs=get;update;patch
func (r *TestClusterGKEReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("testclustergke", req.NamespacedName)

	var testClusterGKE v1alpha1.TestClusterGKE
	if err := r.Get(ctx, req.NamespacedName, &testClusterGKE); err != nil {
		log.Error(err, "unable to fetch testclustergke")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

func (r *TestClusterGKEReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clustersv1alpha1.TestClusterGKE{}).
		Complete(r)
}
