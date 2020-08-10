// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package iam

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

_name:      "\(resource.metadata.name)"
_namespace: "\(defaults.metadata.namespace)" | *"\(resource.metadata.namespace)"

_project: "\(defaults.spec.project)" | *"\(resource.spec.project)"

_adminServiceAccountName:  "\(_name)-admin"
_adminServiceAccountEmail: "\(_adminServiceAccountName)@\(_project).iam.gserviceaccount.com"
_adminServiceAccountRef:   "serviceAccount:\(_project).svc.id.goog[\(_namespace)/\(_adminServiceAccountName)]"

#ClusterAccessResources: {
	kind:       "List"
	apiVersion: "v1"
	items: [
		{
			apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
			kind:       "IAMServiceAccount"
			metadata: {
				name:      "\(_adminServiceAccountName)"
				namespace: "\(_namespace)"
				labels: cluster: "\(_name)"
				annotations: {
					"cnrm.cloud.google.com/project-id": "\(_project)"
				}
			}
		},
		{
			apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
			kind:       "IAMPolicyMember"
			metadata: {
				name: "\(_name)-workload-identity"
				labels: cluster: "\(_name)"
				namespace: "\(_namespace)"
				annotations: {
					"cnrm.cloud.google.com/project-id": "\(_project)"
				}
			}
			spec: {
				member: "\(_adminServiceAccountRef)"
				role:   "roles/iam.workloadIdentityUser"
				resourceRef: {
					apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
					kind:       "IAMServiceAccount"
					name:       "\(_adminServiceAccountName)"
				}
			}
		},
		{
			apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
			kind:       "IAMPolicyMember"
			metadata: {
				name: "\(_name)-cluster-admin"
				labels: cluster: "\(_name)"
				namespace: "\(_namespace)"
				annotations: {
					"cnrm.cloud.google.com/project-id": "\(_project)"
				}
			}
			spec: {
				member: "serviceAccount:\(_adminServiceAccountEmail)"
				role:   "roles/container.clusterAdmin"
				resourceRef: {
					// At the moment ContainerCluster cannot be referenced here, so it's at project level for now
					// (see https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/248)
					// Note that the project-level clusterAdmin role doesn't get automatic access
					// to Kubernetes API in all clusters, unless there is a CRB inside Kubernetes
					apiVersion: "resourcemanager.cnrm.cloud.google.com/v1beta1"
					kind:       "Project"
					external:   "projects/\(_project)"
				}
			}
		},
		{
			apiVersion: "v1"
			kind:       "ServiceAccount"
			metadata: {
				name:      "\(_adminServiceAccountName)"
				namespace: "\(_namespace)"
				labels: cluster: "\(_name)"
				annotations: {
					"iam.gke.io/gcp-service-account":   "\(_adminServiceAccountEmail)"
					"cnrm.cloud.google.com/project-id": "\(_project)"
				}
			}
		},
	]
}

defaults: v1alpha1.#TestClusterGKE

resource: v1alpha1.#TestClusterGKE

template: #ClusterAccessResources
