package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envfreeze"
)

func init() {
	var (
		existingFile string
		allowKeys    string
		denyKeys     string
		overwrite    bool
		outputJSON   bool
	)

	cmd := &cobra.Command{
		Use:   "envfreeze [file]",
		Short: "Freeze env values to lock them from future changes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			var existing map[string]string
			if existingFile != "" {
				existing, err = envfile.Parse(existingFile)
				if err != nil {
					return fmt.Errorf("parse existing: %w", err)
				}
			}

			opts := envfreeze.Options{
				OverwriteExisting: overwrite,
			}
			if allowKeys != "" {
				opts.AllowKeys = strings.Split(allowKeys, ",")
			}
			if denyKeys != "" {
				opts.DenyKeys = strings.Split(denyKeys, ",")
			}

			res, err := envfreeze.Freeze(src, existing, opts)
			if err != nil {
				return err
			}

			if outputJSON {
				return json.NewEncoder(os.Stdout).Encode(res.Frozen)
			}

			for _, k := range sortedStringKeys(res.Frozen) {
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, res.Frozen[k])
			}
			if len(res.Skipped) > 0 {
				fmt.Fprintf(os.Stderr, "skipped (empty): %s\n", strings.Join(res.Skipped, ", "))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&existingFile, "existing", "", "Path to existing frozen env file")
	cmd.Flags().StringVar(&allowKeys, "allow", "", "Comma-separated keys to freeze (allowlist)")
	cmd.Flags().StringVar(&denyKeys, "deny", "", "Comma-separated keys to exclude")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing frozen values")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")

	rootCmd.AddCommand(cmd)
}

func sortedStringKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	import_sort_strings(keys)
	return keys
}
