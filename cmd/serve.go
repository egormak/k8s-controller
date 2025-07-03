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

	"k8s-controller/internal/infrastructure/server"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  `Start the HTTP server for the Kubernetes controller API`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting HTTP server...")

		// Get port from config
		port := viper.GetInt("server.port")
		if port == 0 {
			port = 8080 // default port
		}

		// Create and configure server
		srv := server.NewServer(port)

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

		// Log server startup
		slog.Info("HTTP server starting", "port", port)

		// Start server (blocks until shutdown or error)
		if err := srv.Start(); err != nil {
			slog.Error("Server error", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Add any serve-specific flags
	serveCmd.Flags().IntP("port", "p", 0, "Port to run the HTTP server on (default is from config or 8080)")
	if err := viper.BindPFlag("server.port", serveCmd.Flags().Lookup("port")); err != nil {
		panic(fmt.Errorf("failed to bind server.port flag: %w", err))
	}
}
