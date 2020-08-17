import (
	"encoding/hex"
	"crypto/sha256"
)

constants: {
	name:               "gke-test-cluster-logview"
	ingressRoutePrefix: "/\(hex.Encode(sha256.Sum256(constants.name+parameters.ingressRoutePrefixSalt)))"
}

_workload: {
	apiVersion: "apps/v1"
	kind:       "Deployment"
	metadata: {
		name: "\(constants.name)"
		labels: name: "\(constants.name)"
		namespace: "\(parameters.namespace)"
		annotations: "fluxcd.io/automated": "true"
	}
	spec: _workloadSpec
}

_workloadSpec: {
	selector: matchLabels: name: "\(constants.name)"
	template: metadata: labels: name: "\(constants.name)"
	replicas: 2
	template: {
		metadata: labels: name: "\(constants.name)"
		spec: {
			serviceAccount: "\(constants.name)"
			containers: [{
				name: "logview"
				command: ["/usr/bin/gke-test-cluster-logview"]
				image: "\(parameters.image)"
				env: [{
					name: "NAMESPACE"
					valueFrom: fieldRef: fieldPath: "metadata.namespace"
				}]
				resources: {
					limits: {
						cpu:    "100m"
						memory: "100Mi"
					}
					requests: {
						cpu:    "100m"
						memory: "100Mi"
					}
				}
				ports: [{
					name:          "http"
					containerPort: 8080
				}]
			}]
			terminationGracePeriodSeconds: 10
		}
	}
}

#WorkloadTemplate: {
	kind:       "List"
	apiVersion: "v1"
	items: [{
		apiVersion: "v1"
		kind:       "ServiceAccount"
		metadata: {
			name: "\(constants.name)"
			labels: name: "\(constants.name)"
			namespace: "\(parameters.namespace)"
		}
	}, {
		apiVersion: "rbac.authorization.k8s.io/v1beta1"
		kind:       "ClusterRoleBinding"
		metadata: {
			name: "\(parameters.namespace)-\(constants.name)"
			labels: name: "\(constants.name)"
		}
		roleRef: {
			kind:     "ClusterRole"
			name:     "\(constants.name)"
			apiGroup: "rbac.authorization.k8s.io"
		}
		subjects: [{
			kind:      "ServiceAccount"
			name:      "\(constants.name)"
			namespace: "\(parameters.namespace)"
		}]
	},
		{
			apiVersion: "rbac.authorization.k8s.io/v1"
			kind:       "Role"
			metadata: {
				name:      "\(constants.name)"
				namespace: "\(parameters.namespace)"
				labels: name: "\(constants.name)"
			}
			rules: [{
				apiGroups: [""]
				resources: ["pods/log"]
				verbs: ["get"]
			}]
		},
		{
			apiVersion: "rbac.authorization.k8s.io/v1beta1"
			kind:       "RoleBinding"
			metadata: {
				name:      "\(constants.name)"
				namespace: "\(parameters.namespace)"
				labels: name: "\(constants.name)"
			}
			roleRef: {
				kind:     "Role"
				name:     "\(constants.name)"
				apiGroup: "rbac.authorization.k8s.io"
			}
			subjects: [{
				kind:      "ServiceAccount"
				name:      "\(constants.name)"
				namespace: "\(parameters.namespace)"
			}]
		},
		_workload,
		{
			apiVersion: "v1"
			kind:       "Service"
			metadata: {
				name: "\(constants.name)"
				labels: name: "\(constants.name)"
				namespace: "\(parameters.namespace)"
			}
			spec: {
				selector: name: "\(constants.name)"
				ports: [{
					name:       "http"
					port:       80
					targetPort: 8080
				}]
			}
		},
		{
			apiVersion: "projectcontour.io/v1"
			kind:       "HTTPProxy"
			metadata: {
				name: "\(constants.name)"
				labels: name: "\(constants.name)"
				namespace: "\(parameters.namespace)"
			}
			spec: {
				routes: [{
					conditions: [{
						prefix: constants.ingressRoutePrefix
					}]
					services: [{
						name: "\(constants.name)"
						port: 80
					}]
					pathRewritePolicy: {
						replacePrefix: [{
							prefix:      constants.ingressRoutePrefix
							replacement: "/"
						}]
					}
				}]
			}
		},
		{
			apiVersion: "v1"
			kind:       "ConfigMap"
			metadata: {
				name: "\(constants.name)"
				labels: name: "\(constants.name)"
				namespace: "\(parameters.namespace)"
			}
			stringData: {
				ingressRoutePrefix: constants.ingressRoutePrefix
			}
		},
	]
}

#WorkloadParameters: {
	namespace:              string
	image:                  string
	ingressRoutePrefixSalt: string
}

parameters: #WorkloadParameters
template:   #WorkloadTemplate
