// package cmd provides command-line interface for the application
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"k8s-controller/internal/infrastructure/config"
	"k8s-controller/internal/infrastructure/server"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  `Start the HTTP server for the Kubernetes controller API`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting HTTP server...")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			slog.Error("Failed to load configuration", "error", err)
			cfg = config.Default() // Use default config on error
		}

		// Create and configure server
		srv := server.NewServer(cfg.ServerPort)

		// Setup routes - this will also connect to Kubernetes
		slog.Info("Setting up routes and connecting to Kubernetes...")
		srv.SetupRoutes()
		slog.Info("Routes configured successfully")

		// Handle graceful shutdown
		go func() {
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

			<-sigCh
			slog.Info("Received shutdown signal")

			if err := srv.Shutdown(); err != nil {
				slog.Error("Error shutting down server", "error", err)
			}
		}()

		// Start the server
		slog.Info("Starting server", "port", cfg.ServerPort)
		if err := srv.Start(); err != nil {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Add flags for the serve command
	serveCmd.Flags().Int("port", 8080, "Port to run the server on")

	// Bind flags to viper config
	viper.BindPFlag("server.port", serveCmd.Flags().Lookup("port"))
}
