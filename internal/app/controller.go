package app

import (
	"context"
	"log/slog"

	"k8s-controller/internal/app/handlers"
	"k8s-controller/internal/domain"
	"k8s-controller/internal/infrastructure/config"
	"k8s-controller/internal/infrastructure/kubernetes"
)

// KubernetesController is the main application service that orchestrates the controller operations
type KubernetesController struct {
	client          kubernetes.Client
	resourceService domain.ResourceService
	resourceHandler *handlers.ResourceHandler
	config          *config.Config
}

// NewKubernetesController creates a new instance of the Kubernetes controller
func NewKubernetesController() *KubernetesController {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
	}

	// Initialize client with configuration
	k8sClient := kubernetes.NewClient()
	resourceService := domain.NewResourceService(k8sClient)
	resourceHandler := handlers.NewResourceHandler(resourceService)

	// Setup the event handler
	k8sClient.SetEventHandler(resourceHandler)

	return &KubernetesController{
		client:          k8sClient,
		resourceService: resourceService,
		resourceHandler: resourceHandler,
		config:          cfg,
	}
}

// Run starts the controller
func (c *KubernetesController) Run() error {
	slog.Info("Starting Kubernetes controller")

	ctx := context.Background()
	if err := c.client.Connect(ctx); err != nil {
		return err
	}

	// Start watching for resources
	return c.client.WatchResources(ctx)
}
