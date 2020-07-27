// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package job

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

_name:        "\(resource.metadata.name)"
_namespace:   "\(defaults.metadata.namespace)" | *"\(resource.metadata.namespace)"
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
					volumes: []
					initContainers: [{
						name: "get-credentials"
						command: [
							"bash",
							"-c",
							"until gcloud auth list \"--format=value(account)\" | grep \(_name)-admin@\(_project).iam.gserviceaccount.com ; do sleep 1 ; done",
						]
						image: "google/cloud-sdk:slim@sha256:a2bade78228faad59a16c36d440f10cfef58a6055cd997d19e258c59c78a409a"
						volumeMounts: []
					}]
					containers: [{
						name: "test-runner"
						command: [
							"bash",
							"-l",
						]
						tty:   true
						image: "\(_runnerImage)"
						volumeMounts: []
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
