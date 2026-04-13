package cmd

import (
	"fmt"
	"os"
	"strings"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/promote"

	"github.com/spf13/cobra"
)

func init() {
	var (
		keys      string
		overwrite bool
		dryRun    bool
	)

	cmd := &cobra.Command{
		Use:   "promote <src> <dst>",
		Short: "Promote env variables from one environment file to another",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcPath, dstPath := args[0], args[1]

			src, err := envfile.Parse(srcPath)
			if err != nil {
				return fmt.Errorf("parse source: %w", err)
			}
			dst, err := envfile.Parse(dstPath)
			if err != nil {
				return fmt.Errorf("parse destination: %w", err)
			}

			var keyList []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					if k = strings.TrimSpace(k); k != "" {
						keyList = append(keyList, k)
					}
				}
			}

			opts := promote.Options{
				Keys:      keyList,
				Overwrite: overwrite,
				DryRun:    dryRun,
			}

			out, res, err := promote.Promote(src, dst, opts)
			if err != nil {
				return err
			}

			fmt.Fprintln(os.Stderr, res.Summary())

			if dryRun {
				fmt.Fprintln(os.Stderr, "[dry-run] no changes written")
				return nil
			}

			f, err := os.Create(dstPath)
			if err != nil {
				return fmt.Errorf("open destination for writing: %w", err)
			}
			defer f.Close()
			for k, v := range out {
				fmt.Fprintf(f, "%s=%s\n", k, v)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&keys, "keys", "", "Comma-separated list of keys to promote (default: all)")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing keys in destination")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing to destination")

	rootCmd.AddCommand(cmd)
}
