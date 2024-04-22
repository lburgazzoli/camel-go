
CONTAINER_REGISTRY ?= quay.io
CONTAINER_REGISTRY_REPOSITORY ?= lburgazzoli/camel-go
CONTAINER_TAG ?= latest
CONTAINER_IMAGE ?= $(CONTAINER_REGISTRY)/$(CONTAINER_REGISTRY_REPOSITORY):$(CONTAINER_TAG)

WASM_CONTAINER_REGISTRY ?= quay.io
WASM_CONTAINER_REGISTRY_REPOSITORY ?= lburgazzoli/camel-go-wasm
WASM_CONTAINER_TAG ?= latest
WASM_CONTAINER_IMAGE ?= $(WASM_CONTAINER_REGISTRY)/$(WASM_CONTAINER_REGISTRY_REPOSITORY):$(WASM_CONTAINER_TAG)

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
LOCAL_BIN_PATH := ${PROJECT_PATH}/bin

KO_CONFIG_PATH ?= ${PROJECT_PATH}/etc/ko.yaml
KO_DOCKER_REPO ?= "quay.io/lburgazzoli/camel-go"
WASM_CONTAINER_IMAGE_REPO ?= quay.io/lburgazzoli/camel-go-wasm:latest

CGO_ENABLED := 0
BUILD_TAGS := -tags components_all -tags steps_all

LINT_GOGC := 10
LINT_DEADLINE := 10m

## Tools
GOIMPORT ?= $(LOCALBIN)/goimports
GOIMPORT_VERSION ?= latest
KO ?= $(LOCALBIN)/ko
KO_VERSION ?= main
TINYGO_VERSION ?= 0.31.2
GOLANGCI ?= $(LOCALBIN)/golangci-lint
GOLANGCI_VERSION ?= v1.55.2
CODEGEN_VERSION ?= v0.28.4
KUSTOMIZE_VERSION ?= v5.0.1
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_TOOLS_VERSION ?= v0.13.0
KIND_VERSION ?= v0.20.0
KIND ?= $(LOCALBIN)/kind
OPERATOR_SDK_VERSION ?= v1.32.0
OPERATOR_SDK ?= $(LOCALBIN)/operator-sdk
OPM_VERSION ?= v1.33.0
OPM ?= $(LOCALBIN)/opm
YQ ?= $(LOCALBIN)/yq
KUBECTL ?= kubectl
DAPR_VERSION ?= 1.12.0


# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif


ifndef ignore-not-found
  ignore-not-found = false
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

.PHONY: clean
clean:
	go clean -x
	go clean -x -testcache
	rm -f $(LOCAL_BIN_PATH)/camel

.PHONY: fmt
fmt: goimport
	$(GOIMPORT) -l -w .
	go fmt ./...

.PHONY: test
test:
	go test $(BUILD_TAGS) ./pkg/... ./test/...

.PHONY: deps
deps:
	go mod tidy

.PHONY: check/lint
check: check/lint

.PHONY: check/lint
check/lint: golangci-lint
	@$(GOLANGCI) run \
		--config .golangci.yml \
		--out-format tab \
		--skip-dirs etc \
		--deadline $(LINT_DEADLINE)

.PHONY: build
build: fmt
	CGO_ENABLED=0 go build -o $(LOCAL_BIN_PATH)/camel $(BUILD_TAGS) cmd/camel/main.go

.PHONY: image
image: ko
	KO_DOCKER_REPO=$(CONTAINER_REGISTRY)/$(CONTAINER_REGISTRY_REPOSITORY) \
	KO_CONFIG_PATH=$(PROJECT_PATH)/etc/ko.yaml \
	$(KO) build \
		--bare \
		--local \
		--tags $(CONTAINER_TAG) \
		--sbom none \
		./cmd/camel

.PHONY: image/publish
image/publish: ko
	KO_DOCKER_REPO=$(CONTAINER_REGISTRY)/$(CONTAINER_REGISTRY_REPOSITORY) \
	KO_CONFIG_PATH=$(PROJECT_PATH)/etc/ko.yaml \
	$(KO) build \
		--bare \
		--tags $(CONTAINER_TAG) \
		--sbom none \
		./cmd/camel

.PHONY: run/operator
run/operator: install
	go run -ldflags="$(GOLDFLAGS)" cmd/camel/main.go operator --leader-election=false --zap-devel

.PHONY: wasm/build
wasm/build:
	TINYGO_VERSION=$(TINYGO_VERSION) ./hack/scripts/build_wasm.sh $(PROJECT_PATH) etc/wasm/fn/simple_process.go etc/wasm/fn/simple_process.wasm
	TINYGO_VERSION=$(TINYGO_VERSION) ./hack/scripts/build_wasm.sh $(PROJECT_PATH) etc/wasm/fn/simple_logger.go etc/wasm/fn/simple_logger.wasm
	TINYGO_VERSION=$(TINYGO_VERSION) ./hack/scripts/build_wasm.sh $(PROJECT_PATH) etc/wasm/fn/to_upper.go etc/wasm/fn/to_upper.wasm
	TINYGO_VERSION=$(TINYGO_VERSION) ./hack/scripts/build_wasm.sh $(PROJECT_PATH) etc/wasm/fn/to_lower.go etc/wasm/fn/to_lower.wasm

