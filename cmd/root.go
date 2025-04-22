// cmd/root.go

// Package cmd contains the command implementations for the changie CLI application.
//
// This package uses the cobra library to define commands, subcommands, and flags.
// Each command is self-contained and follows a consistent pattern:
// - Command declaration with descriptive help text
// - Flag initialization with proper defaults and help text
// - Viper binding for configuration management
// - Command execution logic
//
// The package follows a hierarchical structure with root as the base command,
// and various subcommands for specific functionality.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/peiman/changie/internal/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// cfgFile holds the path to the configuration file
	cfgFile string

	// Version contains the current version of the application.
	// This is populated at build time via ldflags.
	Version = "dev"

	// Commit contains the Git commit hash of the build.
	// This is populated at build time via ldflags.
	Commit = ""

	// Date contains the build timestamp.
	// This is populated at build time via ldflags.
	Date = ""

	// binaryName is the name of the binary.
	// This is populated at build time via ldflags from Taskfile.yml.
	binaryName = "changie"
)

// RootCmd represents the base command when called without any subcommands.
// It is exported so that tests in other packages can manipulate it.
var RootCmd = &cobra.Command{
	Use:   binaryName,
	Short: "A professional changelog management tool for SemVer projects",
	Long: fmt.Sprintf(`%s is a command-line tool for managing changelogs following the Keep a Changelog format and Semantic Versioning.

It helps you automate changelog entries, version bumping, and Git tag management while maintaining a clean, consistent format.
For more information on the format, see https://keepachangelog.com`, binaryName),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := initConfig(); err != nil {
			return err
		}
		if err := logger.Init(nil); err != nil {
			return fmt.Errorf("failed to initialize logger: %w", err)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and runs it.
// This is called by main.main(). It returns an error if there was
// a problem during execution.
func Execute() error {
	RootCmd.Version = fmt.Sprintf("%s, commit %s, built at %s", Version, Commit, Date)
	return RootCmd.Execute()
}

// init is called automatically when the package is imported.
// It initializes the command flags and binds them to viper.
func init() {
	// Define persistent flags for the root command
	// These flags are available to all subcommands
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("Config file (default is $HOME/.%s.yaml)", binaryName))
	if err := viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config")); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind 'config' flag")
	}

	RootCmd.PersistentFlags().String("log-level", "info", "Set the log level (trace, debug, info, warn, error, fatal, panic)")
	if err := viper.BindPFlag("app.log_level", RootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind 'log-level'")
	}
}

// initConfig reads in config file and ENV variables if set.
// It follows a priority order:
// 1. Command-line flags
// 2. Environment variables
// 3. Configuration file
// 4. Default values
func initConfig() error {
	// Use config file from the flag if provided
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Otherwise look for config in home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(fmt.Sprintf(".%s", binaryName))
	}

	// Configure viper to read environment variables
	// Environment variables are converted using this pattern:
	// app.log_level becomes APP_LOG_LEVEL
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("app.log_level", "info")

	// Try to read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error and use defaults
			log.Info().Msg("No config file found, using defaults and environment variables")
		} else {
			// Config file was found but another error occurred
			log.Error().Err(err).Msg("Failed to read config file")
			return fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		log.Info().Str("config_file", viper.ConfigFileUsed()).Msg("Using config file")
	}

	return nil
}
