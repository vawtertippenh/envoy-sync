package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/merge"
)

func init() {
	var strategy string

	mergeCmd := &cobra.Command{
		Use:   "merge <file1> <file2> [file3...]",
		Short: "Merge multiple .env files into one",
		Long: `Merge two or more .env files into a single output.

Conflict resolution strategies:
  first  — keep the value from the earliest file (default)
  last   — keep the value from the latest file
  error  — abort if any key has conflicting values`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			sources := make([]map[string]string, 0, len(args))
			for _, path := range args {
				env, err := envfile.Parse(path)
				if err != nil {
					return fmt.Errorf("parsing %s: %w", path, err)
				}
				sources = append(sources, env)
			}

			res, err := merge.Merge(sources, merge.Strategy(strategy))
			if err != nil {
				return err
			}

			if len(res.Conflicts) > 0 {
				fmt.Fprintln(os.Stderr, "# Conflicts detected:")
				for _, c := range res.Conflicts {
					fmt.Fprintf(os.Stderr, "#   %s: %v\n", c.Key, c.Values)
				}
			}

			for k, v := range res.Env {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
			}
			return nil
		},
	}

	mergeCmd.Flags().StringVarP(&strategy, "strategy", "s", "first",
		"conflict resolution strategy: first, last, or error")

	rootCmd.AddCommand(mergeCmd)
}
