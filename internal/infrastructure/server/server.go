// package server provides HTTP server functionality using Fiber
package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"k8s-controller/internal/infrastructure/kubernetes"
)

// Server represents an HTTP server
type Server struct {
	app            *fiber.App
	port           int
	kubeClient     kubernetes.Client
	deploymentCtrl *DeploymentController
}

// NewServer creates a new HTTP server instance
func NewServer(port int) *Server {
	// Create Kubernetes client
	kubeClient := kubernetes.NewClient()

	// Initialize controllers
	deploymentCtrl := NewDeploymentController(kubeClient)

	app := fiber.New(fiber.Config{
		AppName: "K8s Controller API",
	})

	// Add middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	return &Server{
		app:            app,
		port:           port,
		kubeClient:     kubeClient,
		deploymentCtrl: deploymentCtrl,
	}
}

// SetupRoutes configures all the routes for the server
func (s *Server) SetupRoutes() {
	// Connect to Kubernetes with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.kubeClient.Connect(ctx); err != nil {
		slog.Error("Failed to connect to Kubernetes cluster", "error", err)
		// We'll continue setup but K8s-related endpoints will fail
	} else {
		slog.Info("Successfully connected to Kubernetes cluster")
	}

	// Initialize informers with 30s resync period in a separate goroutine
	// This prevents blocking the server startup
	go func() {
		informerCtx := context.Background()
		// Default to "default" namespace, but could be configured to watch more
		namespaces := []string{"default"}
		if err := s.kubeClient.InitializeInformers(informerCtx, namespaces); err != nil {
			slog.Error("Failed to initialize informers", "error", err)
		} else {
			slog.Info("Successfully initialized informers")
		}
	}()

	// API group
	api := s.app.Group("/api")

	// Health check endpoint
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// Version endpoint
	api.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"version": "v0.1.0",
		})
	})

	// Kubernetes resource endpoints
	k8s := api.Group("/k8s")

	// Deployment endpoints
	deployments := k8s.Group("/deployments")
	deployments.Get("/", s.deploymentCtrl.ListDeployments)

	// Root route - redirect to /api/health
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/api/health")
	})
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	slog.Info("Starting HTTP server", "port", s.port)

	address := fmt.Sprintf(":%d", s.port)
	return s.app.Listen(address)
}

// Shutdown gracefully stops the HTTP server
func (s *Server) Shutdown() error {
	slog.Info("Shutting down HTTP server")
	return s.app.Shutdown()
}
