// package cmd provides command-line interface for the application
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"k8s-controller/internal/app"
)

// controlCmd represents the control command
var controlCmd = &cobra.Command{
	Use:   "control",
	Short: "Start the Kubernetes controller",
	Long:  `Start the Kubernetes controller to watch and manage resources`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Kubernetes controller...")

		// Create the controller
		controller := app.NewKubernetesController()

		// Handle graceful shutdown
		go func() {
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

			<-sigCh
			slog.Info("Received shutdown signal, stopping controller")
			controller.Stop()
		}()

		// Run the controller (blocks until stopped or error)
		if err := controller.Run(); err != nil {
			slog.Error("Controller error", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(controlCmd)

	// Add any control-specific flags here if needed
}
