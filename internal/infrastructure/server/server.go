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
		AppName:               "K8s Controller API",
		DisableStartupMessage: false,
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

// SetupRoutes configures the HTTP routes
func (s *Server) SetupRoutes() {
	// API version prefix
	api := s.app.Group("/api/v1")

	// Connect to Kubernetes
	if err := s.kubeClient.Connect(context.Background()); err != nil {
		slog.Error("Failed to connect to Kubernetes", "error", err)
		return
	}

	// Initialize informers for default namespace
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.kubeClient.InitializeInformers(ctx, []string{"default"}); err != nil {
		slog.Warn("Failed to initialize informers", "error", err)
		// Continue anyway, we'll use direct API calls
	}

	// Health check
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "timestamp": time.Now().Format(time.RFC3339)})
	})

	// Deployments
	api.Get("/deployments", s.deploymentCtrl.ListDeployments)
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	return s.app.Listen(fmt.Sprintf(":%d", s.port))
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
