package cmd

import (
	"fmt"
	"os"
	"strings"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/reorder"

	"github.com/spf13/cobra"
)

func init() {
	var strategy string
	var order []string
	var unknownLast bool

	cmd := &cobra.Command{
		Use:   "reorder <file>",
		Short: "Reorder keys in a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			opts := reorder.Options{
				Strategy:       reorder.Strategy(strategy),
				Order:          order,
				PutUnknownLast: unknownLast,
			}

			res, err := reorder.Reorder(env, opts)
			if err != nil {
				return fmt.Errorf("reorder error: %w", err)
			}

			w := os.Stdout
			for _, k := range res.Ordered {
				v := res.Env[k]
				if strings.ContainsAny(v, " \t") {
					v = `"` + v + `"`
				}
				fmt.Fprintf(w, "%s=%s\n", k, v)
			}

			if len(res.Unknown) > 0 && !unknownLast {
				fmt.Fprintf(os.Stderr, "note: %d key(s) not in order list: %s\n",
					len(res.Unknown), strings.Join(res.Unknown, ", "))
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&strategy, "strategy", "s", "alpha", "reorder strategy: alpha, template, custom")
	cmd.Flags().StringSliceVarP(&order, "order", "o", nil, "comma-separated key order (template/custom)")
	cmd.Flags().BoolVar(&unknownLast, "unknown-last", true, "append unknown keys at end (template/custom)")

	rootCmd.AddCommand(cmd)
}
