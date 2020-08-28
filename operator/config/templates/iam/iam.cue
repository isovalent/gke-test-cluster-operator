// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package iam

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha2"

_generatedName: resource.metadata.name | *resource.status.clusterName

_namespace: defaults.metadata.namespace | *resource.metadata.namespace

_project: defaults.spec.project | *resource.spec.project

_commonLabels: cluster: resource.metadata.name

_commonAnnotations: "cnrm.cloud.google.com/project-id": _project

_adminServiceAccountName:  "\(_generatedName)-admin"
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
				name:        _adminServiceAccountName
				labels:      _commonLabels
				namespace:   _namespace
				annotations: _commonAnnotations
			}
		},
		{
			apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
			kind:       "IAMPolicyMember"
			metadata: {
				name:        "\(_generatedName)-workload-identity"
				labels:      _commonLabels
				namespace:   _namespace
				annotations: _commonAnnotations
			}
			spec: {
				member: _adminServiceAccountRef
				role:   "roles/iam.workloadIdentityUser"
				resourceRef: {
					apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
					kind:       "IAMServiceAccount"
					name:       _adminServiceAccountName
				}
			}
		},
		{
			apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
			kind:       "IAMPolicyMember"
			metadata: {
				name:        "\(_generatedName)-cluster-admin"
				labels:      _commonLabels
				namespace:   _namespace
				annotations: _commonAnnotations
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
				name:      _adminServiceAccountName
				namespace: _namespace
				labels:    _commonLabels
				annotations: {
					"iam.gke.io/gcp-service-account":   _adminServiceAccountEmail
					"cnrm.cloud.google.com/project-id": _project
				}
			}
		},
	]
}

defaults: v1alpha2.#TestClusterGKE

resource: v1alpha2.#TestClusterGKE

template: #ClusterAccessResources
