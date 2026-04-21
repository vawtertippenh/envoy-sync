package cmd

import (
	"fmt"
	"os"
	"strings"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envlookup"

	"github.com/spf13/cobra"
)

func init() {
	var (
		keys          []string
		maskSensitive bool
		caseFold      bool
		extraPatterns []string
	)

	cmd := &cobra.Command{
		Use:   "envlookup <file>",
		Short: "Look up specific keys in a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			opts := envlookup.Options{
				Keys:              keys,
				MaskSensitive:     maskSensitive,
				CaseFold:          caseFold,
				SensitivePatterns: extraPatterns,
			}

			results := envlookup.Lookup(env, opts)

			missing := 0
			for _, r := range results {
				if !r.Found {
					missing++
				}
			}

			fmt.Print(envlookup.Render(results))

			if missing > 0 {
				fmt.Fprintf(os.Stderr, "\n%d key(s) not found\n", missing)
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&keys, "key", "k", nil,
		"Keys to look up (comma-separated or repeated flag; omit for all)")
	cmd.Flags().BoolVarP(&maskSensitive, "mask", "m", false,
		"Mask values of sensitive keys")
	cmd.Flags().BoolVar(&caseFold, "case-fold", false,
		"Case-insensitive key matching")
	cmd.Flags().StringSliceVar(&extraPatterns, "sensitive", nil,
		"Extra patterns to treat as sensitive (comma-separated)")

	_ = strings.ToLower // satisfy import
	rootCmd.AddCommand(cmd)
}
