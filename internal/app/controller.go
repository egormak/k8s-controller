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
func NewKubernetesController() *KubernetesController {
	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		cfg = &config.Config{} // Use default config on error
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

// Run starts the controller's main loop
func (c *KubernetesController) Run() error {
	slog.Info("Initializing Kubernetes controller")

	// Connect to the Kubernetes cluster
	if err := c.client.Connect(c.ctx); err != nil {
		return err
	}

	slog.Info("Kubernetes controller started")

	// Start watching resources
	go func() {
		if err := c.resourceService.WatchResources(c.ctx); err != nil {
			if c.ctx.Err() == nil { // Only log if not due to context cancellation
				slog.Error("Error watching resources", "error", err)
			}
		}
	}()

	// Start a periodic health check
	go c.startPeriodicHealthCheck()

	// Block until context is cancelled
	<-c.ctx.Done()

	slog.Info("Kubernetes controller stopped")
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
			slog.Info("Controller health check", "status", "running")
		case <-c.ctx.Done():
			return
		}
	}
}
