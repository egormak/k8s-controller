package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"k8s-controller/internal/infrastructure/kubernetes"
)

var namespace string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Kubernetes resources",
	Long:  `List Kubernetes resources like deployments, services, pods`,
}

// deploymentCmd represents the deployment subcommand
var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "List deployments",
	Long:  `List deployments in the specified namespace`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Listing deployments in namespace: %s\n", namespace)

		// Create Kubernetes client
		client := kubernetes.NewClient()

		// Connect to cluster
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := client.Connect(ctx); err != nil {
			slog.Error("Failed to connect to Kubernetes cluster", "error", err)
			os.Exit(1)
		}

		// List deployments
		deployments, err := client.ListDeployments(ctx, namespace)
		if err != nil {
			slog.Error("Failed to list deployments", "error", err, "namespace", namespace)
			os.Exit(1)
		}

		// Display results
		if len(deployments) == 0 {
			fmt.Printf("No deployments found in namespace '%s'\n", namespace)
			return
		}

		fmt.Printf("Found %d deployment(s) in namespace '%s':\n", len(deployments), namespace)
		fmt.Printf("%-30s %-10s %-10s %-10s\n", "NAME", "READY", "UP-TO-DATE", "AVAILABLE")
		fmt.Println("--------------------------------------------------------------------------------")

		for _, deployment := range deployments {
			fmt.Printf("%-30s %-10d %-10d %-10d\n",
				deployment.Name,
				deployment.ReadyReplicas,
				deployment.UpdatedReplicas,
				deployment.AvailableReplicas)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(deploymentCmd)

	// Add namespace flag to both list and deployment commands
	listCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")
	deploymentCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")

	// Bind flags to viper
	if err := viper.BindPFlag("kubernetes.namespace", listCmd.PersistentFlags().Lookup("namespace")); err != nil {
		panic(fmt.Errorf("failed to bind namespace flag: %w", err))
	}
}
