package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/yourusername/envoy-sync/internal/envfile"
	"github.com/yourusername/envoy-sync/internal/envsanitize"
	"github.com/spf13/cobra"
)

func init() {
	var (
		trimWS       bool
		stripCtrl    bool
		replaceNL    bool
		maxLen       int
		outputFile   string
	)

	cmd := &cobra.Command{
		Use:   "envsanitize <file>",
		Short: "Sanitize env file values (trim whitespace, strip control chars, etc.)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			opts := envсаnitize.Options{
				TrimWhitespace:    trimWS,
				StripControlChars: stripCtrl,
				ReplaceNewlines:   replaceNL,
				MaxLength:         maxLen,
			}

			res := envсаnitize.Sanitize(env, opts)

			if len(res.Changed) > 0 {
				sort.Strings(res.Changed)
				fmt.Fprintf(os.Stderr, "sanitized %d key(s): %v\n", len(res.Changed), res.Changed)
			}

			w := os.Stdout
			if outputFile != "" {
				f, err := os.Create(outputFile)
				if err != nil {
					return fmt.Errorf("create output: %w", err)
				}
				defer f.Close()
				w = f
			}

			keys := make([]string, 0, len(res.Env))
			for k := range res.Env {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Fprintf(w, "%s=%s\n", k, res.Env[k])
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&trimWS, "trim", false, "trim leading/trailing whitespace from values")
	cmd.Flags().BoolVar(&stripCtrl, "strip-control", false, "strip non-printable control characters")
	cmd.Flags().BoolVar(&replaceNL, "replace-newlines", false, "replace embedded newlines with \\n")
	cmd.Flags().IntVar(&maxLen, "max-length", 0, "truncate values to this length (0 = unlimited)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "write result to file instead of stdout")

	rootCmd.AddCommand(cmd)
}
