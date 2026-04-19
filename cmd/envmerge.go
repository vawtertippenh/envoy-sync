package cmd

import (
	"fmt"
	"os"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envmerge"

	"github.com/spf13/cobra"
)

func init() {
	var strategy string
	var prefix string

	cmd := &cobra.Command{
		Use:   "envmerge <file1> <file2> [files...]",
		Short: "Merge multiple .env files with conflict resolution",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var sources []map[string]string
			for _, path := range args {
				env, err := envfile.Parse(path)
				if err != nil {
					return fmt.Errorf("parse %s: %w", path, err)
				}
				sources = append(sources, env)
			}

			opts := envmerge.Options{
				Strategy: envmerge.Strategy(strategy),
				Prefix:   prefix,
			}

			res, err := envmerge.Merge(sources, opts)
			if err != nil {
				return err
			}

			if len(res.Conflicts) > 0 {
				fmt.Fprintln(os.Stderr, "# conflicts:")
				for _, c := range res.Conflicts {
					fmt.Fprintf(os.Stderr, "#   %s: %v\n", c.Key, c.Values)
				}
			}

			for k, v := range res.Env {
				fmt.Printf("%s=%s\n", k, v)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&strategy, "strategy", "last", "conflict strategy: first|last|error")
	cmd.Flags().StringVar(&prefix, "prefix", "", "prefix to prepend to all keys")
	rootCmd.AddCommand(cmd)
}
