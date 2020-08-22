// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package basic

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha2"

_name:        "\(resource.metadata.name)"
_namespace:   "\(defaults.metadata.namespace)" | *"\(resource.metadata.namespace)"
_project:     "\(defaults.spec.project)" | *"\(resource.spec.project)"
_location:    "\(defaults.spec.location)" | *"\(resource.spec.location)"
_region:      "\(defaults.spec.region)" | *"\(resource.spec.region)"
_nodes:       defaults.spec.nodes | *resource.spec.nodes
_machineType: "\(defaults.spec.machineType)" | *"\(resource.spec.machineType)"

#ClusterCoreResources: {
	kind:       "List"
	apiVersion: "v1"
	items: [
		{
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
		},
		{
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
				initialNodeCount: _nodes
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
		},
		{
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
		},
		{
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
		},
	]
}

defaults: v1alpha2.#TestClusterGKE

resource: v1alpha2.#TestClusterGKE

template: #ClusterCoreResources
