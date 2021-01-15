package requesteraccess

import "strings"

_name: "\(parameters.namespace)-ci"

_commonLabels: name: _name

_commonAnnotations: "cnrm.cloud.google.com/project-id": parameters.project

_roleName: strings.Replace(_name, "-", "", -1)
_serviceAccountName: "\(_name)@\(parameters.project).iam.gserviceaccount.com"

_clusterRoleName: "gke-test-cluster-operator-ci"

#RequesterAccessResources: {
	apiVersion: "v1"
	kind:       "List"
	items: [{
		apiVersion: "rbac.authorization.k8s.io/v1beta1"
		kind:       "RoleBinding"
		metadata: {
			labels:    _commonLabels
			name:      _name
			namespace: parameters.namespace
		}
		roleRef: {
			apiGroup: "rbac.authorization.k8s.io"
			kind:     "ClusterRole"
			name:     _clusterRoleName
		}
		subjects: [{
			kind: "User"
			name: _serviceAccountName
		}]
	}, {
		apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
		kind:       "IAMServiceAccount"
		metadata: {
			annotations: _commonAnnotations
			labels:      _commonLabels
			name:        _name
			namespace:   parameters.namespace
		}
	}, {
		apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
		kind:       "IAMCustomRole"
		metadata: {
			annotations: _commonAnnotations
			labels:      _commonLabels
			name:        _roleName
			namespace:   parameters.namespace
		}
		spec: {
			title:       "Role for CI access to clusters in namespace \(parameters.namespace)"
			description: "This role only contains permissions required for CI access"
			permissions: [
				"container.thirdPartyObjects.create",
				"container.thirdPartyObjects.get",
				"container.pods.list",
				"container.pods.get",
				"container.pods.getLogs",
				"container.jobs.list",
				"container.jobs.get",
				"container.clusters.list",
				"container.clusters.get",
				"container.clusters.getCredentials",
				"container.configMaps.create",
			]
			stage: "GA"
		}
	}, {
		apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
		kind:       "IAMPolicyMember"
		metadata: {
			annotations: _commonAnnotations
			labels:      _commonLabels
			name:        _name
			namespace:   parameters.namespace
		}
		spec: {
			member: "serviceAccount:\(_serviceAccountName)"
			resourceRef: {
				apiVersion: "resourcemanager.cnrm.cloud.google.com/v1beta1"
				external:   "projects/\(parameters.project)"
				kind:       "Project"
			}
			role: "projects/\(parameters.project)/roles/\(_roleName)"
		}
	}]
}

#RequesterAccessParameters: {
	project: string
	namespace: string
}

parameters: #RequesterAccessParameters
template:   #RequesterAccessResources
