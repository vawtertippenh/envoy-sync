// Package cmd implements the CLI commands for envoy-sync.
// It uses the cobra library to define the root command and subcommands.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// maskSecrets controls whether secret values are masked in output
	maskSecrets bool
	// secretKeys holds the list of key patterns to treat as secrets
	secretKeys []string
)

// rootCmd is the base command for the envoy-sync CLI.
var rootCmd = &cobra.Command{
	Use:   "envoy-sync",
	Short: "Diff and sync .env files across environments",
	Long: `envoy-sync is a CLI tool for comparing and synchronizing .env files
across different environments. It supports secret masking to prevent
sensitive values from being exposed in output or logs.

Examples:
  envoy-sync diff .env.local .env.production
  envoy-sync diff --mask .env.staging .env.production
  envoy-sync sync .env.local .env.production`,
	SilenceUsage: true,
}

// Execute runs the root command and exits with a non-zero status on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Persistent flags are available to all subcommands
	rootCmd.PersistentFlags().BoolVarP(
		&maskSecrets,
		"mask", "m",
		false,
		"mask secret values in output",
	)
	rootCmd.PersistentFlags().StringSliceVarP(
		&secretKeys,
		"secret-keys", "s",
		[]string{"PASSWORD", "SECRET", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL"},
		"comma-separated list of key substrings to treat as secrets (case-insensitive)",
	)
}
