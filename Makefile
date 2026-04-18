# Copyright 2024 HAMi Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in" BASIS,
# WITHOUT WARRANTIESn# limitations under the License.

B= hami
VERshell git describe")
GIT_COMMITshell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go configuration
GO ?= go
GOFLAGS ?= -trimpath
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT) -X main.buildDate=$(BUILD_DATE)"

# Container image configuration
REGISTRY ?= ghcr.io/hami-io
IMAGE_NAME ?= $(REGISTRY)/$(BINARY_NAME)
IMAGE_TAG ?= $(VERSION)

# Directories
OUTPUT_DIR ?= bin
CMD_DIR ?= cmd

.PHONY: all build clean test lint fmt vet docker-build docker-push help

all: build

## build: Build all binaries
build:
	@echo "Building $(BINARY_NAME) version=$(VERSION)"
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(OUTPUT_DIR)/ ./$(CMD_DIR)/...

## test: Run unit tests
test:
	$(GO) test -v -race -coverprofile=coverage.out ./...

## test-coverage: Show test coverage report
test-coverage: test
	$(GO) tool cover -html=coverage.out

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## fmt: Format Go source files
fmt:
	$(GO) fmt ./...

## vet: Run go vet
vet:
	$(GO) vet ./...

## tidy: Tidy Go modules
tidy:
	$(GO) mod tidy

## docker-build: Build container image
docker-build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(IMAGE_NAME):$(IMAGE_TAG) \
		.

## docker-push: Push container image to registry
docker-push:
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

## clean: Remove build artifacts
clean:
	@rm -rf $(OUTPUT_DIR)
	@rm -f coverage.out

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
