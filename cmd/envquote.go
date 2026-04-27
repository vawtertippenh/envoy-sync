package cmd

import (
	"fmt"
	"os"

	"github.com/yourorg/envoy-sync/internal/envfile"
	"github.com/yourorg/envoy-sync/internal/envquote"
	"github.com/spf13/cobra"
)

func init() {
	var (
		style         string
		forceAll      bool
		stripExisting bool
		outputFile    string
	)

	cmd := &cobra.Command{
		Use:   "envquote <file>",
		Short: "Apply consistent quoting to .env file values",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			opts := envquote.Options{
				Style:         envquote.Style(style),
				ForceAll:      forceAll,
				StripExisting: stripExisting,
			}

			result := envquote.Quote(env, opts)
			output := envquote.Render(result)

			if outputFile != "" {
				if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
					return fmt.Errorf("write: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "wrote %s — %d quoted, %d skipped\n",
					outputFile, result.Quoted, result.Skipped)
			} else {
				fmt.Fprint(cmd.OutOrStdout(), output)
				fmt.Fprintf(cmd.ErrOrStderr(), "# %d quoted, %d skipped\n",
					result.Quoted, result.Skipped)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&style, "style", "double", "quoting style: double, single, auto")
	cmd.Flags().BoolVar(&forceAll, "force", false, "quote every value regardless of content")
	cmd.Flags().BoolVar(&stripExisting, "strip", false, "strip existing quotes before re-quoting")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "write result to file instead of stdout")

	rootCmd.AddCommand(cmd)
}
