BINARY_NAME=calendar-go
GO=go
GOFLAGS=-v -trimpath
CGO_ENABLED=0

VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BRANCH=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
LDFLAGS=-ldflags "-s -w -extldflags '-static' -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Branch=$(BRANCH) -X main.BuildTime=$(BUILD_TIME)"
BUILDTAGS=-tags netgo,timetzdata

.PHONY: all build run test clean fmt vet lint deps tidy help

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOFLAGS) $(LDFLAGS) $(BUILDTAGS) -o $(BINARY_NAME) ./main.go

run:
	@echo "Running $(BINARY_NAME)..."
	$(GO) run ./main.go

test:
	@echo "Running tests..."
	$(GO) test -v -race -cover ./...

clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -f $(BINARY_NAME)

fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

vet:
	@echo "Vetting code..."
	$(GO) vet ./...

lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"; \
	fi

deps:
	@echo "Downloading dependencies..."
	$(GO) mod download

tidy:
	@echo "Tidying dependencies..."
	$(GO) mod tidy

help:
	@echo "Available targets:"
	@echo "  all     - Build the binary (default)"
	@echo "  build   - Build the binary"
	@echo "  run     - Run the application"
	@echo "  test    - Run tests with race detection and coverage"
	@echo "  clean   - Clean build artifacts"
	@echo "  fmt     - Format code"
	@echo "  vet     - Run go vet"
	@echo "  lint    - Run golangci-lint (if installed)"
	@echo "  deps    - Download dependencies"
	@echo "  tidy    - Tidy go.mod and go.sum"
	@echo "  help    - Show this help message"
