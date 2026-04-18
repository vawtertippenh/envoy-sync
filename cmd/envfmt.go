package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envfmt"
)

func init() {
	var uppercase bool
	var sortKeys bool
	var quoteAll bool
	var spaceEq bool
	var write bool

	cmd := &cobra.Command{
		Use:   "envfmt [file]",
		Short: "Format a .env file according to style options",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			style := envfmt.Style{
				UppercaseKeys:    uppercase,
				SortKeys:         sortKeys,
				QuoteAllValues:   quoteAll,
				SpaceAroundEqual: spaceEq,
			}

			result := envfmt.Format(env, style)
			output := envfmt.Render(result)

			if write {
				if err := os.WriteFile(args[0], []byte(output), 0644); err != nil {
					return fmt.Errorf("write: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "formatted %s (%d/%d keys changed)\n", args[0], result.Changed, result.Total)
			} else {
				fmt.Fprint(cmd.OutOrStdout(), output)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&uppercase, "uppercase", false, "Uppercase all keys")
	cmd.Flags().BoolVar(&sortKeys, "sort", false, "Sort keys alphabetically")
	cmd.Flags().BoolVar(&quoteAll, "quote", false, "Quote all values")
	cmd.Flags().BoolVar(&spaceEq, "space-eq", false, "Add spaces around = sign")
	cmd.Flags().BoolVarP(&write, "write", "w", false, "Write result back to file")

	rootCmd.AddCommand(cmd)
}
