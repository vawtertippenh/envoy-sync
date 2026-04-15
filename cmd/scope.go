package cmd

import (
	"fmt"
	"os"
	"strings"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/scope"

	"github.com/spf13/cobra"
)

func init() {
	var prefixes []string
	var suffixes []string
	var strip bool
	var showUnmatched bool

	var scopeCmd = &cobra.Command{
		Use:   "scope <file>",
		Short: "Filter env keys by prefix or suffix patterns",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			r := scope.Scope(env, scope.Options{
				Prefixes: prefixes,
				Suffixes: suffixes,
				Strip:    strip,
			})

			fmt.Fprintf(os.Stderr, "# %s\n", r.Summary())

			target := r.Matched
			if showUnmatched {
				target = r.Unmatched
			}

			for k, v := range target {
				if strings.ContainsAny(v, " \t#") {
					fmt.Printf("%s=\"%s\"\n", k, v)
				} else {
					fmt.Printf("%s=%s\n", k, v)
				}
			}
			return nil
		},
	}

	scopeCmd.Flags().StringSliceVar(&prefixes, "prefix", nil, "Key prefix patterns to match (comma-separated)")
	scopeCmd.Flags().StringSliceVar(&suffixes, "suffix", nil, "Key suffix patterns to match (comma-separated)")
	scopeCmd.Flags().BoolVar(&strip, "strip", false, "Strip matched prefix from output keys")
	scopeCmd.Flags().BoolVar(&showUnmatched, "unmatched", false, "Output unmatched keys instead of matched")

	rootCmd.AddCommand(scopeCmd)
}
