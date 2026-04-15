package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/resolve"
)

func init() {
	var overrideFiles []string
	var overrideNames []string
	var showSummary bool

	cmd := &cobra.Command{
		Use:   "resolve <base.env>",
		Short: "Resolve an env file against one or more override sources",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			base, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parsing base file: %w", err)
			}

			if len(overrideFiles) != len(overrideNames) {
				return fmt.Errorf("--override and --name flags must be provided the same number of times")
			}

			sources := make([]resolve.Source, 0, len(overrideFiles))
			for i, f := range overrideFiles {
				vals, err := envfile.Parse(f)
				if err != nil {
					return fmt.Errorf("parsing override file %q: %w", f, err)
				}
				sources = append(sources, resolve.Source{Name: overrideNames[i], Values: vals})
			}

			result := resolve.Resolve(base, sources)

			if showSummary {
				lines := resolve.Summary(result)
				if len(lines) == 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "# no overrides applied")
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), "# overrides applied:")
					for _, l := range lines {
						fmt.Fprintln(cmd.OutOrStdout(), "#  "+l)
					}
				}
			}

			for k, v := range result.Env {
				if strings.ContainsAny(v, " \t") {
					v = `"` + v + `"`
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&overrideFiles, "override", nil, "override env file (repeatable)")
	cmd.Flags().StringArrayVar(&overrideNames, "name", nil, "name for each override source (repeatable)")
	cmd.Flags().BoolVar(&showSummary, "summary", false, "print a comment summary of applied overrides")

	rootCmd.AddCommand(cmd)
}
