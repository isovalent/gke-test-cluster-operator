# syntax=docker/dockerfile:1.1-experimental

# Copyright 2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

ARG GOLANG_IMAGE=docker.io/library/golang:1.14@sha256:ede9a57fa6d862ab87f5abcea707c3d55e445ff01d806334a1cb7aae45ec73bb
ARG GCLOUD_IMAGE=quay.io/isovalent/gke-test-cluster-gcloud:803ff83d3786eb38ef05c95768060b0c7ae0fc4d

FROM ${GOLANG_IMAGE} as builder

WORKDIR /src

RUN --mount=type=bind,target=/src --mount=target=/root/.cache,type=cache --mount=target=/go/pkg/mod,type=cache \
    go mod download

RUN mkdir -p /out/usr/bin

RUN --mount=type=bind,target=/src --mount=target=/root/.cache,type=cache --mount=target=/go/pkg/mod,type=cache \
    go vet ./...

RUN --mount=type=bind,target=/src --mount=target=/root/.cache,type=cache --mount=target=/go/pkg/mod,type=cache \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -ldflags '-s -w' -o /out/usr/bin/gke-test-cluster-promview ./

FROM ${GCLOUD_IMAGE}
COPY --from=builder /out /

ENTRYPOINT ["/usr/bin/gke-test-cluster-promview"]
