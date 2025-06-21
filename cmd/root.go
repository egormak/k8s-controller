package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-controller",
	Short: "k8s-controller is a CLI tool for managing Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to k8s-controller. Use --help for usage.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		slog.Error("error in command", "error", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/k8s-controller.yaml)")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	// viper.BindEnv("config", "CONFIG")

	// cfg := viper.GetString("config")
	// fmt.Println("Using file:", cfg)
	// fmt.Println("SERVER_PORT:", viper.GetString("SERVER_PORT"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find current directory.
		currentDir, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in home directory with name "k8s-config" (without extension).
		viper.AddConfigPath(currentDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("k8s-config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
