package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "k8s-controller",
	Short: "K8s Controller is a CLI tool for k8s",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to k8s-controller. Use --help for usage.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("error in command", "error", err)
		os.Exit(1)
	}
}
