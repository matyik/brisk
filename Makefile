BINARY_NAME=brisk

.PHONY: all build test lint clean Run

all: lint test build

build:
	@echo "Building production binary..."
	go build -ldflags="-s -w" -o dist/$(BINARY_NAME) ./cmd/brisk/main.go

run:
	go run ./cmd/brisk/main.go $(file)

test:
	@echo "Running unit tests..."
	go test -v -race ./...

lint:
	@echo "Running linters..."
	golangci-lint run

clean:
	@echo "Cleaning build artifacts..."
	rm -rf dist/
	go clean