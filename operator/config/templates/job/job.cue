// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package job

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

_project: "cilium-ci"

_name:            "\(resource.metadata.name)"
_namespace:       "\(defaults.metadata.namespace)" | *"\(resource.metadata.namespace)"
_location:        "\(defaults.spec.location)" | *"\(resource.spec.location)"
_runnerImage:     "\(defaults.spec.jobSpec.runner.image)" | *"\(resource.spec.jobSpec.runner.image)"
_runnerInitImage: "\(defaults.spec.jobSpec.runner.initImage)" | *"\(resource.spec.jobSpec.runner.initImage)"

_runnerCommand: [...string]

if len(resource.spec.jobSpec.runner.command) > 0 {
	_runnerCommand: resource.spec.jobSpec.runner.command
}

_authInfoEnv: [
	{
		name: "SERVICE_ACCOUNT"
		value: "\(_name)-admin@\(_project).iam.gserviceaccount.com"
	},
	{
		name: "CLUSTER_LOCATION"
		value: "\(_location)"
	},
	{
		name: "CLUSTER_NAME"
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

_extraVolumes: [...{}]

_extraVolumeMounts: [...{}]

if resource.spec.jobSpec.runner.configMap != "" {
	_extraVolumes: [{
		name: "config"
		configMap: name: "\(resource.spec.jobSpec.runner.configMap)"
	}]
	_extraVolumeMounts: [{
		name:      "config"
		mountPath: "/config"
	}]
}

#JobTemplate: {
	kind:       "List"
	apiVersion: "v1"
	items: [{
		apiVersion: "batch/v1"
		kind:       "Job"
		metadata: {
			name: "test-runner-\(_name)"
			labels: cluster: "\(_name)"
			namespace: "\(_namespace)"
		}
		spec: {
			backoffLimit: 0
			template: {
				metadata:
					labels:
						cluster: "\(_name)"
				spec: {
					volumes: [_kubeconfigVolume] + _extraVolumes
					initContainers: [{
						name:         "init-runner"
						image:        "\(_runnerInitImage)"
						env:          [_kubeconfigEnv] + _authInfoEnv + _extraEnv
						volumeMounts: [_kubeconfigVolumeMount] + _extraVolumeMounts
					}]
					containers: [{
						name:         "test-runner"
						command:      _runnerCommand
						image:        "\(_runnerImage)"
						env:          [_kubeconfigEnv] + _authInfoEnv + _extraEnv
						volumeMounts: [_kubeconfigVolumeMount] + _extraVolumeMounts
					}]
					dnsPolicy:          "ClusterFirst"
					restartPolicy:      "Never"
					serviceAccountName: "\(_name)-admin"
				}
			}
		}
	}]
}

defaults: v1alpha1.#TestClusterGKE

resource: v1alpha1.#TestClusterGKE

template: #JobTemplate
