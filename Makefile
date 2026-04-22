.PHONY: build test test-cover lint lint-fix fmt fmt-diff clean tools mod

BINARY_NAME=dendrite
BUILD_DIR=bin
GOLANGCI_LINT=go tool -modfile tools/go.mod golangci-lint

# Build
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/dendrite

# Test
test:
	go test -v -race -shuffle=on ./...

test-cover:
	go test -v -race -coverprofile=coverage.out -coverpkg=./... ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint
lint:
	$(GOLANGCI_LINT) run ./...

lint-fix:
	$(GOLANGCI_LINT) run --fix ./...

fmt:
	$(GOLANGCI_LINT) fmt ./...

fmt-diff:
	$(GOLANGCI_LINT) fmt ./... --diff

# Clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Tools
tools:
	cd tools && go mod tidy

# Go mod
mod:
	go mod tidy
