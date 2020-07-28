// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package job

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

_name:        "\(resource.metadata.name)"
_namespace:   "\(defaults.metadata.namespace)" | *"\(resource.metadata.namespace)"
_location:    "\(defaults.spec.location)" | *"\(resource.spec.location)"
_runnerImage: "\(defaults.spec.jobSpec.runnerImage)" | *"\(resource.spec.jobSpec.runnerImage)"

_project: "cilium-ci"

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
					volumes: [
						{
							name: "credentials"
							emptyDir: {}
						}
					]
					initContainers: [{
						name: "get-credentials"
						command: [
							"gcloud-auth-init.sh",
							"\(_name)-admin@\(_project).iam.gserviceaccount.com",
							"\(_name)",
							"\(_location)",
						]
						image: "docker.io/errordeveloper/gke-test-cluster-job-runner-init:1b1b875acb5fa546f9bf827f73c615f7db4f28dd"
						env: [
							{
								name: "KUBECONFIG"
								value: "/credentials/kubeconfig"
							}
						]
						volumeMounts: [
							{
								name: "credentials"
								mountPath: "/credentials"
							}
						]
					}]
					containers: [{
						name: "test-runner"
						command: [
							"bash",
							"-l",
						]
						tty:   true
						image: "\(_runnerImage)"
						env: [
							{
								name: "KUBECONFIG"
								value: "/credentials/kubeconfig"
							}
						]
						volumeMounts: [
							{
								name: "credentials"
								mountPath: "/credentials"
							}
						]
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
