#!/bin/bash

# Copyright 2017-2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o pipefail
set -o nounset

image="${1}"
use_namespace="${NAMESPACE:-kube-system}"

cat > config/operator/instances.json << EOF
{
  "instances": [
    {
      "output": "operator-test.yaml",
      "parameters": {
        "namespace": "${use_namespace}",
        "image": "${image}",
        "test": true
      }
    },
    {
      "output": "operator.yaml",
      "parameters": {
        "namespace": "${use_namespace}",
        "image": "${image}",
        "test": false
      }
    }
  ]
}
EOF

cat > config/logview/instances.json << EOF
{
  "instances": [
    {
      "output": "test-clusters-atlantis/logview.yaml",
      "parameters": {
        "namespace": "test-clusters-atlantis",
        "image": "${image}"
      }
    }
  ]
}
EOF

if [ -n "${GOPATH+x}" ] ; then
  export PATH="${PATH}:${GOPATH}/bin"
fi

kg -input-directory config/operator -output-directory config/operator
kg -input-directory config/logview -output-directory config/logview
