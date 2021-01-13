GOBIN = $(shell go env GOPATH)/bin

CONTROLLER_GEN ?= $(GOBIN)/controller-gen
IMAGINE ?= $(GOBIN)/imagine
KG ?= $(GOBIN)/kg

ifeq ($(MAKER_CONTAINER),true)
  IMAGINE=imagine
  KG=kg
endif

REGISTRY ?= quay.io/isovalent
imagine_push_or_export = --export
ifeq ($(PUSH),true)
imagine_push_or_export = --push
endif
KIND_CLUSTER_NAME = operator-kind

.buildx_builder:
	docker buildx create --platform linux/amd64 > $@

images.all: images.operator images.initutil images.logview images.promview images.requester

images.operator: .buildx_builder
	$(IMAGINE) build \
		--builder $$(cat .buildx_builder) \
		--base ./ \
		--name gke-test-cluster-operator \
		--upstream-branch origin/main \
		--registry $(REGISTRY) \
		$(imagine_push_or_export)
	$(IMAGINE) image \
		--base ./ \
		--name gke-test-cluster-operator \
		--upstream-branch origin/main \
		--registry $(REGISTRY) \
		> image-gke-test-cluster-operator.tag

images.gcloud: .buildx_builder
	$(IMAGINE) build \
		--builder $$(cat .buildx_builder) \
		--base ./gcloud \
		--name gke-test-cluster-gcloud \
		--test \
		--upstream-branch origin/main \
		--registry $(REGISTRY) \
		$(imagine_push_or_export)
	$(IMAGINE) image \
		--base ./gcloud  \
		--name gke-test-cluster-gcloud \
		--upstream-branch origin/main \
		--registry $(REGISTRY) \
		> image-gke-test-cluster-gcloud.tag

images.initutil: .buildx_builder
	$(IMAGINE) build \
		--builder $$(cat .buildx_builder) \
		--base ./initutil \
		--name gke-test-cluster-initutil \
		--test \
		--upstream-branch origin/main \
		--registry $(REGISTRY) \
		$(imagine_push_or_export)
	$(IMAGINE) image \
		--base ./initutil  \
		--name gke-test-cluster-initutil \
		--upstream-branch origin/main \
		--registry $(REGISTRY) \
		> image-gke-test-cluster-initutil.tag

images.logview: .buildx_builder
	$(IMAGINE) build \
		--builder $$(cat .buildx_builder) \
		--base ./logview \
		--name gke-test-cluster-logview \
		--upstream-branch origin/main \
		--registry $(REGISTRY) \
		$(imagine_push_or_export)
	$(IMAGINE) image \
		--base ./logview \
		--name gke-test-cluster-logview \
		--upstream-branch origin/main \
		--registry $(REGISTRY) \
		> image-gke-test-cluster-logview.tag

images.promview: .buildx_builder
	$(IMAGINE) build \
		--builder $$(cat .buildx_builder) \
		--base ./promview \
		--name gke-test-cluster-promview \
		--upstream-branch origin/main \
		--registry $(REGISTRY) \
		$(imagine_push_or_export)
	$(IMAGINE) image \
		--base ./promview \
		--name gke-test-cluster-promview \
		--upstream-branch origin/main \
		--registry $(REGISTRY) \
		> image-gke-test-cluster-promview.tag

images.requester: .buildx_builder
	$(IMAGINE) build \
		--builder $$(cat .buildx_builder) \
		--base ./ \
		--dockerfile ./requester/Dockerfile \
		--name gke-test-cluster-requester \
		--upstream-branch origin/main \
		--registry $(REGISTRY) \
		$(imagine_push_or_export)
	$(IMAGINE) image \
		--base ./ \
		--name gke-test-cluster-requester \
		--upstream-branch origin/main \
		--registry $(REGISTRY) \
		> image-gke-test-cluster-requester.tag

manifests.generate:
	./generate-manifests.sh "$$(cat  image-gke-test-cluster-operator.tag)" "$$(cat image-gke-test-cluster-logview.tag)"

manifests.promote:
	cp config/operator/operator.yaml ../config/kube-system/operator.yaml
	cp config/templates/prom/prom.yaml ../config/prom/prom.yaml
	cp config/crd/*.yaml config/rbac/*.yaml ../config/kube-system/
	for ns_dir in config/logview/*/ ; do cp $${ns_dir}/logview.yaml ../config/$$(basename $${ns_dir})/logview.yaml ; done

test.controllers-local: images.all
	docker load -i image-gke-test-cluster-operator.oci
	$(MAKE) test.controllers

test.controllers:
	./scripts/test-controllers.sh "$$(cat image-gke-test-cluster-operator.tag)" "$$(cat image-gke-test-cluster-logview.tag)"

test.unit:
	go test ./pkg/...

misc.generate:
	$(CONTROLLER_GEN) object:headerFile=".license_header.go.txt" crd:trivialVersions=false object paths="./api/..."
	$(CONTROLLER_GEN) object:headerFile=".license_header.go.txt" rbac:roleName=gke-test-cluster-operator webhook paths="./..."
	go generate ./api/...
	go generate ./pkg/...
