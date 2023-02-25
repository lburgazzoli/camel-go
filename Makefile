
# Image URL to use all building/pushing image targets
IMG ?= quay.io/lburgazzoli/camel-go:latest


MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
LOCAL_BIN_PATH := ${PROJECT_PATH}/bin
KO_CONFIG_PATH := ${PROJECT_PATH}/etc/ko.yaml
KO_DOCKER_REPO := "quay.io/lburgazzoli/cos-fleetshard"
CGO_ENABLED := 0
BUILD_TAGS := -tags components_all -tags steps_all

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

.PHONY: fmt
fmt: goimport
	$(LOCALBIN)/goimports -l -w .
	go fmt ./...

.PHONY: test
test:
	go test $(BUILD_TAGS) ./...

.PHONY: deps
deps:
	go mod tidy

.PHONY: lint
lint: golangci-lint
	$(LOCALBIN)/golangci-lint run --config .golangci.yml --out-format tab

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

.PHONY: examples
examples:
	docker run \
		--rm \
		--volume $(PROJECT_PATH):/src \
		tinygo/tinygo:0.27.0 \
			tinygo build \
			-o wasm.wasm \
			-target=wasm examples/wasm/export


##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

## Tool Binaries
GOIMPORT ?= $(LOCALBIN)/goimports
KO ?= $(LOCALBIN)/ko
GOLANGCILINT ?=  $(LOCALBIN)/golangci-lint
TINYGO ?=  $(LOCALBIN)/tinygo

.PHONY: goimport
goimport: $(GOIMPORT)
$(GOIMPORT): $(LOCALBIN)
	@test -s $(LOCALBIN)/goimport || \
	GOBIN=$(LOCALBIN) go install golang.org/x/tools/cmd/goimports@latest

.PHONY: golangci-lint
golangci-lint: $(GOLANGCILINT)
$(GOLANGCILINT): $(LOCALBIN)
	@test -s $(LOCALBIN)/golangci-lint || \
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.1

.PHONY: ko
ko: $(KO)
$(KO): $(LOCALBIN)
	@test -s $(LOCALBIN)/ko || GOBIN=$(LOCALBIN) go install github.com/google/ko@main
