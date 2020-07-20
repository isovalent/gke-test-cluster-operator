GOBIN := $(shell go env GOPATH)/bin

CONTROLLER_GEN ?= $(GOBIN)/controller-gen
IMAGINE ?= $(GOBIN)/imagine
KG ?= $(GOBIN)/kg

.buildx_builder:
	# see https://github.com/docker/buildx/issues/308
	mkdir -p ../.buildx
	docker buildx create --platform linux/amd64 > $@

image: .buildx_builder
	$(IMAGINE) build \
		--builder $$(cat .buildx_builder) \
		--base ./ \
		--name gke-test-cluster-operator \
		--registry docker.io/isovalent \
		--export \
		--cleanup
	$(IMAGINE) image \
		--base ./ \
		--name gke-test-cluster-operator \
		--registry docker.io/isovalent \
		> image-gke-test-cluster-operator.tag

manifests:
	jq -n \
	  --arg image $$(cat  image-gke-test-cluster-operator.tag) \
	    '.instances=[({} | .output="operator.json" | .parameters=({} | .namespace="kube-system" | .image="\($$image)" | .test=false))]' \
	  > config/operator/instances.json
	$(KG) -input-directory config/operator -output-directory config/operator

test.local: image
	docker load -i image-gke-test-cluster-operator.oci
	./test-controllers.sh "$$(cat image-gke-test-cluster-operator.tag)"

test.unit:
	go test ./pkg/...

generate:
	$(CONTROLLER_GEN) object:headerFile=".license_header.go.txt" crd:trivialVersions=false rbac:roleName=gke-test-cluster-operator webhook paths="./..."
