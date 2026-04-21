package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envreplace"
)

func init() {
	var filePath string
	var rules []string
	var regexMode bool
	var outputPath string

	cmd := &cobra.Command{
		Use:   "envreplace",
		Short: "Find and replace values in an env file",
		Long:  `Apply one or more find-and-replace rules to values in a .env file. Keys are never modified.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if filePath == "" {
				return fmt.Errorf("--file is required")
			}
			if len(rules) == 0 {
				return fmt.Errorf("at least one --rule is required (format: FIND=REPLACE)")
			}

			env, err := envfile.Parse(filePath)
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			parsedRules := make([]envreplace.Rule, 0, len(rules))
			for _, r := range rules {
				parts := strings.SplitN(r, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid rule %q: expected FIND=REPLACE", r)
				}
				parsedRules = append(parsedRules, envreplace.Rule{
					Find:    parts[0],
					Replace: parts[1],
					Regex:   regexMode,
				})
			}

			res, err := envreplace.Replace(env, parsedRules)
			if err != nil {
				return fmt.Errorf("replace error: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Changed keys (%d): %s\n",
				len(res.ChangedKeys), strings.Join(res.ChangedKeys, ", "))

			dest := outputPath
			if dest == "" {
				dest = filePath
			}

			f, err := os.Create(dest)
			if err != nil {
				return fmt.Errorf("write error: %w", err)
			}
			defer f.Close()

			for _, k := range sortedEnvKeys(res.Env) {
				fmt.Fprintf(f, "%s=%s\n", k, res.Env[k])
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to .env file")
	cmd.Flags().StringArrayVarP(&rules, "rule", "r", nil, "Find=Replace rule (repeatable)")
	cmd.Flags().BoolVar(&regexMode, "regex", false, "Treat Find as a regular expression")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path (defaults to input file)")

	rootCmd.AddCommand(cmd)
}

func sortedEnvKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
