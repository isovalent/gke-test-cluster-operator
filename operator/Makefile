GOBIN := $(shell go env GOPATH)/bin

CONTROLLER_GEN ?= $(GOBIN)/controller-gen
KG ?= $(GOBIN)/kg

image:
	make -C .. operator-image

manifests:
	jq \
	  --arg image $$(cat  ../image-gke-test-cluster-operator.tag) \
	  -n '.instances=[{} | .output="operator.json" | .parameters=({} | .namespace="kube-system" | .image="\($$image)")]' \
	   > config/operator/instances.json
	$(KG) -input-directory config/operator

test.local: image
	docker load -i ../image-gke-test-cluster-operator.oci
	make -C .. test.operator-controllers

test.unit:
	go test ./pkg/...

generate:
	$(CONTROLLER_GEN) object:headerFile=".license_header.go.txt" crd:trivialVersions=false rbac:roleName=gke-test-cluster-operator webhook paths="./..."
