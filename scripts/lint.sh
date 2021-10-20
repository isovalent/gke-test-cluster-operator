#!/bin/bash

# Copyright 2017-2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o pipefail
set -o nounset

MAKER_IMAGE="${MAKER_IMAGE:-quay.io/cilium/image-maker:68cf0628c7a77124319cb8c33434b8cdb42a3865}"

root_dir="$(git rev-parse --show-toplevel)"

if [ -z "${MAKER_CONTAINER+x}" ] ; then
   exec docker run --rm --volume "${root_dir}:/src" --workdir /src "${MAKER_IMAGE}" "/src/scripts/$(basename "${0}")"
fi

cd "${root_dir}"

find ./ -name '*.sh' -not -path './.gopath/**' -exec shellcheck {} +
find ./ -name Dockerfile -not -path './.gopath/**' -exec hadolint {} +
