#!/bin/bash -x

set -o errexit
set -o pipefail
set -o nounset

CILIUM_VERSION="1.9.1"
CNRM_VERSION="1.33.0"

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "${root_dir}/config/test-deps"

curl --fail --show-error --silent --location \
  --output cilium.all.yaml \
  "https://github.com/cilium/cilium/blob/v${CILIUM_VERSION}/install/kubernetes/quick-install.yaml?raw=true"


curl --fail --show-error --silent --location \
  --output cnrm.crd.yaml \
  "https://github.com/GoogleCloudPlatform/k8s-config-connector/blob/${CNRM_VERSION}/install-bundles/install-bundle-namespaced/crds.yaml?raw=true"


