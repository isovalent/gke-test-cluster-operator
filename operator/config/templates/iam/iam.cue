// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package iam

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

_name:     "\(resource.metadata.name)"
_namespace: "\(defaults.metadata.namespace)" | *"\(resource.metadata.namespace)"

_project: "cilium-ci"

_serviceAccountName: "\(_name)-admin"
_serviceAccountRef:  "serviceAccount:\(_project).svc.id.goog[\(_namespace)/\(_serviceAccountName)]"

#ClusterAccessResources: {
	kind:       "List"
	apiVersion: "v1"
	items: [
		{
			apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
			kind:       "IAMServiceAccount"
			metadata: {
				name:      "\(_name)"
				namespace: "\(_namespace)"
				labels: cluster: "\(_name)"
				annotations: {
					"cnrm.cloud.google.com/project-id": "\(_project)"
				}
			}
			spec: displayName: "\(_serviceAccountName)"
		},
		{
			apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
			kind:       "IAMPolicy"
			metadata: {
				name:      "\(_name)"
				namespace: "\(_namespace)"
				labels: cluster: "\(_name)"
				annotations: {
					"cnrm.cloud.google.com/project-id": "\(_project)"
				}
			}
			spec: {
				resourceRef: {
					apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
					kind:       "IAMServiceAccount"
					name:       "\(_serviceAccountName)"
				}
				// TODO: is GCP access supposed to be defined here as well?
				bindings: [{
					role: "roles/iam.workloadIdentityUser"
					members: [_serviceAccountRef]
				}]
			}
		},
		// Alternative mode of declaration using IAMPolicyMember
		// {
		//  apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
		//  kind:       "IAMPolicyMember"
		//  metadata: {
		//   name: "\(_name)"
		//   labels: cluster: "\(_name)"
		//   namespace: "\(_namespace)"
		//  }
		//  spec: {
		//   member: "\(_serviceAccountRef)"
		//   role:   "roles/iam.workloadIdentityUser"
		//   resourceRef: {
		//    apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
		//    kind:       "IAMServiceAccount"
		//    name:       "\(_serviceAccountName)"
		//   }
		//  }
		// },
		//   //# TODO: should we still have IAMCustomRole and binding above, or there is a role
		//   //# we can use already?
		//   //# ---
		//   //# apiVersion: iam.cnrm.cloud.google.com/v1beta1
		//   //# kind: IAMCustomRole
		//   //# metadata:
		//   //#   name: containerclusteradmin
		//   //#   namespace: \(_namespace)
		//   //# spec:
		//   //#   title: Admin role for GKE clusters
		//   //#   description: This role only contains permissions to access GKE clusters
		//   //#   permissions:
		//   //#     # TOOD
		//   //#     - container.clusters.get
		//   //#   stage: GA
		{
			apiVersion: "v1"
			kind:       "ServiceAccount"
			metadata: {
				name:      "\(_serviceAccountName)"
				namespace: "\(_namespace)"
				labels: cluster: "\(_name)"
				annotations: {
					"iam.gke.io/gcp-service-account":   "\(_serviceAccountName)@\(_project).iam.gserviceaccount.com"
					"cnrm.cloud.google.com/project-id": "\(_project)"
				}
			}
		},
	]
}

defaults: v1alpha1.#TestClusterGKE

resource: v1alpha1.#TestClusterGKE

template: #ClusterAccessResources
