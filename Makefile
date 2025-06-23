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

# Help
help:
	@echo "Available targets:"
	@echo "  all           - Clean, build, and test"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  run           - Run the application"
	@echo "  run-controller - Run the Kubernetes controller"
	@echo "  run-server    - Run the HTTP server"
	@echo "  run-debug     - Run with debug logs"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-push   - Push Docker image to registry"
	@echo "  help          - Show this help message"