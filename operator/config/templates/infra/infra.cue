// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package infra

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

_name:      "\(resource.metadata.name)"
_namespace: "\(defaults.metadata.namespace)" | *"\(resource.metadata.namespace)"

_project:  "\(defaults.spec.project)" | *"\(resource.spec.project)"
_location: "\(defaults.spec.location)" | *"\(resource.spec.location)"

_runnerImage:     "\(defaults.spec.jobSpec.runner.image)" | *"\(resource.spec.jobSpec.runner.image)"
_runnerInitImage: "\(defaults.spec.jobSpec.runner.initImage)" | *"\(resource.spec.jobSpec.runner.initImage)"

_runnerCommand: [...string]

if len(resource.spec.jobSpec.runner.command) > 0 {
	_runnerCommand: resource.spec.jobSpec.runner.command
}

_authInfoEnv: [
	{
		name:  "SERVICE_ACCOUNT"
		value: "\(_name)-admin@\(_project).iam.gserviceaccount.com"
	},
	{
		name:  "CLUSTER_LOCATION"
		value: "\(_location)"
	},
	{
		name:  "CLUSTER_NAME"
		value: "\(_name)"
	},
]

_kubeconfigEnv: {
	name:  "KUBECONFIG"
	value: "/credentials/kubeconfig"
}

_extraEnv: [...{}]

if len(resource.spec.jobSpec.runner.env) > 0 {
	_extraEnv: resource.spec.jobSpec.runner.env
}

_kubeconfigVolume: {
	name: "credentials"
	emptyDir: {}
}

_kubeconfigVolumeMount: {
	name:      "credentials"
	mountPath: "/credentials"
}

_systemConfigVolume: {
	name: "config-system"
	configMap: {
		optional: true
		name:     "\(_name)-system"
	}
}
_systemConfigVolumeMount: {
	name:      "config-system"
	mountPath: "/config/system"
}

_extraVolumes: [...{}]

_extraVolumeMounts: [...{}]

if resource.spec.jobSpec.runner.configMap != "" {
	_extraVolumes: [{
		name: "config-user"
		configMap: name: "\(resource.spec.jobSpec.runner.configMap)"
	}]
	_extraVolumeMounts: [{
		name:      "config-user"
		mountPath: "/config/user"
	}]
}

_commonInitContainer: {
	name:         "initutil"
	image:        "\(_runnerInitImage)"
	env:          [_kubeconfigEnv] + _authInfoEnv + _extraEnv
	volumeMounts: [_kubeconfigVolumeMount, _systemConfigVolumeMount] + _extraVolumeMounts
}

_promviewImage: "quay.io/isovalent/gke-test-cluster-promview:7695938dcf3a6e4f0e7fb9537091103259aed46e"

_promviewWorkload: {
	apiVersion: "apps/v1"
	kind:       "Deployment"
	metadata: {
		name: "\(_name)-promview"
		labels: {
			cluster:   "\(_name)"
			component: "promview"
		}
		namespace: "\(_namespace)"
	}
	spec: _promviewWorkloadSpec
}

_promviewWorkloadSpec: {
	selector:
		matchLabels: {
			cluster:   "\(_name)"
			component: "promview"
		}
	template: metadata: {
		labels: {
			cluster:   "\(_name)"
			component: "promview"
		}
		annotations: {
			// do not scrape the pod directly, use service and label seletor
			"prometheus.io.scrape": "false"
		}
	}
	replicas: 2
	template: {
		metadata: labels: {
			cluster:   "\(_name)"
			component: "promview"
		}
		spec: {
			serviceAccountName:           "\(_name)-admin"
			automountServiceAccountToken: false
			enableServiceLinks:           false
			volumes:                      [_kubeconfigVolume, _systemConfigVolume] + _extraVolumes
			initContainers: [_commonInitContainer]
			containers: [{
				name: "promview"
				command: ["/usr/bin/gke-test-cluster-promview"]
				image:        _promviewImage
				env:          [_kubeconfigEnv] + _authInfoEnv + _extraEnv
				volumeMounts: [_kubeconfigVolumeMount, _systemConfigVolumeMount] + _extraVolumeMounts
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

_promviewService: {
	apiVersion: "v1"
	kind:       "Service"
	metadata: {
		name: "\(_name)-promview"
		labels: {
			cluster:   "\(_name)"
			component: "promview"
		}
		namespace: "\(_namespace)"
	}
	spec: {
		selector: {
			cluster:   "\(_name)"
			component: "promview"
		}
		ports: [{
			name:       "promview"
			port:       80
			targetPort: 8080
		}]
	}
}

_testRunnerJob: {
	apiVersion: "batch/v1"
	kind:       "Job"
	metadata: {
		name: "test-runner-\(_name)"
		labels: {
			cluster:   "\(_name)"
			component: "test-runner"
		}
		namespace: "\(_namespace)"
	}
	spec: _testRunnerJobSpec
}

_testRunnerJobSpec: {
	backoffLimit: 0
	template: {
		metadata: labels: cluster: "\(_name)"
		spec: {
			serviceAccountName:           "\(_name)-admin"
			automountServiceAccountToken: false
			enableServiceLinks:           false
			volumes:                      [_kubeconfigVolume, _systemConfigVolume] + _extraVolumes
			initContainers: [_commonInitContainer]
			containers: [{
				name:         "test-runner"
				command:      _runnerCommand
				image:        "\(_runnerImage)"
				env:          [_kubeconfigEnv] + _authInfoEnv + _extraEnv
				volumeMounts: [_kubeconfigVolumeMount, _systemConfigVolumeMount] + _extraVolumeMounts
			}]
			dnsPolicy:     "ClusterFirst"
			restartPolicy: "Never"
		}
	}
}

#TestWorkloadTemplate: {
	kind:       "List"
	apiVersion: "v1"
	items: [
		_testRunnerJob,
		_promviewWorkload,
		_promviewService,
	]
}

defaults: v1alpha1.#TestClusterGKE

resource: v1alpha1.#TestClusterGKE

template: #TestWorkloadTemplate
