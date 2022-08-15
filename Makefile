# Adds $GOPATH/bit to $PATH
export PATH := $(PATH):$(shell go env GOPATH)/bin

help: ## Display this help message
	@echo "Usage:"
	@grep -E "^[a-zA-Z_-]+:.*? ## .+$$" $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-24s\033[0m %s\n", $$1, $$2}'

lint: ## Run golangci-lint
	golangci-lint run

lint-install: ## Install golangci-lint
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint-list: ## Run golangci-lint linters (print enabled & disabled linters)
	golangci-lint linters

test: ## Run tests
	mkdir -p bin
	go test -coverprofile bin/cover.out \
	$(shell go list ./... | grep -v internal)

cover: test ## Run tests & show coverage
	mkdir -p bin
	go tool cover -func bin/cover.out

race: ## Run tests with race flag
	go test -race -count=1 ./...

pre-commit: test lint ## Run pre commit checks

.PHONY: help lint lint-install lint-list test cover race pre-commit
