GOBIN = $(shell go env GOPATH)/bin

CONTROLLER_GEN ?= $(GOBIN)/controller-gen
IMAGINE ?= $(GOBIN)/imagine
KG ?= $(GOBIN)/kg

ifeq ($(MAKER_CONTAINER),true)
  IMAGINE=imagine
  KG=kg
endif

REGISTRY ?= docker.io/isovalent
imagine_push_or_export = --export
ifeq ($(PUSH),true)
imagine_push_or_export = --push
endif

.buildx_builder:
	# see https://github.com/docker/buildx/issues/308
	mkdir -p ../.buildx
	docker buildx create --platform linux/amd64 > $@

operator-image: .buildx_builder
	$(IMAGINE) build \
		--builder $$(cat .buildx_builder) \
		--base ./ \
		--name gke-test-cluster-operator \
		--registry $(REGISTRY) \
		$(imagine_push_or_export) \
		--cleanup
	$(IMAGINE) image \
		--base ./ \
		--name gke-test-cluster-operator \
		--registry $(REGISTRY) \
		> image-gke-test-cluster-operator.tag

job-runner-init-image: .buildx_builder
	$(IMAGINE) build \
		--builder $$(cat .buildx_builder) \
		--base ./config/job-runner-init-image \
		--name gke-test-cluster-job-runner-init \
		--test \
		--registry $(REGISTRY) \
		$(imagine_push_or_export) \
		--cleanup
	$(IMAGINE) image \
		--base ./config/job-runner-init-image  \
		--name gke-test-cluster-operator \
		--registry $(REGISTRY) \
		> image-gke-test-cluster-job-runner-init.tag

manifests:
	./generate-manifests.sh "$$(cat  image-gke-test-cluster-operator.tag)"

test.controllers-local: operator-image
	docker load -i image-gke-test-cluster-operator.oci
	$(MAKE) test.controllers

test.controllers:
	./test-controllers.sh "$$(cat image-gke-test-cluster-operator.tag)"

test.unit:
	go test ./pkg/...

generate:
	$(CONTROLLER_GEN) object:headerFile=".license_header.go.txt" crd:trivialVersions=false rbac:roleName=gke-test-cluster-operator webhook paths="./..."
