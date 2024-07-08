BINARY  := molen
PKG     := $(shell go list -m)
LOCAL_BIN_DIR := $(PWD)/bin

.DEFAULT_GOAL := help

.PHONY: run
run: ## Run the binary file
	go run cmd/molen/main.go

.PHONY: env
env: ## Create environment in foreground
	docker compose -p molen -f env/docker-compose.yml up -d

.PHONY: env-down
env-down: ## Destroy environment
	docker compose -p molen down

.PHONY: build
build: docs ## Build the binary file
	goreleaser build --snapshot --rm-dist --single-target

.PHONY: docs
docs: ## Generate swagger documentation
	swag init --pd -g internal/server/server.go

.PHONY: lint
lint: ## Lint Go files
	@golangci-lint --version
	@GOPATH="$(shell dirname $(PWD))" $(LOCAL_BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION) run ./...

.PHONY: test
test: ## Run unit tests
	@go test -v -race ./...

.PHONY: coverage
coverage: ## Run unit tests with coverage
	@go test -v -race -cover -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
