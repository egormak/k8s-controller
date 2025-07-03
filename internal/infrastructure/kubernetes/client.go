package kubernetes

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"k8s-controller/internal/domain"
)

// Using the ResourceEventHandler interface defined in informer.go

// Client implements the domain.ResourceClient interface
type Client interface {
	domain.ResourceClient
	SetEventHandler(handler ResourceEventHandler)
	ListDeployments(ctx context.Context, namespace string) ([]domain.Deployment, error)
	GetDeploymentInformer(namespace string) (cache.SharedIndexInformer, error)
	InitializeInformers(ctx context.Context, namespaces []string) error
	SetNamespaces(namespaces []string)
	SetWatchedResources(resources []string)
}

// kubeClient is a concrete implementation of the Client interface
type kubeClient struct {
	clientset         *kubernetes.Clientset
	eventHandler      ResourceEventHandler
	informerFactories map[string]informers.SharedInformerFactory
	namespaces        []string
	watchedResources  []string
}

// NewClient creates a new Kubernetes client with sensible defaults
func NewClient() Client {
	return &kubeClient{
		informerFactories: make(map[string]informers.SharedInformerFactory),
		namespaces:        []string{"default"},
		watchedResources:  []string{"deployments", "services", "pods"},
	}
}

// SetNamespaces sets the namespaces to watch
func (c *kubeClient) SetNamespaces(namespaces []string) {
	if len(namespaces) > 0 {
		c.namespaces = namespaces
	}
}

// SetWatchedResources sets the types of resources to watch
func (c *kubeClient) SetWatchedResources(resources []string) {
	if len(resources) > 0 {
		c.watchedResources = resources
	}
}

// SetEventHandler sets the handler for resource events
func (c *kubeClient) SetEventHandler(handler ResourceEventHandler) {
	c.eventHandler = handler
}

// Connect establishes a connection to the Kubernetes cluster
func (c *kubeClient) Connect(ctx context.Context) error {
	slog.Info("Connecting to Kubernetes cluster")

	// Use default kubeconfig location
	home := homedir.HomeDir()
	kubeconfigPath := filepath.Join(home, ".kube", "config")

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

	// First initialize informers to ensure cache is ready
	if err := c.InitializeInformers(ctx, c.namespaces); err != nil {
		return err
	}

	// Then start watching resources with event handlers
	return c.startInformers(ctx, c.namespaces, c.watchedResources, c.eventHandler)
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

// ListDeployments retrieves all deployments in the specified namespace using the informer cache
func (c *kubeClient) ListDeployments(ctx context.Context, namespace string) ([]domain.Deployment, error) {
	slog.Debug("Listing deployments from cache", "namespace", namespace)

	if c.clientset == nil {
		return nil, fmt.Errorf("kubernetes client not connected")
	}

	// Check if we have an informer for this namespace
	factory, ok := c.informerFactories[namespace]
	if !ok {
		slog.Warn("No informer factory for namespace, falling back to direct API call", "namespace", namespace)
		// Fall back to direct API call if no informer is available
		deploymentList, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			slog.Error("Failed to list deployments", "error", err, "namespace", namespace)
			return nil, err
		}

		var deployments []domain.Deployment
		for _, dep := range deploymentList.Items {
			deployment := domain.Deployment{
				Name:              dep.Name,
				Namespace:         dep.Namespace,
				ReadyReplicas:     dep.Status.ReadyReplicas,
				UpdatedReplicas:   dep.Status.UpdatedReplicas,
				AvailableReplicas: dep.Status.AvailableReplicas,
				Replicas:          *dep.Spec.Replicas,
				Labels:            dep.Labels,
				CreationTimestamp: dep.CreationTimestamp.Format("2006-01-02 15:04:05"),
			}
			deployments = append(deployments, deployment)
		}
		return deployments, nil
	}

	// Get the store from the informer
	lister := factory.Apps().V1().Deployments().Lister()
	deploymentList, err := lister.Deployments(namespace).List(labels.Everything())
	if err != nil {
		slog.Error("Failed to list deployments from cache", "error", err, "namespace", namespace)
		return nil, err
	}

	var deployments []domain.Deployment
	for _, dep := range deploymentList {
		deployment := domain.Deployment{
			Name:              dep.Name,
			Namespace:         dep.Namespace,
			ReadyReplicas:     dep.Status.ReadyReplicas,
			UpdatedReplicas:   dep.Status.UpdatedReplicas,
			AvailableReplicas: dep.Status.AvailableReplicas,
			Replicas:          *dep.Spec.Replicas,
			Labels:            dep.Labels,
			CreationTimestamp: dep.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		deployments = append(deployments, deployment)
	}

	slog.Info("Successfully listed deployments", "count", len(deployments), "namespace", namespace)
	return deployments, nil
}

// InitializeInformers initializes informer factories for specified namespaces
func (c *kubeClient) InitializeInformers(ctx context.Context, namespaces []string) error {
	if c.clientset == nil {
		return fmt.Errorf("kubernetes client not connected")
	}

	slog.Info("Initializing informer factories", "namespaces", namespaces)

	// If no namespaces provided, use default
	if len(namespaces) == 0 {
		namespaces = []string{"default"}
	}

	// Create a factory for each namespace with 30s resync period
	for _, namespace := range namespaces {
		factory := informers.NewSharedInformerFactoryWithOptions(
			c.clientset,
			30*time.Second, // resync period
			informers.WithNamespace(namespace),
		)

		// Pre-create some commonly used informers to ensure they are available
		// This doesn't start watching yet, just creates the informers
		factory.Apps().V1().Deployments().Informer()

		// Store the factory
		c.informerFactories[namespace] = factory

		// Start the informer factory with a background context
		// This ensures we don't block even if the parent context is cancelled
		factory.Start(ctx.Done())

		slog.Info("Started informer factory", "namespace", namespace)
	}

	// Wait for the initial sync to complete with a reasonable timeout
	syncCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for namespace, factory := range c.informerFactories {
		deployInformer := factory.Apps().V1().Deployments().Informer()
		if !cache.WaitForCacheSync(syncCtx.Done(), deployInformer.HasSynced) {
			slog.Warn("Timeout waiting for deployment cache to sync", "namespace", namespace)
			// Continue despite timeout - the cache will eventually sync
		} else {
			slog.Info("Deployment cache synced", "namespace", namespace)
		}
	}

	return nil
}

// GetDeploymentInformer returns the deployment informer for the given namespace
func (c *kubeClient) GetDeploymentInformer(namespace string) (cache.SharedIndexInformer, error) {
	factory, ok := c.informerFactories[namespace]
	if !ok {
		return nil, fmt.Errorf("no informer factory for namespace %s", namespace)
	}

	// Return the deployment informer
	return factory.Apps().V1().Deployments().Informer(), nil
}
