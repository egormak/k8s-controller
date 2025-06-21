# k8s-controller

A CLI tool for Kubernetes operations.

## Overview

k8s-controller is a command-line interface tool built in Go, designed to help with Kubernetes operations. It uses the Cobra library to provide a robust command-line interface.

## Installation

### Prerequisites

- Go 1.22.3 or higher
- Access to a Kubernetes cluster (for actual operations)

### Build from source

```bash
git clone <repository-url>
cd k8s-controller
go build
```

## Usage

Basic usage:

```bash
./k8s-controller
```

For help and available commands:

```bash
./k8s-controller --help
```

## Project Structure

- `cmd/` - Contains command definitions using Cobra
- `main.go` - Application entry point

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework

## License

[Insert License Information]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
