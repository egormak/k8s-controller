package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"k8s-controller/internal/app"
)

var cfgFile string
var logLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-controller",
	Short: "k8s-controller is a CLI tool for managing Kubernetes",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupLogger()
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to k8s-controller. Use --help for usage.")
		slog.Info("k8s-controller started")

		controller := app.NewKubernetesController()
		if err := controller.Run(); err != nil {
			slog.Error("Failed to run controller", "error", err)
			os.Exit(1)
		}
	},
}

// setupLogger configures the slog logger based on the specified log level
func setupLogger() {
	var level slog.Level
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}
	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Debug("Logger initialized", "level", level.String())
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
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "INFO", "Set the logging level (DEBUG, INFO, WARN, ERROR)")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))

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

	// Override log level from config if present
	if viper.IsSet("log.level") {
		logLevel = viper.GetString("log.level")
	}
}
