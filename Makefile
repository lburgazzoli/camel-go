
# Image URL to use all building/pushing image targets
IMG ?= quay.io/lburgazzoli/camel-go:latest


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
TINYGO_VERSION ?= 0.29.0
GOLANGCI ?= $(LOCALBIN)/golangci-lint
GOLANGCI_VERSION ?= v1.52.2

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development


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

##@ Build

.PHONY: build
build: fmt
	CGO_ENABLED=0 go build -o $(LOCAL_BIN_PATH)/camel $(BUILD_TAGS) cmd/camel/main.go

.PHONY: image
image: ko
	KO_DOCKER_REPO=quay.io/lburgazzoli/camel-go \
	$(KO) build \
		--bare \
		--local \
		--tags latest \
		--platform=linux/amd64,linux/arm64 \
		./cmd/camel

.PHONY: image/publish
image/publish: ko
	KO_DOCKER_REPO=quay.io/lburgazzoli/camel-go \
	$(KO) build \
		--bare \
		--tags latest \
		--platform=linux/amd64,linux/arm64 \
		./cmd/camel


.PHONY: image/wasm
image/wasm:
	 oras push --verbose $(WASM_CONTAINER_IMAGE_REPO) \
 		etc/wasm/fn/simple_process.wasm:application/vnd.module.wasm.content.layer.v1+wasm \
 		etc/wasm/fn/simple_logger.wasm:application/vnd.module.wasm.content.layer.v1+wasm \
 		etc/wasm/fn/to_upper.wasm:application/vnd.module.wasm.content.layer.v1+wasm \
		etc/wasm/fn/to_lower.wasm:application/vnd.module.wasm.content.layer.v1+wasm


.PHONY: run/examples/dapr/flow
run/examples/dapr/flow:
	dapr run \
		 --app-id flow \
         --log-level debug \
         --app-protocol http \
	  	 --app-port 8080 \
		 --dapr-http-port 3500 \
         --resources-path ./etc/examples/dapr/config \
         -- go run cmd/camel/main.go run --dev --route ./etc/examples/dapr/dapr.yaml

.PHONY: run/examples/dapr/pub
run/examples/dapr/pub:
	dapr run \
		 --app-id pub \
         --log-level info \
         --resources-path ./etc/examples/dapr/config \
         -- go run cmd/camel/main.go dapr pub --pubsub-name sensors --topic iot source=sensor-1 data=foo

.PHONY: build/wasm
build/wasm:
	@docker run \
		--rm \
		-ti \
		-v $(PROJECT_PATH):/src:Z \
		-w /src \
		tinygo/tinygo:$(TINYGO_VERSION) \
		tinygo build \
			-target=wasi \
			-scheduler=none \
			-o etc/wasm/fn/simple_process.wasm  \
			etc/wasm/fn/simple_process.go

	@docker run \
		--rm \
		-ti \
		-v $(PROJECT_PATH):/src:Z \
		-w /src \
		tinygo/tinygo:$(TINYGO_VERSION) \
		tinygo build \
			-target=wasi \
			-scheduler=none \
			-o etc/wasm/fn/simple_logger.wasm  \
			etc/wasm/fn/simple_logger.go

	@docker run \
		--rm \
		-ti \
		-v $(PROJECT_PATH):/src:Z \
		-w /src \
		tinygo/tinygo:$(TINYGO_VERSION) \
		tinygo build \
			-target=wasi \
			-scheduler=none \
			-o etc/wasm/fn/to_upper.wasm  \
			etc/wasm/fn/to_upper.go

	@docker run \
		--rm \
		-ti \
		-v $(PROJECT_PATH):/src:Z \
		-w /src \
		tinygo/tinygo:$(TINYGO_VERSION) \
		tinygo build \
			-target=wasi \
			-scheduler=none \
			-o etc/wasm/fn/to_lower.wasm  \
			etc/wasm/fn/to_lower.go

##@ Build Dependencies

## Location to install dependencies to
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

