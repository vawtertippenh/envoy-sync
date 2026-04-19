package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envoy-sync/internal/envdrop"
	"github.com/user/envoy-sync/internal/envfile"
)

func init() {
	var file string
	var keys, prefixes, suffixes []string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "envdrop",
		Short: "Drop keys from a .env file by name, prefix, or suffix",
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(file)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			res := envdrop.Drop(env, envdrop.Options{
				Keys:     keys,
				Prefixes: prefixes,
				Suffixes: suffixes,
				DryRun:   dryRun,
			})

			if dryRun {
				fmt.Printf("Would drop %d key(s): %s\n", len(res.Dropped), strings.Join(res.Dropped, ", "))
				return nil
			}

			for _, k := range sortedMapKeys(res.Out) {
				fmt.Printf("%s=%s\n", k, res.Out[k])
			}
			if len(res.Dropped) > 0 {
				fmt.Fprintf(cmd.ErrOrStderr(), "# dropped: %s\n", strings.Join(res.Dropped, ", "))
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", ".env", "Path to .env file")
	cmd.Flags().StringSliceVarP(&keys, "key", "k", nil, "Exact key names to drop")
	cmd.Flags().StringSliceVar(&prefixes, "prefix", nil, "Drop keys with this prefix")
	cmd.Flags().StringSliceVar(&suffixes, "suffix", nil, "Drop keys with this suffix")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview which keys would be dropped")

	rootCmd.AddCommand(cmd)
}

func sortedMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	import_sort_strings(keys)
	return keys
}
