#!/bin/bash

# Copyright 2017-2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

set -o xtrace
set -o errexit
set -o pipefail
set -o nounset

if [ "$#" -ne 2 ] ; then
  echo "$0 expects exactly two argument"
fi

service_account="$1"
cluster_name="$2"

until gcloud auth list "--format=value(account)" | grep "${service_account}" ; do sleep 1 ; done

gcloud auth login "${service_account}"
