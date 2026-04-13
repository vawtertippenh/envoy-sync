package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/validate"
)

func init() {
	var requiredKeys []string

	validateCmd := &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate an .env file for common issues",
		Long: `Checks an .env file for:
  - Non-uppercase keys (warning)
  - Empty values (warning)
  - Missing required keys (error)

Example:
  envoy-sync validate .env --require DATABASE_URL --require PORT`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			env, err := envfile.Parse(filePath)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", filePath, err)
			}

			// Flatten required keys in case comma-separated values are passed.
			var required []string
			for _, r := range requiredKeys {
				for _, part := range strings.Split(r, ",") {
					part = strings.TrimSpace(part)
					if part != "" {
						required = append(required, part)
					}
				}
			}

			report := validate.Validate(filePath, env, required)
			fmt.Println(report.Summary())

			if report.HasErrors() {
				os.Exit(1)
			}
			return nil
		},
	}

	validateCmd.Flags().StringArrayVar(
		&requiredKeys, "require", []string{},
		"Keys that must be present (can be specified multiple times)",
	)

	rootCmd.AddCommand(validateCmd)
}
