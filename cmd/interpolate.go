package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/interpolate"
)

func init() {
	var outputJSON bool

	interpolateCmd := &cobra.Command{
		Use:   "interpolate [file]",
		Short: "Resolve variable references in a .env file",
		Long: `Reads a .env file and resolves $VAR or ${VAR} references
within values, printing the fully interpolated environment.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			env, err := envfile.Parse(path)
			if err != nil {
				return fmt.Errorf("parsing %s: %w", path, err)
			}

			res := interpolate.Interpolate(env)

			for _, w := range res.Warnings {
				fmt.Fprintf(os.Stderr, "warning: %s\n", w)
			}

			if outputJSON {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(res.Env)
			}

			for k, v := range res.Env {
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
			}
			return nil
		},
	}

	interpolateCmd.Flags().BoolVar(&outputJSON, "json", false, "Output result as JSON")
	rootCmd.AddCommand(interpolateCmd)
}
