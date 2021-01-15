#!/bin/bash

# Copyright 2017-2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o pipefail
set -o nounset

if [ "$#" -ne 3 ] ; then
  echo "$0 requires exactly 2 arguments - project, namespace & logview image"
  exit 1
fi

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "${root_dir}"

project="${1}"
namespace="${2}"
logview_image="${3}"

ingress_route_prefix_salt="${INGRESS_ROUTE_PREFIX_SALT:-""}"

if [ -z "${ingress_route_prefix_salt}" ] ; then
  echo "WARNING: INGRESS_ROUTE_PREFIX_SALT is not set, it's required for production"
fi

cat > config/logview/instances.cue << EOF
package logview

instances: [{
  output: "logview.yaml"
  parameters: {
    namespace: "${namespace}"
    image: "${logview_image}"
    ingressRoutePrefixSalt: "${ingress_route_prefix_salt}"
  }
}]
EOF

cat > config/requester-access/instances.cue << EOF
package requesteraccess

instances: [{
  output: "access.yaml"
  parameters: {
    namespace: "${namespace}"
    project: "${project}"
  }
}]
EOF

if [ -n "${GOPATH+x}" ] ; then
  export PATH="${PATH}:${GOPATH}/bin"
fi

kg -input-directory ./config/logview -output-directory "ns-${namespace}"
kg -input-directory ./config/requester-access -output-directory "ns-${namespace}"

cat > "ns-${namespace}/0-namespace.yaml" << EOF
apiVersion: v1
kind: Namespace
metadata:
  name: ${namespace}
EOF
