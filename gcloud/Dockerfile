# Copyright 2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

ARG CLOUD_SDK_IMAGE=docker.io/google/cloud-sdk:slim@sha256:a2bade78228faad59a16c36d440f10cfef58a6055cd997d19e258c59c78a409a
ARG TESTER_IMAGE=docker.io/cilium/image-tester:70724309b859786e0a347605e407c5261f316eb0@sha256:89cc1f577d995021387871d3dbeb771b75ab4d70073d9bcbc42e532792719781

FROM ${CLOUD_SDK_IMAGE} as builder

COPY install-deps.sh /tmp/install-deps.sh
RUN /tmp/install-deps.sh

FROM ${TESTER_IMAGE} as test
COPY --from=builder / /
COPY test /test
RUN /test/bin/cst

FROM scratch
LABEL maintainer="maintainer@cilium.io"
# work-around for quay.io 400 error: https://gist.github.com/errordeveloper/7a078c53d4694f0cad9a727255fb23a6
COPY empty /tmp/empty
COPY --from=builder / /
CMD [ "/usr/local/bin/init.sh" ]
