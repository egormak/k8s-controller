# k8s-controller

A CLI tool for Kubernetes operations.

## Overview

k8s-controller is a command-line interface tool built in Go, designed to help with Kubernetes operations. It uses the Cobra library to provide a robust command-line interface and Viper for configuration management.

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

### Configuration

The application supports multiple configuration methods in the following priority order:

1. Command-line flags
2. Environment variables
3. Configuration file
4. Default values

#### Configuration File

By default, the application looks for a `k8s-config.yaml` file in the current directory. You can specify a custom configuration file using the `--config` flag:

```bash
./k8s-controller --config=/path/to/custom-config.yaml
```

Example configuration file format:

```yaml
server:
  port: 8080
```

#### Environment Variables

You can configure the application using environment variables:

1. Config file location:
   ```bash
   export CONFIG=/path/to/custom-config.yaml
   ./k8s-controller
   ```

2. Application settings (when used with `viper.AutomaticEnv()`):
   ```bash
   export SERVER_PORT=9090
   ./k8s-controller
   ```

## Project Structure

- `cmd/` - Contains command definitions using Cobra
- `main.go` - Application entry point

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management

## Advanced Configuration

### Setting Current Directory as Default

The application is configured to search for configuration files in the current working directory rather than the home directory.

### Using Viper's AutomaticEnv

The application uses `viper.AutomaticEnv()` to automatically bind environment variables to configuration keys. For example, if you have a configuration key `server.port`, you can set it using the environment variable `SERVER_PORT`.

## License

[Insert License Information]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
