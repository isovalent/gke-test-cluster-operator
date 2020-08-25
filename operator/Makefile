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
	# see https://github.com/docker/buildx/issues/308
	mkdir -p ../.buildx
	docker buildx create --platform linux/amd64 > $@

images.all: images.operator images.job-runner-init images.logview images.promview images.requester

images.operator: .buildx_builder
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

images.job-runner-init: .buildx_builder
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
		--name gke-test-cluster-job-runner-init \
		--registry $(REGISTRY) \
		> image-gke-test-cluster-job-runner-init.tag

images.logview: .buildx_builder
	$(IMAGINE) build \
		--builder $$(cat .buildx_builder) \
		--base ./logview \
		--name gke-test-cluster-logview \
		--registry $(REGISTRY) \
		$(imagine_push_or_export) \
		--cleanup
	$(IMAGINE) image \
		--base ./logview \
		--name gke-test-cluster-logview \
		--registry $(REGISTRY) \
		> image-gke-test-cluster-logview.tag

images.promview: .buildx_builder
	$(IMAGINE) build \
		--builder $$(cat .buildx_builder) \
		--base ./promview \
		--name gke-test-cluster-promview \
		--registry $(REGISTRY) \
		$(imagine_push_or_export) \
		--cleanup
	$(IMAGINE) image \
		--base ./promview \
		--name gke-test-cluster-promview \
		--registry $(REGISTRY) \
		> image-gke-test-cluster-promview.tag

images.requester: .buildx_builder
	$(IMAGINE) build \
		--builder $$(cat .buildx_builder) \
		--base ./ \
		--dockerfile ./requester/Dockerfile \
		--name gke-test-cluster-requester \
		--registry $(REGISTRY) \
		$(imagine_push_or_export) \
		--cleanup
	$(IMAGINE) image \
		--base ./ \
		--name gke-test-cluster-requester \
		--registry $(REGISTRY) \
		> image-gke-test-cluster-requester.tag

manifests.generate:
	./generate-manifests.sh "$$(cat  image-gke-test-cluster-operator.tag)" "$$(cat image-gke-test-cluster-logview.tag)"

manifests.promote:
	cp config/operator/operator.yaml ../config/kube-system/operator.yaml
	cp config/crd/*.yaml config/rbac/*.yaml ../config/kube-system/
	for ns_dir in config/logview/*/ ; do cp $${ns_dir}/logview.yaml ../config/$$(basename $${ns_dir})/logview.yaml ; done

test.controllers-local: images.all
	docker load -i image-gke-test-cluster-operator.oci
	$(MAKE) test.controllers

test.controllers:
	./test-controllers.sh "$$(cat image-gke-test-cluster-operator.tag)" "$$(cat image-gke-test-cluster-logview.tag)"

test.unit:
	go test ./pkg/...

misc.generate:
	$(CONTROLLER_GEN) object:headerFile=".license_header.go.txt" crd:trivialVersions=false object paths="./api/..."
	$(CONTROLLER_GEN) object:headerFile=".license_header.go.txt" rbac:roleName=gke-test-cluster-operator webhook paths="./..."
	go generate ./api/...
	go generate ./pkg/...

.PHONY: kind-cluster
kind-cluster:
	if ! kind get clusters | grep $(KIND_CLUSTER_NAME) ; then \
		kind create cluster --name $(KIND_CLUSTER_NAME) ; \
	fi
	kind export kubeconfig --name $(KIND_CLUSTER_NAME)
	kubectl apply -f ../config/cnrm-system/crds.yaml

test.controllers-kind: kind-cluster images.operator
	kind load image-archive image-gke-test-cluster-operator.oci --name $(KIND_CLUSTER_NAME)
	$(MAKE) test.controllers
