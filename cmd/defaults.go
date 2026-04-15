package cmd

import (
	"fmt"
	"os"
	"strings"

	"envoy-sync/internal/defaults"
	"envoy-sync/internal/envfile"

	"github.com/spf13/cobra"
)

func init() {
	var overrideAll bool
	var rulesFlag []string

	cmd := &cobra.Command{
		Use:   "defaults [file]",
		Short: "Apply default values to missing or empty keys in an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			var rules []defaults.Rule
			for _, r := range rulesFlag {
				parts := strings.SplitN(r, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid rule %q: expected KEY=VALUE", r)
				}
				rules = append(rules, defaults.Rule{
					Key:      parts[0],
					Value:    parts[1],
					Override: overrideAll,
				})
			}

			res := defaults.Apply(env, rules)

			for _, k := range res.Applied {
				fmt.Fprintf(os.Stderr, "applied default: %s=%s\n", k, res.Env[k])
			}
			for _, k := range res.Skipped {
				fmt.Fprintf(os.Stderr, "skipped (exists): %s\n", k)
			}

			for k, v := range res.Env {
				fmt.Printf("%s=%s\n", k, v)
			}

			fmt.Fprintln(os.Stderr, res.Summary())
			return nil
		},
	}

	cmd.Flags().BoolVar(&overrideAll, "override", false, "Override existing values with defaults")
	cmd.Flags().StringArrayVarP(&rulesFlag, "rule", "r", nil, "Default rule in KEY=VALUE format (repeatable)")

	rootCmd.AddCommand(cmd)
}
