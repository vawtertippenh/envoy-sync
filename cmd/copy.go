package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	envcopy "envoy-sync/internal/copy"
	"envoy-sync/internal/envfile"
)

func init() {
	var keys, prefix string
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "copy <src> <dst>",
		Short: "Copy keys from one .env file into another",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parsing src: %w", err)
			}
			dst, err := envfile.Parse(args[1])
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("parsing dst: %w", err)
			}
			if dst == nil {
				dst = map[string]string{}
			}

			var keyList []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					if k = strings.TrimSpace(k); k != "" {
						keyList = append(keyList, k)
					}
				}
			}

			out, res := envcopy.Copy(src, dst, envcopy.Options{
				Keys:      keyList,
				Overwrite: overwrite,
				Prefix:    prefix,
			})

			f, err := os.Create(args[1])
			if err != nil {
				return err
			}
			defer f.Close()
			for k, v := range out {
				fmt.Fprintf(f, "%s=%s\n", k, v)
			}

			fmt.Fprintf(os.Stderr, "copied: %d  skipped: %d\n", len(res.Copied), len(res.Skipped))
			return nil
		},
	}

	cmd.Flags().StringVar(&keys, "keys", "", "comma-separated keys to copy (default: all)")
	cmd.Flags().StringVar(&prefix, "prefix", "", "prefix to add to copied keys")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys in dst")
	rootCmd.AddCommand(cmd)
}
