package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/lint"
)

func init() {
	var failOnError bool

	lintCmd := &cobra.Command{
		Use:   "lint [file]",
		Short: "Lint an .env file for style and correctness issues",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			env, err := envfile.Parse(path)
			if err != nil {
				return fmt.Errorf("failed to parse env file: %w", err)
			}

			result := lint.Lint(env)

			fmt.Println(result.Summary())

			if failOnError && result.HasErrors() {
				os.Exit(1)
			}

			return nil
		},
	}

	lintCmd.Flags().BoolVar(&failOnError, "fail-on-error", false, "Exit with code 1 if any errors are found")

	rootCmd.AddCommand(lintCmd)
}
