package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envgroup"
)

func init() {
	var prefixFlags []string
	var remainder string

	cmd := &cobra.Command{
		Use:   "envgroup [file]",
		Short: "Group env keys by prefix",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			prefixes := make(map[string]string)
			for _, p := range prefixFlags {
				parts := strings.SplitN(p, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid prefix flag %q, expected name=PREFIX_", p)
				}
				prefixes[parts[0]] = parts[1]
			}

			groups := envgroup.GroupBy(env, envgroup.Options{
				Prefixes:  prefixes,
				Remainder: remainder,
			})

			for _, g := range groups {
				fmt.Printf("[%s] (%d keys)\n", g.Name, len(g.Keys))
				for _, k := range g.Keys {
					fmt.Printf("  %s=%s\n", k, g.Env[k])
				}
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&prefixFlags, "prefix", "p", nil, "Group prefix mapping name=PREFIX_ (repeatable)")
	cmd.Flags().StringVar(&remainder, "remainder", "", "Group name for unmatched keys (omit to drop them)")
	rootCmd.AddCommand(cmd)
}
