# Kubernetes Controller

A Go-based Kubernetes controller built using Domain-Driven Design (DDD) principles for better separation of concerns, maintainability, and extensibility.

## Project Structure

The project follows a DDD architecture with clear separation between layers:

```
k8s-controller/
├── cmd/              # Command-line entry points
│   └── root.go       # Root command implementation
├── internal/         # Internal packages (not importable from outside)
│   ├── app/          # Application services
│   │   ├── controller.go      # Main controller orchestration
│   │   └── handlers/          # Application event handlers
│   │       └── resource_handler.go
│   ├── domain/       # Domain model and services
│   │   ├── models.go          # Domain entities and value objects
│   │   └── resource_service.go # Domain service interfaces
│   └── infrastructure/ # Infrastructure implementations
│       ├── config/           # Configuration handling
│       │   └── config.go
│       └── kubernetes/       # Kubernetes client implementation
│           ├── client.go
│           └── informer.go
├── k8s-config.sample.yaml # Sample configuration file
├── main.go          # Application entry point
└── README.md        # Project documentation
```

## Domain-Driven Design (DDD) Layers

1. **Domain Layer** (`internal/domain/`)
   - Contains the core business logic and domain model
   - Defines interfaces that are implemented by infrastructure layer
   - Independent of external frameworks or technologies

2. **Application Layer** (`internal/app/`)
   - Orchestrates domain objects to perform specific use cases
   - Handles the flow of data and coordinates between different components
   - Contains application services and event handlers

3. **Infrastructure Layer** (`internal/infrastructure/`)
   - Implements domain interfaces using specific technologies (Kubernetes, config)
   - Handles technical concerns such as database connections, external APIs, etc.
   - Adapts external systems to the domain model

4. **User Interface Layer** (`cmd/`)
   - Handles user interactions via command line
   - Passes commands to the application layer for processing

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Access to a Kubernetes cluster
- `kubectl` configured with cluster access

### Installation

```bash
# Clone the repository
git clone https://github.com/egorka/k8s-controller.git
cd k8s-controller

# Build the project
go build -o k8s-controller

# Copy sample config
cp k8s-config.sample.yaml k8s-config.yaml

# Edit the configuration as needed
vim k8s-config.yaml
```

### Usage

```bash
# Run the controller with default configuration
./k8s-controller

# Run with specific config file
./k8s-controller --config /path/to/config.yaml

# Run with debug logging
./k8s-controller --log-level DEBUG
```

## Configuration

The controller can be configured using a YAML file. See `k8s-config.sample.yaml` for example configuration options.

Key configuration options:

- `log.level`: Logging level (DEBUG, INFO, WARN, ERROR)
- `kubernetes.kubeconfig`: Path to kubeconfig file (optional)
- `kubernetes.namespaces`: Comma-separated list of namespaces to watch
- `kubernetes.resources`: Comma-separated list of resource types to watch

## Extending the Controller

To add support for new resource types:

1. Update the `setupInformer` method in `internal/infrastructure/kubernetes/informer.go`
2. Add domain-specific handling logic in `internal/app/handlers/resource_handler.go`
3. Update the configuration to include the new resource type

## License

[MIT License](LICENSE)
