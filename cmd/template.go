package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-sync/internal/envfile"
	"github.com/user/envoy-sync/internal/template"
)

func init() {
	var valuesFile string

	cmd := &cobra.Command{
		Use:   "template <template-file>",
		Short: "Fill a .env template with values from another file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tmplPath := args[0]

			tmplEnv, err := envfile.Parse(tmplPath)
			if err != nil {
				return fmt.Errorf("parse template file: %w", err)
			}

			var valEnv map[string]string
			if valuesFile != "" {
				valEnv, err = envfile.Parse(valuesFile)
				if err != nil {
					return fmt.Errorf("parse values file: %w", err)
				}
			} else {
				valEnv = make(map[string]string)
			}

			r := template.Fill(tmplEnv, valEnv)

			for k, v := range r.Filled {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
			}

			summary := template.Summary(r)
			fmt.Fprintln(os.Stderr, summary)

			if len(r.Missing) > 0 {
				return fmt.Errorf("%d required key(s) not filled", len(r.Missing))
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&valuesFile, "values", "v", "", "path to .env file providing values")
	rootCmd.AddCommand(cmd)
}
