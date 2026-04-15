package cmd

import (
	"fmt"
	"os"
	"strings"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/strip"

	"github.com/spf13/cobra"
)

func init() {
	var (
		keys     []string
		prefixes []string
		suffixes []string
		output   string
	)

	cmd := &cobra.Command{
		Use:   "strip <file>",
		Short: "Remove keys from a .env file by name, prefix, or suffix",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			res := strip.Strip(env, strip.Options{
				Keys:           keys,
				RemovePrefixes: prefixes,
				RemoveSuffixes: suffixes,
			})

			if len(res.Removed) > 0 {
				fmt.Fprintf(os.Stderr, "stripped %d key(s): %s\n",
					len(res..Removed, ", "))
			}

			var sb strings.Builder
			for k, v := range res.Env {
				sb.WriteString(k + "=" + v + "\n")
			}

			if output == "" || output == "-" {
				fmt.Print(sb.String())
				return nil
			}

			if err := os.WriteFile(output, []byte(sb.String()), 0o644); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			fmt.Fprintf(os.Stderr, "written to %s\n", output)
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&keys, "key", "k", nil, "explicit key(s) to remove")
	cmd.Flags().StringSliceVar(&prefixes, "prefix", nil, "remove keys with this prefix")
	cmd.Flags().StringSliceVar(&suffixes, "suffix", nil, "remove keys with this suffix")
	cmd.Flags().StringVarP(&output, "output", "o", "-", "output file (default: stdout)")

	rootCmd.AddCommand(cmd)
}
