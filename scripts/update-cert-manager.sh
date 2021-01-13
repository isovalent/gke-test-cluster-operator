#!/bin/bash -x

set -o errexit
set -o pipefail
set -o nounset

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "${root_dir}/config/cert-manager"

curl --fail --show-error --silent --location --remote-name \
  https://github.com/jetstack/cert-manager/releases/latest/download/cert-manager.yaml

