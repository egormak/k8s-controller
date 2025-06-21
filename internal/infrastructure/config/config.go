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
}

// Load retrieves the configuration from viper
func Load() (*Config, error) {
	cfg := &Config{
		LogLevel:           viper.GetString("log.level"),
		KubeconfigPath:     viper.GetString("kubernetes.kubeconfig"),
		ResourceNamespaces: getStringSlice("kubernetes.namespaces"),
		WatchedResources:   getStringSlice("kubernetes.resources"),
	}

	// Set defaults
	if cfg.LogLevel == "" {
		cfg.LogLevel = "INFO"
	}

	if len(cfg.ResourceNamespaces) == 0 {
		cfg.ResourceNamespaces = []string{"default"}
	}

	if len(cfg.WatchedResources) == 0 {
		cfg.WatchedResources = []string{"deployments", "services"}
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

	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
