package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/trim"

	"github.com/spf13/cobra"
)

func init() {
	var (
		allowList   string
		denyList    string
		removeEmpty bool
		outputJSON  bool
	)

	trimCmd := &cobra.Command{
		Use:   "trim [file]",
		Short: "Remove unused or unwanted keys from a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			opts := trim.Options{RemoveEmpty: removeEmpty}
			if allowList != "" {
				opts.AllowList = splitCSV(allowList)
			}
			if denyList != "" {
				opts.DenyList = splitCSV(denyList)
			}

			result := trim.Trim(env, opts)

			if outputJSON {
				return json.NewEncoder(os.Stdout).Encode(result)
			}

			fmt.Fprintf(os.Stdout, "Kept %d key(s), removed %d key(s)\n",
				len(result.Kept), len(result.Removed))
			if len(result.Removed) > 0 {
				fmt.Fprintln(os.Stdout, "Removed:", strings.Join(result.Removed, ", "))
			}
			return nil
		},
	}

	trimCmd.Flags().StringVar(&allowList, "allow", "", "Comma-separated list of keys to keep (all others removed)")
	trimCmd.Flags().StringVar(&denyList, "deny", "", "Comma-separated list of keys to remove")
	trimCmd.Flags().BoolVar(&removeEmpty, "remove-empty", false, "Remove keys with empty values")
	trimCmd.Flags().BoolVar(&outputJSON, "json", false, "Output result as JSON")

	rootCmd.AddCommand(trimCmd)
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
