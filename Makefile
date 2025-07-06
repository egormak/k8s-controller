# Makefile for k8s-controller
.PHONY: all build test clean run docker-build docker-push

# Variables
BINARY_NAME=k8s-controller
IMAGE_NAME=k8s-controller
IMAGE_TAG=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GO_FILES=$(shell find . -name "*.go" -type f)
GO_TEST_FILES=$(shell find . -name "*_test.go" -type f)
DOCKER_REGISTRY=

# Set default goal
.DEFAULT_GOAL := build

# Build the application
all: clean build test

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@go clean

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@go run main.go

# Run the controller
run-controller:
	@echo "Running controller..."
	@go run main.go control

# Run the HTTP server
run-server:
	@echo "Running HTTP server..."
	@go run main.go serve

# Run with debug logs
run-debug:
	@echo "Running with debug logs..."
	@go run main.go --log-level DEBUG

# Build docker image
docker-build:
	@echo "Building Docker image $(IMAGE_NAME):$(IMAGE_TAG)..."
	@docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

# Push docker image
docker-push:
	@echo "Pushing Docker image $(DOCKER_REGISTRY)$(IMAGE_NAME):$(IMAGE_TAG)..."
	@docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(DOCKER_REGISTRY)$(IMAGE_NAME):$(IMAGE_TAG)
	@docker push $(DOCKER_REGISTRY)$(IMAGE_NAME):$(IMAGE_TAG)

# Build release binaries for multiple platforms
build-release:
	@echo "Building release binaries..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(IMAGE_TAG)" -o dist/$(BINARY_NAME)-linux-amd64 .
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(IMAGE_TAG)" -o dist/$(BINARY_NAME)-linux-arm64 .
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(IMAGE_TAG)" -o dist/$(BINARY_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(IMAGE_TAG)" -o dist/$(BINARY_NAME)-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(IMAGE_TAG)" -o dist/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Release binaries built in dist/ directory"

# Create checksums for release binaries
checksums:
	@echo "Creating checksums..."
	@cd dist && sha256sum * > checksums.txt
	@echo "Checksums created in dist/checksums.txt"

# Create a local release (for testing)
release-local: clean build-release checksums
	@echo "Local release created in dist/ directory"

# Tag and push a new release
release-tag:
	@read -p "Enter version tag (e.g., v1.0.0): " tag; \
	git tag $$tag && \
	git push origin $$tag && \
	echo "Tag $$tag pushed. Release workflow will start automatically."

# Help
help:
	@echo "Available targets:"
	@echo "  all             - Clean, build, and test"
	@echo "  build           - Build the application"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  run             - Run the application"
	@echo "  run-controller  - Run the Kubernetes controller"
	@echo "  run-server      - Run the HTTP server"
	@echo "  run-debug       - Run with debug logs"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-push     - Push Docker image to registry"
	@echo "  build-release   - Build release binaries for multiple platforms"
	@echo "  checksums       - Create checksums for release binaries"
	@echo "  release-local   - Create a local release (for testing)"
	@echo "  release-tag     - Tag and push a new release"
	@echo "  help            - Show this help message"