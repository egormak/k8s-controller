// package server provides HTTP server functionality using Fiber
package server

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Server represents an HTTP server
type Server struct {
	app  *fiber.App
	port int
}

// NewServer creates a new HTTP server instance
func NewServer(port int) *Server {
	app := fiber.New(fiber.Config{
		AppName: "K8s Controller API",
	})

	// Add middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	return &Server{
		app:  app,
		port: port,
	}
}

// SetupRoutes configures all the routes for the server
func (s *Server) SetupRoutes() {
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
