package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/export"
	"envoy-sync/internal/mask"
)

var (
	exportFormat string
	exportMask   bool
	exportSorted bool
)

func init() {
	exportCmd := &cobra.Command{
		Use:   "export [file]",
		Short: "Export an .env file in a specified format",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			env, err := envfile.Parse(path)
			if err != nil {
				return fmt.Errorf("parse %s: %w", path, err)
			}

			if exportMask {
				env = mask.MaskMap(env, nil)
			}

			opts := export.Options{
				Format:  export.Format(exportFormat),
				Masked:  exportMask,
				Sorted:  exportSorted,
			}

			out, err := export.Export(env, opts)
			if err != nil {
				return fmt.Errorf("export: %w", err)
			}

			fmt.Fprint(os.Stdout, out)
			return nil
		},
	}

	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "dotenv", "Output format: dotenv, json, shell")
	exportCmd.Flags().BoolVar(&exportMask, "mask", false, "Mask sensitive values before exporting")
	exportCmd.Flags().BoolVar(&exportSorted, "sorted", true, "Sort keys alphabetically")

	rootCmd.AddCommand(exportCmd)
}