.PHONY: wasm/build/ci
wasm/build/ci:
	./hack/scripts/build_wasm_ci.sh $(PROJECT_PATH) etc/wasm/fn/simple_process.go etc/wasm/fn/simple_process.wasm
	./hack/scripts/build_wasm_ci.sh $(PROJECT_PATH) etc/wasm/fn/simple_logger.go etc/wasm/fn/simple_logger.wasm
	./hack/scripts/build_wasm_ci.sh $(PROJECT_PATH) etc/wasm/fn/to_upper.go etc/wasm/fn/to_upper.wasm
	./hack/scripts/build_wasm_ci.sh $(PROJECT_PATH) etc/wasm/fn/to_lower.go etc/wasm/fn/to_lower.wasm

.PHONY: wasm/publish
wasm/publish:
	 oras push --verbose $(WASM_CONTAINER_IMAGE) \
 		etc/wasm/fn/simple_process.wasm:application/vnd.module.wasm.content.layer.v1+wasm \
 		etc/wasm/fn/simple_logger.wasm:application/vnd.module.wasm.content.layer.v1+wasm \
 		etc/wasm/fn/to_upper.wasm:application/vnd.module.wasm.content.layer.v1+wasm \
		etc/wasm/fn/to_lower.wasm:application/vnd.module.wasm.content.layer.v1+wasm

.PHONY: generate
generate: codegen-tools-install
	$(PROJECT_PATH)/hack/scripts/gen_res.sh $(PROJECT_PATH)
	$(PROJECT_PATH)/hack/scripts/gen_client.sh $(PROJECT_PATH)


.PHONY: manifests
manifests: codegen-tools-install
	$(PROJECT_PATH)/hack/scripts/gen_crd.sh $(PROJECT_PATH)


.PHONY: install
install: manifests kustomize
	$(KUSTOMIZE) build config/crd | $(KUBECTL) apply -f -

.PHONY: uninstall
uninstall: manifests kustomize
	$(KUSTOMIZE) build config/crd | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: kind/setup
kind/setup: kind
	$(KIND) create cluster \
		--config $(PROJECT_PATH)/etc/kind/kind-cluster-config.yaml \
		--image=kindest/node:v1.28.0 \
		--name "camel-go"

.PHONY: kind/teardown
kind/teardown: kind
	$(KIND) delete cluster  --name "camel-go"

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

.PHONY: goimport
goimport: $(GOIMPORT)
$(GOIMPORT): $(LOCALBIN)
	@test -s $(GOIMPORT) || \
	GOBIN=$(LOCALBIN) go install golang.org/x/tools/cmd/goimports@$(GOIMPORT_VERSION)

.PHONY: ko
ko: $(KO)
$(KO): $(LOCALBIN)
	@test -s $(KO) || \
	GOBIN=$(LOCALBIN) go install github.com/google/ko@$(KO_VERSION)

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI)
$(GOLANGCI): $(LOCALBIN)
	@test -s $(GOLANGCI) || \
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION)


.PHONY: kustomize
kustomize: $(KUSTOMIZE)
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(LOCALBIN)/kustomize || \
	GOBIN=$(LOCALBIN) GO111MODULE=on go install sigs.k8s.io/kustomize/kustomize/v5@$(KUSTOMIZE_VERSION)

.PHONY: yq
yq: $(YQ)
$(YQ): $(LOCALBIN)
	@test -s $(LOCALBIN)/yq || \
	GOBIN=$(LOCALBIN) go install github.com/mikefarah/yq/v4@latest


.PHONY: kind
kind: $(KIND)
$(KIND): $(LOCALBIN)
	@test -s $(LOCALBIN)/kind || \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/kind@$(KIND_VERSION)

.PHONY: codegen-tools-install
codegen-tools-install: $(LOCALBIN)
	@echo "Installing code gen tools"
	$(PROJECT_PATH)/hack/scripts/install_gen_tools.sh $(PROJECT_PATH) $(CODEGEN_VERSION) $(CONTROLLER_TOOLS_VERSION)

.PHONY: operator-sdk
operator-sdk: $(OPERATOR_SDK)
$(OPERATOR_SDK): $(LOCALBIN)
	@echo "Installing operator-sdk"
	$(PROJECT_PATH)/hack/scripts/install_operator_sdk.sh $(PROJECT_PATH) $(OPERATOR_SDK_VERSION)

.PHONY: opm
opm: $(OPM)
$(OPM): $(LOCALBIN)
	@echo "Installing opm"
	$(PROJECT_PATH)/hack/scripts/install_opm.sh $(PROJECT_PATH) $(OPM_VERSION)

