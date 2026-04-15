package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/split"
)

func init() {
	var (
		prefixFlags []string
		stripPrefix bool
		remainder   string
		outputJSON  bool
	)

	cmd := &cobra.Command{
		Use:   "split <file>",
		Short: "Split an env file into named groups by key prefix",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			prefixes := make(map[string]string, len(prefixFlags))
			for _, p := range prefixFlags {
				parts := strings.SplitN(p, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid prefix mapping %q: expected name=PREFIX_", p)
				}
				prefixes[parts[0]] = parts[1]
			}

			groups, err := split.Split(env, split.Options{
				Prefixes:    prefixes,
				StripPrefix: stripPrefix,
				Remainder:   remainder,
			})
			if err != nil {
				return err
			}

			if outputJSON {
				out := make(map[string]map[string]string, len(groups))
				for _, g := range groups {
					out[g.Name] = g.Env
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(out)
			}

			for _, g := range groups {
				fmt.Printf("[%s] (%d keys)\n", g.Name, len(g.Env))
				for k, v := range g.Env {
					fmt.Printf("  %s=%s\n", k, v)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&prefixFlags, "prefix", "p", nil, "Group prefix mapping: name=PREFIX_ (repeatable)")
	cmd.Flags().BoolVar(&stripPrefix, "strip", false, "Strip matched prefix from output keys")
	cmd.Flags().StringVar(&remainder, "remainder", "", "Group name for unmatched keys (omitted if empty)")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON object")
	_ = cmd.MarkFlagRequired("prefix")

	rootCmd.AddCommand(cmd)
}
