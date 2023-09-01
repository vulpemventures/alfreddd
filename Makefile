.PHONY: build clean lint lint-fix run

## build: build the binary for your platform
build:
	@echo "Building binary..."
	@bash ./scripts/build

## clean: clean golang cache
clean:
	@echo "Cleaning..."
	@go clean

## lint: lint code
lint:
	@echo "Check linting..."
	@golangci-lint run

## lint-fix: lint code
lint-fix:
	@echo "Linting codebase..."
	@golangci-lint run --fix

## run: run in dev mode
run: clean
	@go run ./cmd/alfreddd

## test: run unit tests
test: clean
	@go test -v -count=1 -race ./...
