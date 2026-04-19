package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envsplit"
	"github.com/spf13/cobra"
)

func init() {
	var prefixes []string
	var strip bool
	var remainder bool
	var format string

	cmd := &cobra.Command{
		Use:   "envsplit <file>",
		Short: "Split a .env file into groups by key prefix",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return err
			}
			parts, err := envsplit.Split(env, envsplit.Options{
				Prefixes:      prefixes,
				StripPrefix:   strip,
				KeepRemainder: remainder,
			})
			if err != nil {
				return err
			}
			switch format {
			case "json":
				out := make(map[string]map[string]string)
				for _, p := range parts {
					out[p.Name] = p.Env
				}
				return json.NewEncoder(os.Stdout).Encode(out)
			default:
				for _, p := range parts {
					fmt.Fprintf(os.Stdout, "# --- %s ---\n", p.Name)
					for k, v := range p.Env {
						if strings.ContainsAny(v, " \t") {
							v = `"` + v + `"`
						}
						fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&prefixes, "prefix", "p", nil, "Prefixes to split on (comma-separated or repeated)")
	cmd.Flags().BoolVar(&strip, "strip", false, "Strip the prefix from keys in each group")
	cmd.Flags().BoolVar(&remainder, "remainder", false, "Collect unmatched keys into _remainder group")
	cmd.Flags().StringVar(&format, "format", "dotenv", "Output format: dotenv|json")

	rootCmd.AddCommand(cmd)
}
