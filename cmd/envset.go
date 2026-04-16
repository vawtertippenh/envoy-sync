package cmd

import (
	"fmt"
	"os"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envset"
	"github.com/spf13/cobra"
)

func init() {
	var op string

	cmd := &cobra.Command{
		Use:   "envset <fileA> <fileB>",
		Short: "Perform set operations (union, intersect, diff) on two .env files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse %s: %w", args[0], err)
			}
			b, err := envfile.Parse(args[1])
			if err != nil {
				return fmt.Errorf("parse %s: %w", args[1], err)
			}

			var result envset.Result
			switch op {
			case "union":
				result = envset.Union(a, b)
			case "intersect":
				result = envset.Intersect(a, b)
			case "diff":
				result = envset.Difference(a, b)
			default:
				return fmt.Errorf("unknown operation %q: use union, intersect, or diff", op)
			}

			for _, k := range result.Keys {
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, result.Env[k])
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&op, "op", "o", "union", "Set operation: union, intersect, diff")
	rootCmd.AddCommand(cmd)
}
