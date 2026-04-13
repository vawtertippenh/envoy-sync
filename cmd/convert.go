package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/convert"
	"envoy-sync/internal/envfile"
)

func init() {
	var format string

	convertCmd := &cobra.Command{
		Use:   "convert [file]",
		Short: "Convert a .env file to another format",
		Long: `Convert a .env file into dotenv, json, yaml, or export format.

Example:
  envoy-sync convert .env --format json
  envoy-sync convert .env --format yaml
  envoy-sync convert .env --format export`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			env, err := envfile.Parse(filePath)
			if err != nil {
				return fmt.Errorf("failed to parse env file %q: %w", filePath, err)
			}

			result, err := convert.Convert(env, convert.Format(format))
			if err != nil {
				return fmt.Errorf("conversion failed: %w", err)
			}

			fmt.Fprint(os.Stdout, result.Output)
			return nil
		},
	}

	convertCmd.Flags().StringVarP(&format, "format", "f", "dotenv",
		"Output format: dotenv, json, yaml, export")

	rootCmd.AddCommand(convertCmd)
}
