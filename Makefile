# Tools
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p "$(LOCALBIN)"

GOLANGCI_LINT = $(LOCALBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v2.12.2

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
define go-install-tool
@[ -f "$(1)-$(3)" ] && [ "$$(readlink "$(1)" 2>/dev/null)" = "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f "$(1)" ;\
GOBIN="$(LOCALBIN)" go install $${package} ;\
mv "$(LOCALBIN)/$$(basename "$(1)")" "$(1)-$(3)" ;\
} ;\
ln -sf "$$(realpath "$(1)-$(3)")" "$(1)"
endef

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

# Linting
.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter
	"$(GOLANGCI_LINT)" run

.PHONY: fmt
fmt: golangci-lint ## Run go fmt and fix lint issues
	"$(GOLANGCI_LINT)" fmt
	"$(GOLANGCI_LINT)" run --fix

# Code Generation
.PHONY: ent-gen
ent-gen: ## Generate ent code
	go generate ./ent

# Testing
.PHONY: test
test: ## Run all tests
	go test -v -race -coverprofile=cover.out ./...

.PHONY: test-unit
test-unit: ## Run unit tests only
	go test -v -race $$(go list ./... | grep -v /test) --coverprofile=cover.out

# Building
.PHONY: build
build: ## Build binary
	go build -o bin/dodream-engine ./cmd/dodream-engine

.PHONY: build-linux
build-linux: ## Build binary for Linux amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/dodream-engine-linux-amd64 ./cmd/dodream-engine

# Docker
.PHONY: docker-build
docker-build: ## Build Docker image locally
	docker build -f build/docker/Dockerfile -t dodream-engine:latest .

.PHONY: docker-run
docker-run: ## Run Docker container locally
	docker run --rm -p 8080:8080 dodream-engine:latest

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
