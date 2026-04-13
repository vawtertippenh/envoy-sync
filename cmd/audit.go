package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/audit"
	"envoy-sync/internal/diff"
	"envoy-sync/internal/envfile"
	"envoy-sync/internal/mask"
)

func init() {
	var maskKeys []string

	auditCmd := &cobra.Command{
		Use:   "audit [fileA] [fileB]",
		Short: "Audit differences between two .env files and log sensitive key access",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			envA, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parsing %s: %w", args[0], err)
			}
			envB, err := envfile.Parse(args[1])
			if err != nil {
				return fmt.Errorf("parsing %s: %w", args[1], err)
			}

			log := &audit.Log{}

			results := diff.Compare(envA, envB, maskKeys)
			for _, r := range results {
				switch r.Status {
				case diff.StatusAdded:
					log.Record(r.Key, audit.KindAdded, fmt.Sprintf("value: %s", r.ValueB))
				case diff.StatusRemoved:
					log.Record(r.Key, audit.KindRemoved, fmt.Sprintf("was: %s", r.ValueA))
				case diff.StatusChanged:
					log.Record(r.Key, audit.KindChanged, fmt.Sprintf("%s -> %s", r.ValueA, r.ValueB))
				}
				if mask.IsSensitive(r.Key, maskKeys) {
					log.Record(r.Key, audit.KindMasked, "sensitive key accessed during audit")
				}
			}

			fmt.Fprintln(os.Stdout, log.Summary())
			for _, e := range log.Entries {
				fmt.Fprintf(os.Stdout, "  [%s] %-20s %s\n",
					e.Kind, e.Key, e.Detail)
			}
			return nil
		},
	}

	auditCmd.Flags().StringSliceVar(&maskKeys, "mask-keys", nil,
		"Additional key patterns to treat as sensitive")

	rootCmd.AddCommand(auditCmd)
}
