constants: {
	name:    "gke-test-cluster-operator"
	project: "cilium-ci"
}

_workload: {
	metadata: {
		name: "\(constants.name)"
		labels: name: "\(constants.name)"
		namespace: "\(parameters.namespace)"
		annotations: "fluxcd.io/automated": "true"
	}
	spec: _workloadSpec
}

_workloadSpec: {
	template: {
		metadata: labels: name: "\(constants.name)"
		spec: {
			serviceAccount: "\(constants.name)"
			volumes: [{
				name: "tmp"
				emptyDir: {}
			}]
			containers: [{
				command: _command
				image:   "\(parameters.image)"
				name:    "operator"
				volumeMounts: [{
					name:      "tmp"
					mountPath: "/tmp"
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
			}]
			terminationGracePeriodSeconds: 10
		}
	}
}

_command: [...string]

// you cannot easily append to a list in CUE, so the extra role is declared
// without any rules
_extra_rbac_clusterrole: {
	apiVersion: "rbac.authorization.k8s.io/v1"
	kind:       "ClusterRole"
	metadata:
		name: "\(parameters.namespace)-\(constants.name)-extra"
}

_extra_rbac_clusterrolebiding: {
	apiVersion: "rbac.authorization.k8s.io/v1beta1"
	kind:       "ClusterRoleBinding"
	metadata: {
		name: "\(parameters.namespace)-\(constants.name)-extra"
		labels: name: "\(constants.name)"
	}
	roleRef: {
		kind:     "ClusterRole"
		name:     "\(parameters.namespace)-\(constants.name)-extra"
		apiGroup: "rbac.authorization.k8s.io"
	}
	subjects: [{
		kind:      "ServiceAccount"
		name:      "\(constants.name)"
		namespace: "\(parameters.namespace)"
	}]
}

if !parameters.test {
	_workload: {
		apiVersion: "apps/v1"
		kind:       "Deployment"
	}
	_workloadSpec: {
		replicas: 1
		selector: matchLabels: name: "\(constants.name)"
		template: {
			metadata: labels: name: "\(constants.name)"
			spec: {
				containers: [{
					env: [{
						name: "GITHUB_TOKEN"
						valueFrom:
							secretKeyRef: {
								optional: true
								name:     "\(constants.name)-github-token"
								key:      "token"
							}
					}]
				}]
			}
		}
	}
	_command: [
		"/usr/bin/\(constants.name)",
		"--enable-leader-election",
	]
}

if parameters.test {
	_workload: {
		apiVersion: "batch/v1"
		kind:       "Job"
	}

	_workloadSpec: {
		backoffLimit: 0
		template: spec: {
			restartPolicy: "Never"
			containers: [{
				name: "operator"
				env: [{
					name: "NAMESPACE"
					valueFrom: fieldRef: fieldPath: "metadata.namespace"
				}]
			}]
		}
	}
	_command: [
		"test.gke-test-cluster-operator-controllers",
		"-test.v",
		"-test.timeout=5m",
		"-resource-prefix=\(parameters.namespace)",
	]

	// in principle, there should be no difference in RBAC configuration
	// between test and regular deployments, but it's currently inevitable
	// to allow creation/deletion of namespaces while testing the operator
	_extra_rbac_clusterrole: {
		rules: [{
			apiGroups: [""]
			resources: ["namespaces"]
			verbs: ["create", "delete"]
		}]
	}
}

_adminServiceAccountEmail: "\(constants.name)@\(constants.project).iam.gserviceaccount.com"
_adminServiceAccountRef:   "serviceAccount:\(constants.project).svc.id.goog[\(parameters.namespace)/\(constants.name)]"

_clusterAdminAccess: [
	{
		apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
		kind:       "IAMServiceAccount"
		metadata: {
			name:      "\(constants.name)"
			namespace: "\(parameters.namespace)"
			labels: name: "\(constants.name)"
			annotations: {
				"cnrm.cloud.google.com/project-id": "\(constants.project)"
			}
		}
	},
	{
		apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
		kind:       "IAMPolicyMember"
		metadata: {
			name: "\(constants.name)-workload-identity"
			labels: name: "\(constants.name)"
			namespace: "\(parameters.namespace)"
			annotations: {
				"cnrm.cloud.google.com/project-id": "\(constants.project)"
			}
		}
		spec: {
			member: "\(_adminServiceAccountRef)"
			role:   "roles/iam.workloadIdentityUser"
			resourceRef: {
				apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
				kind:       "IAMServiceAccount"
				name:       "\(constants.name)"
			}
		}
	},
	{
		apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
		kind:       "IAMPolicyMember"
		metadata: {
			name: "\(constants.name)-cluster-admin"
			labels: name: "\(constants.name)"
			namespace: "\(parameters.namespace)"
			annotations: {
				"cnrm.cloud.google.com/project-id": "\(constants.project)"
			}
		}
		spec: {
			member: "serviceAccount:\(_adminServiceAccountEmail)"
			role:   "roles/container.admin" // "roles/owner"
			resourceRef: {
				// At the moment ContainerCluster cannot be referenced here, so it's at project level for now
				// (see https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/248)
				apiVersion: "resourcemanager.cnrm.cloud.google.com/v1beta1"
				kind:       "Project"
				external:   "projects/cilium-ci"
			}
		}

	},
]

#WorkloadTemplate: {
	kind:       "List"
	apiVersion: "v1"
	items:      [{
		apiVersion: "v1"
		kind:       "ServiceAccount"
		metadata: {
			name: "\(constants.name)"
			labels: name: "\(constants.name)"
			namespace: "\(parameters.namespace)"
			annotations: {
				"iam.gke.io/gcp-service-account": "\(_adminServiceAccountEmail)"
			}
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
		_workload,
		_extra_rbac_clusterrole,
		_extra_rbac_clusterrolebiding,
	] + _clusterAdminAccess
}

#WorkloadParameters: {
	namespace: string
	image:     string
	test:      bool
}

parameters: #WorkloadParameters
template:   #WorkloadTemplate
