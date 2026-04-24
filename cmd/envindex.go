package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-sync/internal/envfile"
	"github.com/yourorg/envoy-sync/internal/envindex"
)

func init() {
	var prefix, suffix, substring string

	cmd := &cobra.Command{
		Use:   "envindex <file>",
		Short: "Search and index keys in a .env file by prefix, suffix, or substring",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			idx := envindex.Build(env)

			var results []envindex.Entry
			switch {
			case prefix != "":
				results = idx.ByPrefix(prefix)
			case suffix != "":
				results = idx.BySuffix(suffix)
			case substring != "":
				results = idx.BySubstring(substring)
			default:
				results = idx.All()
			}

			if len(results) == 0 {
				fmt.Fprintln(os.Stderr, "no matching keys found")
				return nil
			}

			for _, e := range results {
				fmt.Printf("%s=%s\n", e.Key, e.Value)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "filter keys by prefix")
	cmd.Flags().StringVar(&suffix, "suffix", "", "filter keys by suffix")
	cmd.Flags().StringVar(&substring, "substring", "", "filter keys by substring")

	rootCmd.AddCommand(cmd)
}
