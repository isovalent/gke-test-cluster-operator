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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	clustersv1alpha1 "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

	"github.com/isovalent/gke-test-cluster-management/operator/pkg/config"
)

// TestClusterGKEReconciler reconciles a TestClusterGKE object
type TestClusterGKEReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	ConfigRenderer *config.Config
}

// +kubebuilder:rbac:groups=clusters.ci.cilium.io,resources=testclustergkes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=clusters.ci.cilium.io,resources=testclustergkes/status,verbs=get;update;patch
func (r *TestClusterGKEReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("Reconcile", req.NamespacedName)

	var testClusterGKE v1alpha1.TestClusterGKE
	if err := r.Get(ctx, req.NamespacedName, &testClusterGKE); err != nil {
		log.Error(err, "unable to fetch object")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	objs, err := r.ConfigRenderer.RenderObjects(&testClusterGKE)
	if err != nil {
		log.Error(err, "unable render config template")
		return ctrl.Result{}, err
	}
	log.Info("generated config", "items", objs.Items)

	// TODO (mvp)
	// - detect event type, error on updates
	// - handle deletion
	// - write a few simple controller tests
	// - update RBAC configs
	// - de-kustomize configs
	// TODO (post-mvp)
	// - implement pool object
	// - implement GCP project annotation
	// - implement job runner pod (use sonoboy as PoC)
	// - use random cluster name, instead of same as test object
	// - wait for cluster to get created, update status

	if err := objs.EachListItem(r.createOrSkip); err != nil {
		log.Error(err, "unable reconcile object")
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

func (r *TestClusterGKEReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clustersv1alpha1.TestClusterGKE{}).
		Complete(r)
}

func (r *TestClusterGKEReconciler) createOrSkip(obj runtime.Object) error {
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		return err
	}

	ctx := context.Background()
	log := r.Log.WithValues("createOrSkip", key)

	// TODO (post-mvp) probably don't need to make a full copy,
	// should be able to copy just TypeMeta and ObjectMeta
	remoteObj := obj.DeepCopyObject()
	getErr := r.Client.Get(ctx, key, remoteObj)
	if apierrors.IsNotFound(getErr) {
		log.Info("will create", "obj", obj)
		return r.Client.Create(ctx, obj)
	}
	if getErr == nil {
		log.Info("already exists", "remoteObj", remoteObj)
	}
	return getErr
}
