package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/user/envoy-sync/internal/envfile"
	"github.com/user/envoy-sync/internal/envsample"
	"github.com/spf13/cobra"
)

func init() {
	var (
		n             int
		seed          int64
		deterministic bool
		forceKeys     []string
		quiet         bool
	)

	cmd := &cobra.Command{
		Use:   "envsample <file>",
		Short: "Sample a random subset of keys from an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			opts := envsample.Options{
				N:             n,
				Seed:          seed,
				Deterministic: deterministic || seed != 0,
				IncludeKeys:   forceKeys,
			}

			result, err := envsample.Sample(env, opts)
			if err != nil {
				return err
			}

			if !quiet {
				fmt.Fprintln(os.Stderr, "#", envsample.Summary(result))
			}

			keys := make([]string, 0, len(result.Env))
			for k := range result.Env {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, result.Env[k])
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&n, "count", "n", 0, "Number of keys to sample (0 = all)")
	cmd.Flags().Int64Var(&seed, "seed", 0, "Random seed for deterministic sampling")
	cmd.Flags().BoolVar(&deterministic, "deterministic", false, "Use deterministic sampling with seed")
	cmd.Flags().StringSliceVar(&forceKeys, "force", nil, "Keys to always include in the sample")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress summary output")

	rootCmd.AddCommand(cmd)
}
