package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-sync/internal/envfile"
	"github.com/yourorg/envoy-sync/internal/envsqueeze"
)

func init() {
	var (
		dedupeValues       bool
		removeEmpty        bool
		removePlaceholders bool
	)

	cmd := &cobra.Command{
		Use:   "squeeze <file>",
		Short: "Remove redundant entries from an .env file",
		Long: `squeeze reads an .env file and removes redundant entries based on the
selected options: empty values, placeholder strings, or keys sharing
the same value (keeping only the first alphabetically).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			opts := envсqueeze.Options{
				DedupeValues:       dedupeValues,
				RemoveEmpty:        removeEmpty,
				RemovePlaceholders: removePlaceholders,
			}
			res := envsqueeze.Squeeze(env, opts)

			if len(res.Dropped) > 0 {
				fmt.Fprintf(os.Stderr, "Dropped %d key(s): %v\n", len(res.Dropped), res.Dropped)
			}

			keys := make([]string, 0, len(res.Env))
			for k := range res.Env {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Printf("%s=%s\n", k, res.Env[k])
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&dedupeValues, "dedupe-values", false, "remove keys with duplicate values, keeping the first alphabetically")
	cmd.Flags().BoolVar(&removeEmpty, "remove-empty", false, "remove keys with empty values")
	cmd.Flags().BoolVar(&removePlaceholders, "remove-placeholders", false, "remove keys with placeholder values (e.g. CHANGE_ME)")

	rootCmd.AddCommand(cmd)
}
