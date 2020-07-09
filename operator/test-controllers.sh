#!/bin/bash

# Copyright 2017-2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o pipefail
set -o nounset

image="${1}"

name="test-gke-test-cluster-operator-controllers"
namespace="${name}-$(date +%s)"

# this is a simple script for runnig controller tests on kind,
# it could eventually become a Go program

echo "INFO: creating test job"

kubectl create --namespace "${namespace}" --filename - <<EOF
apiVersion: v1
kind: List
items:
  - apiVersion: v1
    kind: Namespace
    metadata:
      name: "${namespace}"
      labels:
        test: "${name}"
  - apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: "${name}"
      namespace: "${namespace}"
      labels:
        test: "${name}"
  - apiVersion: rbac.authorization.k8s.io/v1beta1
    kind: ClusterRoleBinding
    metadata:
      name: "${namespace}"
      namespace: "${namespace}"
      labels:
        test: "${name}"
    roleRef:
      kind: ClusterRole
      name: cluster-admin
      apiGroup: rbac.authorization.k8s.io
    subjects:
      - kind: ServiceAccount
        name: "${name}"
        namespace: "${namespace}"
  - apiVersion: batch/v1
    kind: Job
    metadata:
      name: "${name}"
      namespace: "${namespace}"
      labels:
        test: "${name}"
    spec:
      backoffLimit: 0
      template:
        metadata:
          labels:
            test: "${name}"
        spec:
          serviceAccount: "${name}"
          tolerations:
            - effect: NoSchedule
              operator: Exists
          restartPolicy: Never
          dnsPolicy: ClusterFirst
          volumes:
          - name: tmp
            emptyDir: {}
          containers:
            - name: test
              image: "${image}"
              imagePullPolicy: Never
              command:
              - test.gke-test-cluster-operator-controllers
              - -test.v
              - -test.timeout=5m
              - -resource-prefix=\$(POD_NAME)
              - -crd-path=/config/crd/bases
              volumeMounts:
              - mountPath: /tmp
                name: tmp
              env:
              - name: POD_NAME
                valueFrom:
                  fieldRef:
                    fieldPath: metadata.name
EOF

# these tests should quick and there is no point in streaming logs,
# completion should be within a minute or two, otherwise there is
# a problem with the tests

follow_logs() {
  echo "INFO: streaming container logs"
  kubectl logs --namespace="${namespace}" --selector="test=${name}" --follow || true
}

troubleshoot() {
  echo "INFO: gatherig additional information that maybe usefull in troubleshooting the failure"
  kubectl describe pods --namespace="${namespace}" --selector="test=${name}"
}

bail() {
  echo "INFO: cleaning up..."
  kubectl delete namespace "${namespace}" --wait="false"
  kubectl delete clusterrolebinding "${namespace}" --wait="false"
  exit "$1"
}

one_pod_running() {
  pods=($(kubectl get pods --namespace="${namespace}" --selector="test=${name}" --output="jsonpath={range .items[?(.status.phase == \"Running\")]}{.metadata.name}{\"\n\"}{end}"))
  test "${#pods[@]}" -eq 1
}

one_pod_failed() {
  pods=($(kubectl get pods --namespace="${namespace}" --selector="test=${name}" --output="jsonpath={range .items[?(.status.phase == \"Failed\")]}{.metadata.name}{\"\n\"}{end}"))
  test "${#pods[@]}" -eq 1
}

wait_for_pod() {
  echo "INFO: waiting for the test job to start..."
  # kubectl wait job only supports two condtions - complete or failed,
  # so wait for a pod to get schulded
  until kubectl wait pods --namespace="${namespace}" --selector="test=${name}" --for="condition=PodScheduled" --timeout="2m" 2> /dev/null ; do sleep 0.5 ; done
  # kubectl wait doesn't support multiple conditions, and `wait -n` is not
  # available in all common versions of bash, so poll pod status instead
  until one_pod_running || one_pod_failed ; do
    kubectl get pods --namespace="${namespace}" --selector="test=${name}" --show-kind --no-headers
    sleep 0.5
  done
}

container_status() {
  echo "INFO: getting container status..."
  exit_code="$(kubectl get pods --namespace="${namespace}" --selector="test=${name}" --output="jsonpath={.items[0].status.containerStatuses[0].state.terminated.exitCode}")"
  echo "INFO: container existed with ${exit_code}"
  return "${exit_code}"
}

wait_for_pod

follow_logs

if ! container_status ; then
  troubleshoot
  bail 1
fi
bail 0
