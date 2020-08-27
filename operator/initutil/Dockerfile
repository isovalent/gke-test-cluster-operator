# Copyright 2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

ARG GCLOUD_IMAGE=quay.io/isovalent/gke-test-cluster-gcloud:803ff83d3786eb38ef05c95768060b0c7ae0fc4d
ARG TESTER_IMAGE=docker.io/cilium/image-tester:70724309b859786e0a347605e407c5261f316eb0@sha256:89cc1f577d995021387871d3dbeb771b75ab4d70073d9bcbc42e532792719781

FROM ${GCLOUD_IMAGE} as builder

COPY init.sh /usr/local/bin/init.sh

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
