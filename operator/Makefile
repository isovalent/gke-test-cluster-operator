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

image: .buildx_builder
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

manifests:
	jq -n \
	  --arg image $$(cat  image-gke-test-cluster-operator.tag) \
	    '.instances=[({} | .output="operator.json" | .parameters=({} | .namespace="kube-system" | .image="\($$image)" | .test=false))]' \
	  > config/operator/instances.json
	$(KG) -input-directory config/operator -output-directory config/operator

test.controllers-local: image
	docker load -i image-gke-test-cluster-operator.oci
	$(MAKE) test.controllers

test.controllers:
	./test-controllers.sh "$$(cat image-gke-test-cluster-operator.tag)"

test.unit:
	go test ./pkg/...

generate:
	$(CONTROLLER_GEN) object:headerFile=".license_header.go.txt" crd:trivialVersions=false rbac:roleName=gke-test-cluster-operator webhook paths="./..."
