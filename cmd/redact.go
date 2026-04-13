package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/envoy-sync/internal/envfile"
	"github.com/envoy-sync/internal/redact"
	"github.com/spf13/cobra"
)

func init() {
	var rules []string
	var extraPatterns []string
	var outputJSON bool

	redactCmd := &cobra.Command{
		Use:   "redact [file]",
		Short: "Redact sensitive values from an .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			parsedRules := make([]redact.Rule, 0, len(rules))
			for _, r := range rules {
				parts := strings.SplitN(r, "=", 2)
				rule := redact.Rule{Key: parts[0]}
				if len(parts) == 2 {
					rule.Replacement = parts[1]
				}
				parsedRules = append(parsedRules, rule)
			}

			result := redact.Redact(env, parsedRules, extraPatterns)

			if outputJSON {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(result.Redacted)
			}

			for k, v := range result.Redacted {
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
			}
			fmt.Fprintln(os.Stderr, result.Summary())
			return nil
		},
	}

	redactCmd.Flags().StringArrayVarP(&rules, "rule", "r", nil, "Redaction rules as KEY or KEY=replacement")
	redactCmd.Flags().StringArrayVarP(&extraPatterns, "pattern", "p", nil, "Extra sensitive key patterns")
	redactCmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")

	rootCmd.AddCommand(redactCmd)
}
