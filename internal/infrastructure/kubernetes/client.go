package kubernetes

import (
	"context"
	"log/slog"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"k8s-controller/internal/domain"
	"k8s-controller/internal/infrastructure/config"
)

// Client implements the domain.ResourceClient interface
type Client interface {
	domain.ResourceClient
	SetEventHandler(handler ResourceEventHandler)
}

// kubeClient is a concrete implementation of the Client interface
type kubeClient struct {
	clientset    *kubernetes.Clientset
	config       *config.Config
	eventHandler ResourceEventHandler
}

// NewClient creates a new Kubernetes client
func NewClient() Client {
	return &kubeClient{}
}

// SetEventHandler sets the handler for resource events
func (c *kubeClient) SetEventHandler(handler ResourceEventHandler) {
	c.eventHandler = handler
}

// Connect establishes a connection to the Kubernetes cluster
func (c *kubeClient) Connect(ctx context.Context) error {
	slog.Info("Connecting to Kubernetes cluster")

	// Get kubeconfig from the config or use default location
	var kubeconfigPath string
	if c.config != nil && c.config.KubeconfigPath != "" {
		kubeconfigPath = c.config.KubeconfigPath
	} else {
		home := homedir.HomeDir()
		kubeconfigPath = filepath.Join(home, ".kube", "config")
	}

	// Build the config from the kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		slog.Error("Failed to build config from flags", "error", err)
		return err
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		slog.Error("Failed to create Kubernetes client", "error", err)
		return err
	}

	c.clientset = clientset
	slog.Info("Successfully connected to Kubernetes cluster")
	return nil
}

// WatchResources starts watching for resource events
func (c *kubeClient) WatchResources(ctx context.Context) error {
	slog.Info("Starting to watch resources")

	if c.eventHandler == nil {
		slog.Warn("No event handler set, resource events will not be processed")
		return nil
	}

	// Determine namespaces and resources to watch
	namespaces := []string{"default"}
	resources := []string{"deployments", "services", "pods"}

	if c.config != nil {
		if len(c.config.ResourceNamespaces) > 0 {
			namespaces = c.config.ResourceNamespaces
		}

		if len(c.config.WatchedResources) > 0 {
			resources = c.config.WatchedResources
		}
	}

	// Start the informers
	if err := c.startInformers(ctx, namespaces, resources, c.eventHandler); err != nil {
		return err
	}

	// Block until context is cancelled
	<-ctx.Done()
	return ctx.Err()
}

// GetResource retrieves a specific resource
func (c *kubeClient) GetResource(ctx context.Context, kind, name, namespace string) (domain.Resource, error) {
	slog.Debug("Getting resource", "kind", kind, "name", name, "namespace", namespace)

	// Implementation would use the clientset to get the resource
	// This is a placeholder
	return domain.Resource{}, nil
}

// ApplyResource creates or updates a resource
func (c *kubeClient) ApplyResource(ctx context.Context, resource domain.Resource) error {
	slog.Debug("Applying resource", "kind", resource.Kind, "name", resource.Name, "namespace", resource.Namespace)

	// Implementation would use the clientset to apply the resource
	// This is a placeholder
	return nil
}
