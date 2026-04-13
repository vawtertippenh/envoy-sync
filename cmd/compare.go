package cmd

import (
	"fmt"
	"os"

	"envoy-sync/internal/compare"
	"envoy-sync/internal/envfile"

	"github.com/spf13/cobra"
)

func init() {
	var templateFile string

	compareCmd := &cobra.Command{
		Use:   "compare [target]",
		Short: "Compare a target .env file against a template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if templateFile == "" {
				return fmt.Errorf("--template flag is required")
			}

			tmpl, err := envfile.Parse(templateFile)
			if err != nil {
				return fmt.Errorf("parsing template: %w", err)
			}

			target, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parsing target: %w", err)
			}

			r := compare.Against(tmpl, target)
			fmt.Println(r.Summary())

			if len(r.Missing) > 0 || len(r.Mismatch) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	compareCmd.Flags().StringVarP(&templateFile, "template", "t", "", "path to the template .env file")
	rootCmd.AddCommand(compareCmd)
}
