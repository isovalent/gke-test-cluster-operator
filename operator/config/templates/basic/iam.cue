// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package basic

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

_name:        "\(defaults.metadata.name)" | *"\(resource.metadata.name)"
_namespace:   "\(defaults.metadata.namespace)" | *"\(resource.metadata.namespace)"

_project: "cilium-ci"

IAM :: [{
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
		spec: displayName: "\(_name)-admin"
	}, {
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
				name:       "\(_name)-admin"
			}
			// TODO: is GCP access supposed to be defined here as well?
			bindings: [{
				role: "roles/iam.workloadIdentityUser"
				members: [
					"serviceAccount:\(_project).svc.id.goog[\(_namespace)/\(_name)-admin]",
				]
			}]
			//# TODO: should we still have IAMCustomRole and binding above, or there is a role
			//# we can use already?
			//# ---
			//# apiVersion: iam.cnrm.cloud.google.com/v1beta1
			//# kind: IAMCustomRole
			//# metadata:
			//#   name: containerclusteradmin
			//#   namespace: \(_namespace)
			//# spec:
			//#   title: Admin role for GKE clusters
			//#   description: This role only contains permissions to access GKE clusters
			//#   permissions:
			//#     # TOOD
			//#     - container.clusters.get
			//#   stage: GA
			//# ---
			//# apiVersion: iam.cnrm.cloud.google.com/v1beta1
			//# kind: IAMPolicyMember
			//# metadata:
			//#   name: \(_name)-admin-gcpsa
			//#   labels:
			//#     cluster: \(_name)
			//#   namespace: \(_namespace)
			//# spec:
			//#   member: serviceAccount:\(_name)-admin@\(_project).iam.gserviceaccount.com
			//#   role: projects/\(_project)/roles/containerclusteradmin
			//#   resourceRef:
			//#     apiVersion: container.cnrm.cloud.google.com/v1beta1
			//#     kind: ContainerCluster
			//#     name: \(_name)
		}
	}, {
		apiVersion: "v1"
		kind:       "ServiceAccount"
		metadata: {
			name:      "\(_name)-admin"
			namespace: "\(_namespace)"
			labels: cluster: "\(_name)"
			annotations: {
				"iam.gke.io/gcp-service-account":   "\(_name)-admin@\(_project).iam.gserviceaccount.com"
				"cnrm.cloud.google.com/project-id": "\(_project)"
			}
		}
	}]

defaults: v1alpha1.TestClusterGKE

resource: v1alpha1.TestClusterGKE
