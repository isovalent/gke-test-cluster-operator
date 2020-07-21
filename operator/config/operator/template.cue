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
						memory: "30Mi"
					}
					requests: {
						cpu:    "100m"
						memory: "20Mi"
					}
				}
			}]
			terminationGracePeriodSeconds: 10
		}
	}
}

_command: [...string]

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
}

WorkloadTemplate :: {
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
		_workload,
	]
}

WorkloadParameters :: {
	namespace: string
	image:     string
	test:      bool
}

parameters: WorkloadParameters
template:   WorkloadTemplate
