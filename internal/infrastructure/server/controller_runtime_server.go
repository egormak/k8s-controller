package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s-controller/internal/domain"
	"k8s-controller/internal/infrastructure/config"
	"k8s-controller/internal/infrastructure/controller"
)

// ControllerRuntimeServer extends the basic server with controller-runtime functionality
type ControllerRuntimeServer struct {
	*Server
	controllerRuntime *controller.ControllerRuntime
	resourceService   domain.ResourceService
}

// NewControllerRuntimeServer creates a new server with controller-runtime capabilities
func NewControllerRuntimeServer(port int, cfg *config.Config) (*ControllerRuntimeServer, error) {
	// Create base server
	baseServer := NewServer(port)

	// Create controller runtime
	controllerRuntime, err := controller.NewControllerRuntime(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create controller runtime: %w", err)
	}

	// Create resource service using existing client
	resourceService := domain.NewResourceService(baseServer.kubeClient)

	server := &ControllerRuntimeServer{
		Server:            baseServer,
		controllerRuntime: controllerRuntime,
		resourceService:   resourceService,
	}

	return server, nil
}

// RegisterControllers sets up controllers with the manager
func (s *ControllerRuntimeServer) RegisterControllers(ctx context.Context) error {
	// Get the runtime scheme from controller runtime
	scheme := s.controllerRuntime.GetManager().GetScheme()

	// Register deployment reconciler
	deploymentReconciler := controller.NewDeploymentReconciler(
		s.controllerRuntime.GetClient(),
		scheme,
		s.resourceService,
	)

	if err := s.controllerRuntime.RegisterDeploymentController(deploymentReconciler); err != nil {
		return fmt.Errorf("failed to register deployment controller: %w", err)
	}

	slog.Info("Controllers registered successfully")
	return nil
}

// SetupControllerRuntimeRoutes adds controller-runtime specific API endpoints
func (s *ControllerRuntimeServer) SetupControllerRuntimeRoutes() {
	// API version prefix
	api := s.app.Group("/api/v1")

	// Get metrics endpoint info
	metricsEndpoint := s.controllerRuntime.GetMetricsEndpoint()
	healthEndpoint := s.controllerRuntime.GetHealthEndpoint()

	// Add info endpoint about controller-runtime
	api.Get("/controller-runtime", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":           "running",
			"metrics_endpoint": metricsEndpoint,
			"health_endpoint":  healthEndpoint,
		})
	})

	// Add controller-runtime status endpoint
	api.Get("/controller", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      "Controller is running in the background",
			"reconcilers": []string{"deployment"},
		})
	})

	// Deployment endpoints using controller-runtime client
	deploymentAPI := api.Group("/deployments")

	// GET /api/v1/deployments
	deploymentAPI.Get("/", func(c *fiber.Ctx) error {
		// Get namespace from query param, default to "default"
		namespace := c.Query("namespace", "default")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Use controller-runtime client to list deployments
		var deploymentList appsv1.DeploymentList
		if err := s.controllerRuntime.GetClient().List(ctx, &deploymentList, &client.ListOptions{
			Namespace: namespace,
		}); err != nil {
			slog.Error("Failed to list deployments", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to list deployments",
				"details": err.Error(),
			})
		}

		// Convert to domain models
		deployments := make([]domain.Deployment, 0, len(deploymentList.Items))
		for _, d := range deploymentList.Items {
			deployments = append(deployments, domain.Deployment{
				Name:              d.Name,
				Namespace:         d.Namespace,
				Replicas:          *d.Spec.Replicas,
				ReadyReplicas:     d.Status.ReadyReplicas,
				AvailableReplicas: d.Status.AvailableReplicas,
				UpdatedReplicas:   d.Status.UpdatedReplicas,
				Labels:            d.Labels,
				CreationTimestamp: d.CreationTimestamp.String(),
			})
		}

		return c.JSON(fiber.Map{
			"deployments": deployments,
			"count":       len(deployments),
			"namespace":   namespace,
		})
	})

	// GET /api/v1/deployments/:name
	deploymentAPI.Get("/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		namespace := c.Query("namespace", "default")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Use controller-runtime client to get deployment
		var deployment appsv1.Deployment
		if err := s.controllerRuntime.GetClient().Get(ctx, client.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}, &deployment); err != nil {
			slog.Error("Failed to get deployment", "name", name, "namespace", namespace, "error", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Deployment not found",
				"details": err.Error(),
			})
		}

		// Convert to domain model
		deploymentModel := domain.Deployment{
			Name:              deployment.Name,
			Namespace:         deployment.Namespace,
			Replicas:          *deployment.Spec.Replicas,
			ReadyReplicas:     deployment.Status.ReadyReplicas,
			AvailableReplicas: deployment.Status.AvailableReplicas,
			UpdatedReplicas:   deployment.Status.UpdatedReplicas,
			Labels:            deployment.Labels,
			CreationTimestamp: deployment.CreationTimestamp.String(),
		}

		return c.JSON(deploymentModel)
	})

	// Status routes
	api.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service":            "k8s-controller",
			"status":             "running",
			"fiber_version":      fiber.Version,
			"controller_runtime": "active",
		})
	})
}

// Start begins the server and controller manager
func (s *ControllerRuntimeServer) Start() error {
	// Create a background context for controller runtime
	ctx := context.Background()

	// Start controller manager in a goroutine
	go func() {
		if err := s.controllerRuntime.Start(ctx); err != nil {
			slog.Error("Error starting controller manager", "error", err)
		}
	}()

	// Start the fiber server
	return s.Server.Start()
}

// Shutdown gracefully stops both the server and controller manager
func (s *ControllerRuntimeServer) Shutdown() error {
	// Stop controller manager
	s.controllerRuntime.Stop()

	// Stop the fiber server
	return s.Server.Shutdown()
}
