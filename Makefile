LOCALBIN ?= $(shell pwd)/bin

# Linting
.PHONY: lint
lint: ## Run golangci-lint linter
	scripts/lint.sh $(LOCALBIN)

.PHONY: fmt
fmt: ## Run go fmt and fix lint issues
	scripts/lint.sh $(LOCALBIN) fmt
	scripts/lint.sh $(LOCALBIN) run --fix

# Code Generation
.PHONY: ent-gen
ent-gen: ## Generate ent code
	go generate ./ent

.PHONY: proto-gen
proto-gen: ## Generate protobuf code with buf
	scripts/install-proto-tools.sh
	rm -rf pkg/proto/dodream
	PATH="$(shell go env GOPATH)/bin:$(PATH)" buf generate

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
test-unit: test ## Run unit tests only (alias for test)

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

# Hooks
.PHONY: install-hooks
install-hooks: ## Install lefthook git hooks
	~/go/bin/lefthook install

.PHONY: run-hooks
run-hooks: ## Run all lefthook hooks for testing
	~/go/bin/lefthook run pre-commit && ~/go/bin/lefthook run pre-push

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

.DEFAULT_GOAL := help
