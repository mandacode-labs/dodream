# === Code Generation ===
.PHONY: gen ent-gen ogen-gen proto-gen mock-gen
gen: ent-gen ogen-gen mock-gen
ent-gen:   ; go generate ./ent
ogen-gen:  ; scripts/ogen-gen.sh
proto-gen: ; scripts/proto-gen.sh
mock-gen:  ; scripts/mock-gen.sh

# === Linting ===
.PHONY: lint fmt
lint: ; scripts/lint.sh
fmt:  ; scripts/lint.sh fmt ; scripts/lint.sh run --fix

# === Testing ===
.PHONY: test test-integration test-all
test:             ; go test -count=1 $(shell go list ./... | grep -v /test)
test-integration: ; scripts/test-integration.sh
test-all: test    ; $(MAKE) test-integration

# === Build ===
.PHONY: build build-linux
build:       ; go build -o bin/dodream ./cmd/dodream
build-linux: ; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/dodream-linux-amd64 ./cmd/dodream

# === Docker ===
.PHONY: docker-build docker-run
docker-build: ; docker build -f Dockerfile -t dodream:latest .
docker-run:   ; docker run --rm -p 8080:8080 dodream:latest

# === Hooks ===
.PHONY: hooks run-hooks
hooks:     ; scripts/hooks.sh install
run-hooks: ; scripts/hooks.sh run pre-commit && scripts/hooks.sh run pre-push

# === Housekeeping ===
.PHONY: clean tidy
clean: ; rm -rf bin/ cover.out
tidy:  ; go mod tidy

.DEFAULT_GOAL := help
help:  ; @grep -E '^[a-z0-9-]+:' $(MAKEFILE_LIST) | cut -d: -f1 | sort | column
