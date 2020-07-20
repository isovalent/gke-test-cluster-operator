// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package basic

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

_name:        "\(defaults.metadata.name)" | *"\(resource.metadata.name)"
_namespace:   "\(defaults.metadata.namespace)" | *"\(resource.metadata.namespace)"
_location:    "\(defaults.spec.location)" | *"\(resource.spec.location)"
_region:      "\(defaults.spec.region)" | *"\(resource.spec.region)"
_machineType: "\(defaults.spec.machineType)" | *"\(resource.spec.machineType)"

_project: "cilimum-ci"

// TODO (post-mvp): IAM resources are implementation detail of the operator, so should be
// factored into another template or file or package

ClusterTemplate :: {
	kind:       "List"
	apiVersion: "v1"
	items: [{
		apiVersion: "container.cnrm.cloud.google.com/v1beta1"
		kind:       "ContainerCluster"
		metadata: {
			name:      "\(_name)"
			namespace: "\(_namespace)"
			labels: cluster: "\(_name)"
			annotations: {
				"cnrm.cloud.google.com/remove-default-node-pool": "true"
				"cnrm.cloud.google.com/project-id":               "\(_project)"
			}
		}
		spec: {
			initialNodeCount: 1
			location:         "\(_location)"
			loggingService:   "logging.googleapis.com/kubernetes"
			masterAuth: clientCertificateConfig: issueClientCertificate: false
			monitoringService: "monitoring.googleapis.com/kubernetes"
			networkRef: name:    "\(_name)"
			subnetworkRef: name: "\(_name)"
		}
	}, {
		apiVersion: "container.cnrm.cloud.google.com/v1beta1"
		kind:       "ContainerNodePool"
		metadata: {
			name:      "\(_name)"
			namespace: "\(_namespace)"
			labels: cluster: "\(_name)"
			annotations: {
				"cnrm.cloud.google.com/project-id": "\(_project)"
			}
		}
		spec: {
			clusterRef: name: "\(_name)"
			initialNodeCount: 0
			location:         "\(_location)"
			management: {
				autoRepair:  false
				autoUpgrade: false
			}
			nodeConfig: {
				diskSizeGb:  100
				diskType:    "pd-standard"
				machineType: "\(_machineType)"
				metadata: "disable-legacy-endpoints": "true"
				oauthScopes: [
					"https://www.googleapis.com/auth/logging.write",
					"https://www.googleapis.com/auth/monitoring",
				]
			}
		}
	}, {
		apiVersion: "compute.cnrm.cloud.google.com/v1beta1"
		kind:       "ComputeNetwork"
		metadata: {
			name:      "\(_name)"
			namespace: "\(_namespace)"
			labels: cluster: "\(_name)"
			annotations: {
				"cnrm.cloud.google.com/project-id": "\(_project)"
			}
		}
		spec: {
			autoCreateSubnetworks:       false
			deleteDefaultRoutesOnCreate: false
			routingMode:                 "REGIONAL"
		}
	}, {
		apiVersion: "compute.cnrm.cloud.google.com/v1beta1"
		kind:       "ComputeSubnetwork"
		metadata: {
			name:      "\(_name)"
			namespace: "\(_namespace)"
			labels: cluster: "\(_name)"
			annotations: {
				"cnrm.cloud.google.com/project-id": "\(_project)"
			}
		}
		spec: {
			ipCidrRange: "10.128.0.0/20"
			networkRef: name: "\(_name)"
			region: "\(_region)"
		}
	}, {
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
}

defaults: v1alpha1.TestClusterGKE

resource: v1alpha1.TestClusterGKE

template: ClusterTemplate
