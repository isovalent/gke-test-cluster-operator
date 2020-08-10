#!/bin/bash

# Copyright 2017-2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

set -o xtrace
set -o errexit
set -o pipefail
set -o nounset

if [ -z "${SERVICE_ACCOUNT+x}" ]; then
  echo "SERVICE_ACCOUNT must be set"
fi

if [ -z "${CLUSTER_LOCATION+x}" ]; then
  echo "CLUSTER_LOCATION must be set"
fi

if [ -z "${CLUSTER_NAME+x}" ]; then
  echo "CLUSTER_NAME must be set"
fi

until gcloud auth list "--format=value(account)" | grep "${SERVICE_ACCOUNT}" ; do sleep 1 ; done

gcloud config set account "${SERVICE_ACCOUNT}"

gcloud container clusters get-credentials --zone "${CLUSTER_LOCATION}" "${CLUSTER_NAME}"

if [ -f /config/init-manifest ] ; then
  kubectl apply -f /config/init-manifest
fi
