BINARY  := molen
PKG     := $(shell go list -m)
LOCAL_BIN_DIR := $(PWD)/bin

## swaggo configuration
SWAG_VERSION := $(shell grep github.com/swaggo/swag go.mod | xargs echo | cut -d" " -f2)

## golangci configuration
# GOLANGCI_CONFIG_URL   := https://raw.githubusercontent.com/worldline-go/guide/main/lint/.golangci.yml
GOLANGCI_LINT_VERSION := v1.54.2

.DEFAULT_GOAL := help

.PHONY: run env build golangci docs lint test coverage help

run: ## Run the binary file
	go run cmd/molen/main.go

env: ## Run environment in foreground
	docker compose -p molen -f env/docker-compose.yml up

build: docs ## Build the binary file
	goreleaser build --snapshot --rm-dist --single-target

bin/swag-$(SWAG_VERSION):
	@echo "> downloading swag@$(SWAG_VERSION)"
	@GOBIN=$(LOCAL_BIN_DIR) go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)
	@mv $(LOCAL_BIN_DIR)/swag $(LOCAL_BIN_DIR)/swag-$(SWAG_VERSION)

bin/golangci-lint-$(GOLANGCI_LINT_VERSION):
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCAL_BIN_DIR) $(GOLANGCI_LINT_VERSION)
	@mv $(LOCAL_BIN_DIR)/golangci-lint $(LOCAL_BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION)

docs: bin/swag-$(SWAG_VERSION)
	@$(LOCAL_BIN_DIR)/swag-$(SWAG_VERSION) init --pd -g internal/server/server.go

lint: bin/golangci-lint-$(GOLANGCI_LINT_VERSION) ## Lint Go files
	@$(LOCAL_BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION) --version
	@GOPATH="$(shell dirname $(PWD))" $(LOCAL_BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION) run ./...

test: ## Run unit tests
	@go test -v -race ./...

coverage: ## Run unit tests with coverage
	@go test -v -race -cover -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
