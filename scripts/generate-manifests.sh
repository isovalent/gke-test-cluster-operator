#!/bin/bash

# Copyright 2017-2021 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o pipefail
set -o nounset

if [ "$#" -ne 1 ] ; then
  echo "$0 requires exactly 1 argument - operator image"
  exit 1
fi

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "${root_dir}"

operator_image="${1}"
use_namespace="${NAMESPACE:-kube-system}"

logview_domain="${LOGVIEW_DOMAIN:-""}"

if [ -z "${logview_domain}" ] ; then
  echo "WARNING: LOGVIEW_DOMAIN is not set, it's required for production"
fi

cat > config/operator/instances.json << EOF
{
  "instances": [
    {
      "output": "operator-test.yaml",
      "parameters": {
        "namespace": "${use_namespace}",
        "image": "${operator_image}",
        "test": true,
        "certManager": true
      }
    },
    {
      "output": "operator.yaml",
      "parameters": {
        "namespace": "${use_namespace}",
        "image": "${operator_image}",
        "test": false,
        "logviewDomain": "${logview_domain}",
        "certManager": true
      }
    }
  ]
}
EOF

if [ -n "${GOPATH+x}" ] ; then
  export PATH="${PATH}:${GOPATH}/bin"
fi

kg -input-directory ./config/operator -output-directory ./config/operator
kg -input-directory ./config/templates/prom -output-directory ./config/templates/prom
