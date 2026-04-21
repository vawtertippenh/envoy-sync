package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	envsorter "envoy-sync/internal/envsorter"
)

func init() {
	var (
		strategy   string
		descending bool
		prefixSep  string
	)

	cmd := &cobra.Command{
		Use:   "envsorter <file>",
		Short: "Sort environment variables by key, value, length, or prefix",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			opts := envsorter.Options{
				Strategy:   envsorter.Strategy(strategy),
				Descending: descending,
				PrefixSep:  prefixSep,
			}

			res := envsorter.Sort(env, opts)

			for _, k := range res.Keys {
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, res.Env[k])
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&strategy, "strategy", "s", "alpha",
		"Sort strategy: alpha, value, length, prefix")
	cmd.Flags().BoolVarP(&descending, "desc", "d", false,
		"Sort in descending order")
	cmd.Flags().StringVar(&prefixSep, "prefix-sep", "_",
		"Separator used to detect key prefixes (used with --strategy=prefix)")

	rootCmd.AddCommand(cmd)
}
