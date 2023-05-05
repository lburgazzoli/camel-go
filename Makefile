
# Image URL to use all building/pushing image targets
IMG ?= quay.io/lburgazzoli/camel-go:latest


MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
LOCAL_BIN_PATH := ${PROJECT_PATH}/bin
KO_CONFIG_PATH := ${PROJECT_PATH}/etc/ko.yaml
KO_DOCKER_REPO := "quay.io/lburgazzoli/gamel"
CGO_ENABLED := 0
BUILD_TAGS := -tags components_all -tags steps_all
LINT_GOGC := 10
LINT_DEADLINE := 10m

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
	$(LOCALBIN)/goimports -l -w .
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
check/lint:
	@docker run \
		--rm \
		-t \
		-v $(PROJECT_PATH):/app:Z \
		-e GOGC=$(LINT_GOGC) \
		-w /app \
		golangci/golangci-lint:v1.52 golangci-lint run \
			--config .golangci.yml \
			--out-format tab \
			--skip-dirs etc \
			--skip-dirs pkg/wasm/interop \
			--deadline $(LINT_DEADLINE)

##@ Build

.PHONY: build
build: fmt
	CGO_ENABLED=0 go build -o $(LOCAL_BIN_PATH)/camel $(BUILD_TAGS) cmd/camel/main.go

.PHONY: image/publish
image/publish: ko
	KO_CONFIG_PATH=$(KO_CONFIG_PATH) KO_DOCKER_REPO=${KO_DOCKER_REPO} $(KO) build --sbom none --bare ./cmd/camel

.PHONY: image/local
image/local: ko
	KO_CONFIG_PATH=$(KO_CONFIG_PATH) KO_DOCKER_REPO=ko.local $(KO) build --sbom none --bare ./cmd/camel

.PHONY: image/kind
image/kind: ko
	KO_CONFIG_PATH=$(KO_CONFIG_PATH) KO_DOCKER_REPO=kind.local $(KO) build --sbom none --bare ./cmd/camel

.PHONY: image/wasm
image/wasm:
	 oras push --verbose docker.io/lburgazzoli/camel-go:latest \
 		etc/wasm/fn/simple_process.wasm:application/vnd.module.wasm.content.layer.v1+wasm \
 		etc/wasm/fn/simple_logger.wasm:application/vnd.module.wasm.content.layer.v1+wasm \
 		etc/wasm/fn/to_upper.wasm:application/vnd.module.wasm.content.layer.v1+wasm

.PHONY: build/wasm
build/wasm:
	@docker run \
		--rm \
		-ti \
		-v $(PROJECT_PATH):/src:Z \
		-w /src \
		tinygo/tinygo:0.27.0 \
		tinygo build \
			-target=wasi \
			-scheduler=none \
			-gc=leaking \
			-o etc/wasm/fn/simple_process.wasm  \
			etc/wasm/fn/simple_process.go

	@docker run \
		--rm \
		-ti \
		-v $(PROJECT_PATH):/src:Z \
		-w /src \
		tinygo/tinygo:0.27.0 \
		tinygo build \
			-target=wasi \
			-scheduler=none \
			-gc=leaking \
			-o etc/wasm/fn/simple_logger.wasm  \
			etc/wasm/fn/simple_logger.go

	@docker run \
		--rm \
		-ti \
		-v $(PROJECT_PATH):/src:Z \
		-w /src \
		tinygo/tinygo:0.27.0 \
		tinygo build \
			-target=wasi \
			-scheduler=none \
			-gc=leaking \
			-o etc/wasm/fn/to_upper.wasm  \
			etc/wasm/fn/to_upper.go

.PHONY: generate
generate: protoc-gen-go-plugin
	#go run karmem.org/cmd/karmem build --golang -o "pkg/wasm/interop" etc/message.km
	protoc \
		--plugin=$(LOCALBIN)/protoc-gen-go-plugin \
		--proto_path=$(PROJECT_PATH)/etc/proto \
		--go-plugin_out=$(PROJECT_PATH)/pkg \
		$(PROJECT_PATH)/etc/proto/processor.proto


##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

## Tool Binaries
GOIMPORT ?= $(LOCALBIN)/goimports
KO ?= $(LOCALBIN)/ko
PROTOC_PLUGIN ?= $(LOCALBIN)/protoc-gen-go-plugin

.PHONY: goimport
goimport: $(GOIMPORT)
$(GOIMPORT): $(LOCALBIN)
	@test -s $(LOCALBIN)/goimport || \
	GOBIN=$(LOCALBIN) go install golang.org/x/tools/cmd/goimports@latest

.PHONY: ko
ko: $(KO)
$(KO): $(LOCALBIN)
	@test -s $(LOCALBIN)/ko || GOBIN=$(LOCALBIN) go install github.com/google/ko@main

.PHONY: protoc-gen-go-plugin
protoc-gen-go-plugin: $(PROTOC_PLUGIN)
$(PROTOC_PLUGIN): $(LOCALBIN)
	@test -s $(LOCALBIN)/protoc-gen-go-plugin || GOBIN=$(LOCALBIN) go install github.com/knqyf263/go-plugin/cmd/protoc-gen-go-plugin@main
