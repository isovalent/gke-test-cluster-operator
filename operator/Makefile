GOBIN := $(shell go env GOPATH)/bin

CONTROLLER_GEN ?= $(GOBIN)/controller-gen

image:
	make -C .. operator-image

test.local: image
	docker load -i ../image-gke-test-cluster-operator.oci
	make -C .. test.operator-controllers

test.unit:
	go test ./pkg/...

generate:
	$(CONTROLLER_GEN) object:headerFile=".license_header.go.txt" crd:trivialVersions=false rbac:roleName=gke-test-cluster-operator webhook paths="./..."
