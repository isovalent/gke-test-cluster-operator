// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package infra

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha2"

_generatedName: resource.metadata.name | *resource.status.clusterName

_namespace: defaults.metadata.namespace | *resource.metadata.namespace

_project:  defaults.spec.project | *resource.spec.project
_location: defaults.spec.location | *resource.spec.location

_runnerImage:     defaults.spec.jobSpec.runner.image | *resource.spec.jobSpec.runner.image
_runnerInitImage: defaults.spec.jobSpec.runner.initImage | *resource.spec.jobSpec.runner.initImage

_promviewLabels: {
	cluster:   resource.metadata.name
	component: "promview"
}

_runnerLabels: {
	cluster:   resource.metadata.name
	component: "test-runner"
}

_runnerCommand: [...string]

if len(resource.spec.jobSpec.runner.command) > 0 {
	_runnerCommand: resource.spec.jobSpec.runner.command
}

_authInfoEnv: [
	{
		name:  "SERVICE_ACCOUNT"
		value: "\(_generatedName)-admin@\(_project).iam.gserviceaccount.com"
	},
	{
		name:  "CLUSTER_LOCATION"
		value: _location
	},
	{
		name:  "CLUSTER_NAME"
		value: _generatedName
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
		name:     "\(_generatedName)-system"
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
		configMap: name: resource.spec.jobSpec.runner.configMap
	}]
	_extraVolumeMounts: [{
		name:      "config-user"
		mountPath: "/config/user"
	}]
}

_commonInitContainer: {
	name:         "initutil"
	image:        _runnerInitImage
	env:          [_kubeconfigEnv] + _authInfoEnv + _extraEnv
	volumeMounts: [_kubeconfigVolumeMount, _systemConfigVolumeMount] + _extraVolumeMounts
}

_promviewImage: "quay.io/isovalent/gke-test-cluster-promview:7695938dcf3a6e4f0e7fb9537091103259aed46e"

_promviewWorkload: {
	apiVersion: "apps/v1"
	kind:       "Deployment"
	metadata: {
		name:      "\(_generatedName)-promview"
		labels:    _promviewLabels
		namespace: _namespace
	}
	spec: _promviewWorkloadSpec
}

_promviewWorkloadSpec: {
	selector:
		matchLabels: _promviewLabels
	template: metadata: {
		labels: _promviewLabels
		annotations: {
			// do not scrape the pod directly, use service and label seletor
			"prometheus.io.scrape": "false"
		}
	}
	replicas: 2
	template: {
		metadata:
			labels: _promviewLabels
		spec: {
			serviceAccountName:           "\(_generatedName)-admin"
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
		name:      "\(_generatedName)-promview"
		labels:    _promviewLabels
		namespace: _namespace
	}
	spec: {
		selector: _promviewLabels
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
		name:      "test-runner-\(_generatedName)"
		labels:    _runnerLabels
		namespace: _namespace
	}
	spec: _testRunnerJobSpec
}

_testRunnerJobSpec: {
	backoffLimit: 0
	template: {
		metadata: labels: _runnerLabels
		spec: {
			serviceAccountName:           "\(_generatedName)-admin"
			automountServiceAccountToken: false
			enableServiceLinks:           false
			volumes:                      [_kubeconfigVolume, _systemConfigVolume] + _extraVolumes
			initContainers: [_commonInitContainer]
			containers: [{
				name:         "test-runner"
				command:      _runnerCommand
				image:        _runnerImage
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

defaults: v1alpha2.#TestClusterGKE

resource: v1alpha2.#TestClusterGKE

template: #TestWorkloadTemplate
