#!/bin/bash

# Copyright 2017-2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o pipefail
set -o nounset

operator_image="${1}"
logview_image="${2}"
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
        "test": true
      }
    },
    {
      "output": "operator.yaml",
      "parameters": {
        "namespace": "${use_namespace}",
        "image": "${operator_image}",
        "test": false,
        "logviewDomain": "${logview_domain}"
      }
    }
  ]
}
EOF

ingress_route_prefix_salt="${INGRESS_ROUTE_PREFIX_SALT:-""}"

if [ -z "${ingress_route_prefix_salt}" ] ; then
  echo "WARNING: INGRESS_ROUTE_PREFIX_SALT is not set, it's required for production"
fi

cat > config/logview/instances.json << EOF
{
  "instances": [
    {
      "output": "test-clusters-atlantis/logview.yaml",
      "parameters": {
        "namespace": "test-clusters-atlantis",
        "image": "${logview_image}",
        "ingressRoutePrefixSalt": "${ingress_route_prefix_salt}"
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
