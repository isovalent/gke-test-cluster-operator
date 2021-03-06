#!/bin/bash

# Copyright 2017-2021 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o pipefail
set -o nounset

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "${root_dir}"

name="gke-test-cluster-operator"
namespace="${name}-$(date +%s)"

# this is a simple script for running controller tests on kind,
# it could eventually become a Go program

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
  kubectl delete --wait="false" --filename="config/operator/operator-test.yaml"
  # test namespaces are not deleted on test failure, so make sure those are cleaned up here
  kubectl delete namespace --selector="test=${namespace}" --field-selector="status.phase=Active" --wait="false"
  exit "$1"
}

one_pod_running() {
  # shellcheck disable=SC2207
  pods=($(kubectl get pods --namespace="${namespace}" --selector="name=${name}" --output="jsonpath={range .items[?(.status.phase == \"Running\")]}{.metadata.name}{\"\n\"}{end}"))
  test "${#pods[@]}" -eq 1
}

one_pod_failed() {
  # shellcheck disable=SC2207
  pods=($(kubectl get pods --namespace="${namespace}" --selector="name=${name}" --output="jsonpath={range .items[?(.status.phase == \"Failed\")]}{.metadata.name}{\"\n\"}{end}"))
  test "${#pods[@]}" -eq 1
}

wait_for_cert_manager() {
  echo "INFO: waiting for cert-manager to start..."
  kubectl wait deployment --namespace="cert-manager" --for="condition=Available" cert-manager-webhook cert-manager-cainjector cert-manager --timeout=3m
  kubectl wait pods --namespace="cert-manager" --for="condition=Ready" --all --timeout=3m
  kubectl wait apiservice --for="condition=Available" v1beta1.cert-manager.io v1beta1.acme.cert-manager.io --timeout=3m
  until kubectl get secret --namespace="cert-manager" cert-manager-webhook-ca 2> /dev/null ; do sleep 0.5 ; done
}

wait_for_pod() {
  echo "INFO: waiting for the test job to start..."
  # kubectl wait job only supports two conditions - complete or failed,
  # so wait for a pod to get scheduled
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

wait_for_cert_manager

echo "INFO: cleaning up stale resources"

kubectl delete namespace --selector="test" --field-selector="status.phase=Active" --wait="false"
kubectl delete ClusterRole,ClusterRoleBinding,MutatingWebhookConfiguration,ValidatingWebhookConfiguration --selector="name=${name}"

echo "INFO: creating test job"

# since we rely on external dependencies, CI uses a different
# namespace, which shouldn't have to be the case once this
# is re-written in Go and CUE template is rendered directly
if [ -z "${CI+x}" ] ; then
   NAMESPACE="${namespace}" ./scripts/generate-manifests.sh "${@}"
   kubectl create namespace "${namespace}"
   kubectl label namespace "${namespace}" test="${name}"
else
   namespace="kube-system"
fi

kubectl apply --filename="config/rbac/role.yaml"
kubectl apply --filename="config/crd"

# wait_for_cert_manager attempts what it can, yet the readiness of webhook is hard to determined without retrying
until kubectl apply --namespace="${namespace}" --filename="config/operator/operator-test.yaml" ; do sleep 0.5 ; done

wait_for_pod

follow_logs

if ! container_status ; then
  troubleshoot
  bail 1
fi
bail 0
