// package cmd provides command-line interface for the application
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"k8s-controller/internal/app"
	"k8s-controller/internal/infrastructure/config"
)

// controlCmd represents the control command
var controlCmd = &cobra.Command{
	Use:   "control",
	Short: "Start the Kubernetes controller",
	Long: `Start the Kubernetes controller which will watch for resources
and process them according to the defined business logic.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Kubernetes controller...")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			slog.Error("Failed to load configuration", "error", err)
			cfg = config.Default() // Use default config on error
		}

		// Create controller with config
		controller := app.NewKubernetesController(cfg)

		// Start controller
		if err := controller.Start(); err != nil {
			slog.Error("Failed to start controller", "error", err)
			os.Exit(1)
		}

		// Keep running until shutdown signal
		slog.Info("Controller is running. Press Ctrl+C to stop.")

		// Wait for termination signal
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		slog.Info("Shutting down controller...")
		controller.Stop()

		// Give time for graceful shutdown
		time.Sleep(2 * time.Second)
		slog.Info("Controller stopped")
	},
}

func init() {
	rootCmd.AddCommand(controlCmd)

	// Add flags specific to controller functionality
	controlCmd.Flags().StringSlice("namespaces", []string{"default"}, "Namespaces to watch (comma-separated)")
	controlCmd.Flags().StringSlice("resources", []string{"deployments,services,pods"}, "Resources to watch (comma-separated)")

	// Add leader election flags
	controlCmd.Flags().Bool("leader-elect", false, "Enable leader election for controller")
	controlCmd.Flags().String("leader-election-id", "k8s-controller", "ID to use for leader election")
	controlCmd.Flags().String("leader-election-namespace", "default", "Namespace to use for leader election")

	// Bind flags to viper config
	if err := viper.BindPFlag("kubernetes.namespaces", controlCmd.Flags().Lookup("namespaces")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("kubernetes.resources", controlCmd.Flags().Lookup("resources")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("leader-election.enabled", controlCmd.Flags().Lookup("leader-elect")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("leader-election.id", controlCmd.Flags().Lookup("leader-election-id")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("leader-election.namespace", controlCmd.Flags().Lookup("leader-election-namespace")); err != nil {
		panic(err)
	}
}
