# Kubernetes Controller

A lightweight, Go-based Kubernetes controller with HTTP API capabilities. This project demonstrates how to build a clean, modular Kubernetes controller with a well-defined separation of concerns.

## Features

- Kubernetes controller that watches and responds to resource events
- HTTP API for interacting with Kubernetes resources
- Clean separation between business logic and infrastructure
- Configuration through environment variables or config file
- Support for multiple namespaces and resource types

## Project Structure

The project follows a simplified layered architecture:

```
k8s-controller/
├── cmd/              # Command-line entry points
│   ├── control.go    # Kubernetes controller command
│   ├── list.go       # List resources command
│   ├── root.go       # Root command implementation
│   └── serve.go      # HTTP server command
├── internal/         # Internal packages (not importable from outside)
│   ├── app/          # Application services
│   │   ├── controller.go      # Main controller orchestration
│   │   └── handlers/          # Event handlers
│   │       └── resource_handler.go
│   ├── domain/       # Domain model and services
│   │   ├── deployment.go      # Deployment model
│   │   ├── models.go          # Core model entities
│   │   └── resource_service.go # Resource service
│   └── infrastructure/ # Infrastructure implementations
│       ├── config/           # Configuration handling
│       │   └── config.go
│       ├── kubernetes/       # Kubernetes client implementation
│       │   ├── client.go
│       │   └── informer.go
│       └── server/          # HTTP server implementation
│           ├── deployment_controller.go
│           └── server.go
├── manifests/        # Kubernetes manifests for testing
│   └── nginx_deployment.yaml
├── k8s-config.sample.yaml # Sample configuration file
├── Dockerfile       # Container build definition
├── Makefile        # Build automation
└── main.go         # Application entry point
```

## Getting Started

### Prerequisites

- Go 1.24+
- Access to a Kubernetes cluster (local or remote)
- kubectl configured with access to your cluster

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/k8s-controller.git
   cd k8s-controller
   ```

2. Build the application:
   ```bash
   make build
   ```

3. Create a configuration file:
   ```bash
   cp k8s-config.sample.yaml k8s-config.yaml
   # Edit k8s-config.yaml with your settings
   ```

### Usage

#### Starting the HTTP Server

```bash
./k8s-controller serve --port 8080
```

#### Starting the Kubernetes Controller

```bash
./k8s-controller control --namespaces default,kube-system
```

#### Listing Deployments

```bash
./k8s-controller list deployments --namespace default
```

## Configuration

The application can be configured using:

1. Command-line flags
2. Environment variables 
3. Configuration file (YAML)

Example configuration file:

```yaml
log:
  level: INFO
kubernetes:
  namespaces: default,kube-system
  resources: deployments,services,pods
server:
  port: 8080
```

## Development

### Running Tests

```bash
make test
```

### Building Docker Image

```bash
make docker-build
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
