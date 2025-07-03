package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	LogLevel           string
	KubeconfigPath     string
	ResourceNamespaces []string
	WatchedResources   []string
	ServerPort         int
}

// Default returns a configuration with default values
func Default() *Config {
	return &Config{
		LogLevel:           "INFO",
		ResourceNamespaces: []string{"default"},
		WatchedResources:   []string{"deployments", "services"},
		ServerPort:         8080,
	}
}

// Load retrieves the configuration from viper
func Load() (*Config, error) {
	// Create config with default values
	cfg := Default()

	// Override with values from viper if present
	if viper.IsSet("log.level") {
		cfg.LogLevel = viper.GetString("log.level")
	}

	if viper.IsSet("kubernetes.kubeconfig") {
		cfg.KubeconfigPath = viper.GetString("kubernetes.kubeconfig")
	}

	if viper.IsSet("kubernetes.namespaces") {
		cfg.ResourceNamespaces = getStringSlice("kubernetes.namespaces")
	}

	if viper.IsSet("kubernetes.resources") {
		cfg.WatchedResources = getStringSlice("kubernetes.resources")
	}

	if viper.IsSet("server.port") {
		cfg.ServerPort = viper.GetInt("server.port")
	}

	return cfg, nil
}

// getStringSlice safely gets a string slice from viper
func getStringSlice(key string) []string {
	val := viper.GetString(key)
	if val == "" {
		return []string{}
	}

	items := strings.Split(val, ",")
	result := make([]string, 0, len(items))

	// Trim whitespace from each item
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
