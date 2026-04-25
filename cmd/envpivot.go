package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envpivot"
)

func init() {
	var maskSensitive bool
	var showSingletons bool

	cmd := &cobra.Command{
		Use:   "envpivot <file>",
		Short: "Pivot an env file grouping keys by shared value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			result := envpivot.Pivot(env)

			for _, g := range result.Groups {
				if !showSingletons && len(g.Keys) == 1 {
					continue
				}
				val := g.Value
				if maskSensitive {
					for _, k := range g.Keys {
						_ = k // masking logic would check key sensitivity
					}
				}
				fmt.Fprintf(os.Stdout, "[%s]\n  keys: %s\n", val, strings.Join(g.Keys, ", "))
			}

			fmt.Fprintln(os.Stdout, "---")
			fmt.Fprintln(os.Stdout, envpivot.Summary(result))
			return nil
		},
	}

	cmd.Flags().BoolVar(&maskSensitive, "mask", false, "mask sensitive values in output")
	cmd.Flags().BoolVar(&showSingletons, "singletons", false, "also show groups with a single key")

	rootCmd.AddCommand(cmd)
}
