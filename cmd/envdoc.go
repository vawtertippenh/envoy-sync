package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"envoy-sync/internal/envdoc"
	"envoy-sync/internal/envfile"
	"github.com/spf13/cobra"
)

func init() {
	var (
		format        string
		descriptions  []string
		requiredKeys  []string
		sensitiveKeys []string
	)

	cmd := &cobra.Command{
		Use:   "envdoc <file>",
		Short: "Generate documentation from a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			descMap := make(map[string]string)
			for _, d := range descriptions {
				parts := strings.SplitN(d, "=", 2)
				if len(parts) == 2 {
					descMap[parts[0]] = parts[1]
				}
			}

			opts := envdoc.Options{
				Descriptions:  descMap,
				RequiredKeys:  requiredKeys,
				SensitiveKeys: sensitiveKeys,
			}
			result := envdoc.Document(env, opts)

			switch format {
			case "json":
				return json.NewEncoder(os.Stdout).Encode(result.Entries)
			default:
				fmt.Print(envdoc.Render(result))
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "markdown", "Output format: markdown, json")
	cmd.Flags().StringArrayVar(&descriptions, "desc", nil, "Key descriptions as KEY=description")
	cmd.Flags().StringSliceVar(&requiredKeys, "required", nil, "Comma-separated required keys")
	cmd.Flags().StringSliceVar(&sensitiveKeys, "sensitive", nil, "Comma-separated sensitive keys")

	rootCmd.AddCommand(cmd)
}
