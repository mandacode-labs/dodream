# Tools
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p "$(LOCALBIN)"

GOLANGCI_LINT_VERSION ?= v2.11.4
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)

# Use system golangci-lint if available and version matches, otherwise download
.PHONY: golangci-lint
golangci-lint: ## Download golangci-lint locally if necessary.
	@if command -v golangci-lint >/dev/null 2>&1; then \
		SYSTEM_VER=$$(golangci-lint version --format=short 2>/dev/null || echo "unknown"); \
		if [ "$$SYSTEM_VER" = "$(GOLANGCI_LINT_VERSION)" ]; then \
			echo "Using system golangci-lint $(GOLANGCI_LINT_VERSION)"; \
			exit 0; \
		fi; \
	fi; \
	if [ ! -f "$(GOLANGCI_LINT)" ]; then \
		echo "Downloading golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCALBIN) $(GOLANGCI_LINT_VERSION); \
		mv "$(LOCALBIN)/golangci-lint" "$(GOLANGCI_LINT)"; \
	fi

define GOLANGCI_LINT_CMD
$(shell if command -v golangci-lint >/dev/null 2>&1 && [ "$$(golangci-lint version --format=short 2>/dev/null)" = "$(GOLANGCI_LINT_VERSION)" ]; then echo golangci-lint; else echo $(GOLANGCI_LINT); fi)
endef

# Linting
.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter
	$(GOLANGCI_LINT_CMD) run

.PHONY: fmt
fmt: golangci-lint ## Run go fmt and fix lint issues
	$(GOLANGCI_LINT_CMD) fmt
	$(GOLANGCI_LINT_CMD) run --fix

# Code Generation
.PHONY: ent-gen
ent-gen: ## Generate ent code
	go generate ./ent

.PHONY: proto-gen
proto-gen: ## Generate protobuf code
	rm -rf pkg/proto/dodream
	mkdir -p pkg/proto/dodream/engine/v1
	protoc --proto_path=proto \
		--go_out=pkg/proto --go_opt=paths=source_relative \
		--go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative \
		dodream/engine/v1/engine.proto

.PHONY: ogen-gen
ogen-gen: ## Generate OpenAPI code with ogen
	rm -rf pkg/oas
	ogen --clean api/openapi/v1/openapi.yaml
	mkdir -p pkg/oas
	mv api/openapi/v1/oas_*.go pkg/oas/
	sed -i 's/package api/package oas/g' pkg/oas/*.go

.PHONY: mock-gen
mock-gen: ## Generate mocks with mockery
	~/go/bin/mockery

# Testing
.PHONY: test
test: ## Run unit tests only (excludes integration tests)
	go test -v -race -coverprofile=cover.out $$(go list ./... | grep -v /test)

.PHONY: test-unit
test-unit: ## Run unit tests only (alias for test)
	$(MAKE) test

.PHONY: test-integration
test-integration: ## Run integration tests with Docker
	go test -v -race -tags=integration ./test/integration/...

.PHONY: test-all
test-all: test test-integration ## Run all tests including integration

# Building
.PHONY: build
build: ## Build binary
	go build -o bin/dodream ./cmd/dodream

.PHONY: build-linux
build-linux: ## Build binary for Linux amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/dodream-linux-amd64 ./cmd/dodream

# Docker
.PHONY: docker-build
docker-build: ## Build Docker image locally
	docker build -f Dockerfile -t dodream:latest .

.PHONY: docker-run
docker-run: ## Run Docker container locally
	docker run --rm -p 8080:8080 dodream:latest

# All-in-one
.PHONY: all
all: fmt lint test build ## Format, lint, test and build

# Cleanup
.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f cover.out

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy

.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
