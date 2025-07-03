// package app contains the application business logic
package app

import (
	"context"
	"log/slog"
	"time"

	"k8s-controller/internal/app/handlers"
	"k8s-controller/internal/domain"
	"k8s-controller/internal/infrastructure/config"
	"k8s-controller/internal/infrastructure/kubernetes"
)

// KubernetesController is responsible for watching and reacting to Kubernetes resources
type KubernetesController struct {
	client          kubernetes.Client
	resourceService domain.ResourceService
	resourceHandler *handlers.ResourceHandler
	ctx             context.Context
	cancelFunc      context.CancelFunc
	config          *config.Config
}

// NewKubernetesController creates a new controller instance
func NewKubernetesController(cfg *config.Config) *KubernetesController {
	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Use default config if not provided
	if cfg == nil {
		cfg = config.Default()
	}

	// Create client
	client := kubernetes.NewClient()

	// Create domain services
	resourceService := domain.NewResourceService(client)

	// Create handlers
	resourceHandler := handlers.NewResourceHandler(resourceService)

	// Set handler in client
	client.SetEventHandler(resourceHandler)

	return &KubernetesController{
		client:          client,
		resourceService: resourceService,
		resourceHandler: resourceHandler,
		ctx:             ctx,
		cancelFunc:      cancel,
		config:          cfg,
	}
}

// Start initializes and starts the controller
func (c *KubernetesController) Start() error {
	slog.Info("Starting Kubernetes controller")

	// Connect to Kubernetes cluster
	if err := c.client.Connect(c.ctx); err != nil {
		slog.Error("Failed to connect to Kubernetes", "error", err)
		return err
	}

	// Start watching resources in a goroutine
	go func() {
		if err := c.resourceService.WatchResources(c.ctx); err != nil {
			if c.ctx.Err() == nil { // Only log if not due to context cancellation
				slog.Error("Error watching resources", "error", err)
			}
		}
	}()

	// Start a periodic health check
	go c.startPeriodicHealthCheck()

	return nil
}

// Stop gracefully stops the controller
func (c *KubernetesController) Stop() {
	slog.Info("Stopping Kubernetes controller")
	c.cancelFunc()
}

// startPeriodicHealthCheck runs a periodic health check
func (c *KubernetesController) startPeriodicHealthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			slog.Debug("Health check: Controller is running")
			// Add more health check logic here if needed
		case <-c.ctx.Done():
			return
		}
	}
}
