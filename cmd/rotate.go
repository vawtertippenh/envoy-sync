package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/rotate"
)

func init() {
	var (
		keys          []string
		extraPatterns []string
		length        int
		dryRun        bool
		outFile       string
	)

	cmd := &cobra.Command{
		Use:   "rotate <file>",
		Short: "Rotate secret values in an env file with freshly generated random secrets",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			opts := rotate.Options{
				Keys:          keys,
				ExtraPatterns: extraPatterns,
				Length:        length,
				DryRun:        dryRun,
			}

			out, result, err := rotate.Rotate(env, opts)
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Fprintln(cmd.OutOrStdout(), "[dry-run] Keys that would be rotated:")
				for _, k := range result.Rotated {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", k)
				}
				return nil
			}

			dest := args[0]
			if outFile != "" {
				dest = outFile
			}

			var sb strings.Builder
			for k, v := range out {
				sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
			}
			if err := os.WriteFile(dest, []byte(sb.String()), 0o600); err != nil {
				return fmt.Errorf("write: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Rotated %d key(s): %s\n",
				len(result.Rotated), strings.Join(result.Rotated, ", "))
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "Specific keys to rotate (default: all sensitive)")
	cmd.Flags().StringSliceVar(&extraPatterns, "extra-patterns", nil, "Extra regex patterns to identify sensitive keys")
	cmd.Flags().IntVarP(&length, "length", "l", 16, "Byte length of generated secret (hex output = 2x length)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Report what would be rotated without writing changes")
	cmd.Flags().StringVarP(&outFile, "output", "o", "", "Write result to this file instead of overwriting the source")

	rootCmd.AddCommand(cmd)
}
