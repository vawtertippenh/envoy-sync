package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/filter"

	"github.com/spf13/cobra"
)

func init() {
	var patterns []string
	var regex string
	var invert bool
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "filter <file>",
		Short: "Filter env keys by pattern, suffix, or regex",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			res, err := filter.Filter(env, filter.Options{
				Patterns: patterns,
				Regex:    regex,
				Invert:   invert,
			})
			if err != nil {
				return fmt.Errorf("filter: %w", err)
			}

			if outputJSON {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(res.Env)
			}

			keys := make([]string, 0, len(res.Env))
			for k := range res.Env {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, res.Env[k])
			}
			fmt.Fprintf(os.Stderr, "# matched: %d  dropped: %d\n", res.Matched, res.Dropped)
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&patterns, "pattern", "p", nil, "Glob-style pattern(s) to match keys (e.g. DB_*, *_SECRET)")
	cmd.Flags().StringVarP(&regex, "regex", "r", "", "Regular expression to match keys")
	cmd.Flags().BoolVarP(&invert, "invert", "v", false, "Invert match — keep non-matching keys")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output result as JSON")

	rootCmd.AddCommand(cmd)
}
