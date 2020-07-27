// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package basic

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

_name:        "\(defaults.metadata.name)" | *"\(resource.metadata.name)"
_namespace:   "\(defaults.metadata.namespace)" | *"\(resource.metadata.namespace)"
_location:    "\(defaults.spec.location)" | *"\(resource.spec.location)"
_region:      "\(defaults.spec.region)" | *"\(resource.spec.region)"
_machineType: "\(defaults.spec.machineType)" | *"\(resource.spec.machineType)"

_project: "cilium-ci"

#ClusterTemplate: {
	kind:       "List"
	apiVersion: "v1"
	items:      [{
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
	}] + #IAM
}

defaults: v1alpha1.#TestClusterGKE

resource: v1alpha1.#TestClusterGKE

template: #ClusterTemplate
