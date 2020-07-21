#!/bin/bash

# Copyright 2017-2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o pipefail
set -o nounset

image="${1}"

name="gke-test-cluster-operator"
namespace="${name}-$(date +%s)"

# this is a simple script for runnig controller tests on kind,
# it could eventually become a Go program

echo "INFO: creating test job"

cat > config/operator/instances.json << EOF
{
  "instances": [{
      "output": "operator-test.json",
      "parameters": {
        "namespace": "${namespace}",
        "image": "${image}",
        "test": true
      }
  }]
}
EOF

if [ -n "${GOPATH+x}" ] ; then
  export PATH="${PATH}:${GOPATH}/bin"
fi

kg -input-directory config/operator -output-directory config/operator

kubectl apply --filename="config/rbac/role.yaml"

kubectl create namespace "${namespace}"
kubectl label namespace "${namespace}" test="${name}"

kubectl create --namespace "${namespace}" --filename="config/operator/operator-test.json"

# these tests should quick and there is no point in streaming logs,
# completion should be within a minute or two, otherwise there is
# a problem with the tests

follow_logs() {
  echo "INFO: streaming container logs"
  kubectl logs --namespace="${namespace}" --selector="name=${name}" --follow || true
}

troubleshoot() {
  echo "INFO: gatherig additional information that maybe usefull in troubleshooting the failure"
  kubectl describe pods --namespace="${namespace}" --selector="name=${name}"
}

bail() {
  echo "INFO: cleaning up..."
  kubectl delete --wait="false"  --filename="config/operator/operator-test.json"
  # test namespaces are not deleted on test failure, so make sure those are cleaned up here
  kubectl delete namespace --selector="test=${namespace}" --field-selector="status.phase=Active" --wait="false"
  exit "$1"
}

one_pod_running() {
  pods=($(kubectl get pods --namespace="${namespace}" --selector="name=${name}" --output="jsonpath={range .items[?(.status.phase == \"Running\")]}{.metadata.name}{\"\n\"}{end}"))
  test "${#pods[@]}" -eq 1
}

one_pod_failed() {
  pods=($(kubectl get pods --namespace="${namespace}" --selector="name=${name}" --output="jsonpath={range .items[?(.status.phase == \"Failed\")]}{.metadata.name}{\"\n\"}{end}"))
  test "${#pods[@]}" -eq 1
}

wait_for_pod() {
  echo "INFO: waiting for the test job to start..."
  # kubectl wait job only supports two condtions - complete or failed,
  # so wait for a pod to get schulded
  # TODO: this doesn't check if job failed to create pods, which can happen
  # if there is configuration issue
  until kubectl wait pods --namespace="${namespace}" --selector="name=${name}" --for="condition=PodScheduled" --timeout="2m" 2> /dev/null ; do sleep 0.5 ; done
  # kubectl wait doesn't support multiple conditions, and `wait -n` is not
  # available in all common versions of bash, so poll pod status instead
  until one_pod_running || one_pod_failed ; do
    kubectl get pods --namespace="${namespace}" --selector="name=${name}" --show-kind --no-headers
    sleep 0.5
  done
}

get_container_exit_code() {
  kubectl get pods --namespace="${namespace}" --selector="name=${name}" --output="jsonpath={.items[0].status.containerStatuses[0].state.terminated.exitCode}"
}

container_status() {
  echo "INFO: getting container status..."
  # sometimes the value doesn't parse as a number, before this is re-written in Go
  # it will have to be done like this
  until test -n "$(get_container_exit_code)" ; do sleep 0.5 ; done
  exit_code="$(get_container_exit_code)"
  echo "INFO: container exited with ${exit_code}"
  return "${exit_code}"
}

wait_for_pod

follow_logs

if ! container_status ; then
  troubleshoot
  bail 1
fi
bail 0
