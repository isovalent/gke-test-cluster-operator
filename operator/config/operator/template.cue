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

_extra_rbac_ClusterRoleAndBinding: [...{}]

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
_optionalLogviewDomainFlag: [...string]

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
	] + _optionalLogviewDomainFlag
}

if parameters.logviewDomain != null && len(parameters.logviewDomain) > 0{
	_optionalLogviewDomainFlag: [
		"--logview-domain=\(parameters.logviewDomain)",
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
	_extra_rbac_ClusterRoleAndBinding: [
		{
			apiVersion: "rbac.authorization.k8s.io/v1"
			kind:       "ClusterRole"
			metadata:
				name: "\(parameters.namespace)-\(constants.name)-extra"
			rules: [{
				apiGroups: [""]
				resources: ["namespaces"]
				verbs: ["create", "delete"]
			}]
		}, {
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
		},
	]
}

_adminServiceAccountEmail: "\(constants.name)@\(constants.project).iam.gserviceaccount.com"
_adminServiceAccountRef:   "serviceAccount:\(constants.project).svc.id.goog[\(parameters.namespace)/\(constants.name)]"

_serviceAccount: {
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
}

_rbac_ClusterRoleBinding: {
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
}

_iam_clusterAdminAccess: [
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

_extra_certManager_IssuerAndCertificater: [...{}]

if parameters.certManager {
	_extra_certManager_IssuerAndCertificater: [{
		apiVersion: "cert-manager.io/v1beta1"
		kind:       "Issuer"
		metadata: {
			name: constants.name
			labels: name: constants.name
			namespace: parameters.namespace
		}
		spec: selfSigned: {}
	}, {
		apiVersion: "cert-manager.io/v1beta1"
		kind:       "Certificate"
		metadata: {
			name: constants.name
			labels: name: constants.name
			namespace: parameters.namespace
		}
		spec: {
			dnsNames: [
				"\(constants.name)",
				"\(constants.name).\(parameters.namespace).svc",
				"\(constants.name).\(parameters.namespace).svc.cluster.local",
			]
			issuerRef: {
				kind: "Issuer"
				name: "\(constants.name)"
			}
			secretName: "\(constants.name)-webhook-server-cert"
		}
	}]
}

_core_items: [
	_serviceAccount,
	_workload,
	_rbac_ClusterRoleBinding,
]

#WorkloadTemplate: {
	kind:       "List"
	apiVersion: "v1"
	items:
		_core_items +
		_iam_clusterAdminAccess +
		_extra_rbac_ClusterRoleAndBinding +
		_extra_certManager_IssuerAndCertificater +
		[]
}

#WorkloadParameters: {
	namespace:      string
	image:          string
	test:           bool
	logviewDomain?: string
	certManager: bool
}

parameters: #WorkloadParameters
template:   #WorkloadTemplate
