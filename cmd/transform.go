package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/transform"
)

func init() {
	var (
		inputFile string
		ruleFlags []string
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "transform",
		Short: "Apply value transformations (upper, lower, trim, replace) to env keys",
		Example: `  envoy-sync transform -f .env --rule "LOG_LEVEL:upper" --rule "DB_HOST:trim"
  envoy-sync transform -f .env --rule "*:trim" --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(inputFile)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			rules, err := parseRuleFlags(ruleFlags)
			if err != nil {
				return err
			}

			res, err := transform.Transform(env, rules)
			if err != nil {
				return err
			}

			if outputJSON {
				return json.NewEncoder(os.Stdout).Encode(res.Env)
			}

			for _, k := range sortedMapKeys(res.Env) {
				fmt.Printf("%s=%s\n", k, res.Env[k])
			}

			if len(res.Changed) > 0 {
				fmt.Fprintf(os.Stderr, "transformed keys: %s\n", strings.Join(res.Changed, ", "))
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputFile, "file", "f", ".env", "Input .env file")
	cmd.Flags().StringArrayVar(&ruleFlags, "rule", nil, "Rule in format KEY:OP or KEY:replace:FROM:TO")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")

	rootCmd.AddCommand(cmd)
}

// parseRuleFlags converts "KEY:OP" or "KEY:replace:FROM:TO" strings into Rule structs.
func parseRuleFlags(flags []string) ([]transform.Rule, error) {
	rules := make([]transform.Rule, 0, len(flags))
	for _, f := range flags {
		parts := strings.SplitN(f, ":", 4)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid rule %q: expected KEY:OP", f)
		}
		r := transform.Rule{Key: parts[0], Op: parts[1]}
		if r.Op == "replace" {
			if len(parts) < 4 {
				return nil, fmt.Errorf("replace rule %q requires KEY:replace:FROM:TO", f)
			}
			r.From = parts[2]
			r.To = parts[3]
		}
		rules = append(rules, r)
	}
	return rules, nil
}
