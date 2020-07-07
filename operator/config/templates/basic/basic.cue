package basic

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

// TODO: https://github.com/cuelang/cue/discussions/439

ClusterTemplate :: {
	kind: "List"
	apiVersion: "v1"
	items: [{
		apiVersion: "container.cnrm.cloud.google.com/v1beta1"
		kind:       "ContainerCluster"
		metadata: {
			annotations: "cnrm.cloud.google.com/remove-default-node-pool": "true"
			labels: cluster: "\(resource.metadata.name)" | ""
			name:      "\(resource.metadata.name)" | ""
			namespace: "\(resource.metadata.namespace)" | ""
		}
		spec: {
			initialNodeCount: 1
			location:         "\(resource.spec.location)" | ""
			loggingService:   "logging.googleapis.com/kubernetes"
			masterAuth: clientCertificateConfig: issueClientCertificate: false
			monitoringService: "monitoring.googleapis.com/kubernetes"
			networkRef: name: "\(resource.metadata.name)" | ""
			subnetworkRef: name: "\(resource.metadata.name)" | ""
		}
	}, {
		apiVersion: "container.cnrm.cloud.google.com/v1beta1"
		kind:       "ContainerNodePool"
		metadata: {
			labels: cluster: "\(resource.metadata.name)" | ""
			name:      "\(resource.metadata.name)" | ""
			namespace: "\(resource.metadata.namespace)" | ""
		}
		spec: {
			clusterRef: name: "\(resource.metadata.name)" | ""
			initialNodeCount: 0
			location:         "\(resource.spec.location)" | ""
			management: {
				autoRepair:  false
				autoUpgrade: false
			}
			nodeConfig: {
				diskSizeGb:  100
				diskType:    "pd-standard"
				machineType: "\(resource.spec.machineType)" | ""
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
			labels: cluster: "\(resource.metadata.name)" | ""
			name:      "\(resource.metadata.name)" | ""
			namespace: "\(resource.metadata.namespace)" | ""
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
			labels: cluster: "\(resource.metadata.name)" | ""
			name:      "\(resource.metadata.name)" | ""
			namespace: "\(resource.metadata.namespace)"  | ""
		}
		spec: {
			ipCidrRange: "10.128.0.0/20"
			networkRef: name: "\(resource.metadata.name)" | ""
			region: "us-west1" // TODO: split region from location
		}
	}]
}

resource: v1alpha1.TestClusterGKE2

// defaults: v1alpha1.TestClusterGKE
// defaults: {
// 	metadata: name: ""
// }

template: ClusterTemplate
