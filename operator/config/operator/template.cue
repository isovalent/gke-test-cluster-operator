constants: {
	name: "gke-test-cluster-operator"
}

_workload: {
	metadata: {
		name: "\(constants.name)"
		labels: name: "\(constants.name)"
		namespace: "\(parameters.namespace)"
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
		selector: matchLabels: name: "\(constants.name)"
		template: metadata: labels: name: "\(constants.name)"
		replicas: 1
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

	// it's not normally desired to have difference in RBAC configuration
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

#WorkloadTemplate: {
	kind:       "List"
	apiVersion: "v1"
	items: [
		{
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
		_workload,
		_extra_rbac_clusterrole,
		_extra_rbac_clusterrolebiding,
	]
}

#WorkloadParameters: {
	namespace: string
	image:     string
	test:      bool
}

parameters: #WorkloadParameters
template:   #WorkloadTemplate
