package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-sync/internal/flatten"
)

func init() {
	var separator string
	var prefix string
	var upper bool
	var outputFile string

	cmd := &cobra.Command{
		Use:   "flatten <input.json>",
		Short: "Flatten a nested JSON file into a dotenv-style KEY=VALUE file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("reading input: %w", err)
			}

			opts := flatten.Options{
				Separator: separator,
				Prefix:    prefix,
				UpperCase: upper,
			}

			flat, err := flatten.FlattenJSON(data, opts)
			if err != nil {
				return err
			}

			rendered := flatten.Render(flat)

			if outputFile != "" {
				if err := os.WriteFile(outputFile, []byte(rendered), 0o644); err != nil {
					return fmt.Errorf("writing output: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Written %d keys to %s\n", len(flat), outputFile)
			} else {
				fmt.Fprint(cmd.OutOrStdout(), rendered)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&separator, "separator", "s", "_", "Separator between nested key segments")
	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Prefix to prepend to all output keys")
	cmd.Flags().BoolVarP(&upper, "upper", "u", true, "Convert keys to upper-case")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write output to file instead of stdout")

	rootCmd.AddCommand(cmd)
}
