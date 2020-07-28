#!/bin/bash

# Copyright 2017-2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

set -o xtrace
set -o errexit
set -o pipefail
set -o nounset

if [ "$#" -ne 3 ] ; then
  echo "$0 expects exactly three argument"
  exit 1
fi

service_account="$1"
cluster_name="$2"
cluster_location="$3"

until gcloud auth list "--format=value(account)" | grep "${service_account}" ; do sleep 1 ; done

gcloud config set account "${service_account}"

gcloud container clusters get-credentials --zone "${cluster_location}" "${cluster_name}"
