package main

import (
	"fmt"
	"strings"
)

func makeMakefile(app string, noApiSpec bool) []byte {
	content := fmt.Sprintf(`
.PHONY: build clean cov help intergrationtest lint run test vet proto proto-lint

## build: build for all platforms
build:
	@echo "Building %s binary..."
	@bash ./scripts/build

## clean: cleans the binary
clean:
	@echo "Cleaning..."
	@go clean

## cov: generates coverage report
cov:
	@echo "Coverage..."
	@go test -cover ./...

## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## intergrationtest: runs integration tests
integrationtest:
	@echo "Running integration tests..."
	@go test -v -count=1 -race ./... $(go list ./... | grep internal/test)

## lint: lint codebase
lint:
	@echo "Linting code..."
	@golangci-lint run --fix

## run: run in dev mode
run: clean
	@echo "Running %s in dev mode..."
	@go run ./cmd/%s

## test: runs unit and component tests
test:
	@echo "Running unit tests..."
	@go test -v -count=1 -race ./... $(go list ./... | grep -v internal/test)

## vet: code analysis
vet:
	@echo "Running code analysis..."
	@go vet ./...
	`, app, app, app)

	if !noApiSpec {
		content += `
	
## proto: compile proto stubs
proto: proto-lint
	@echo "Compiling stubs..."
	@buf generate

## proto-lint: lint protos
proto-lint:
	@echo "Linting protos..."
	@buf lint
		`
	}

	return []byte(strings.Trim(content, "\n"))
}
