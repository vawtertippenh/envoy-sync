package cmd

import (
	"fmt"
	"os"
	"sort"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/mask"

	"github.com/spf13/cobra"
)

var maskCmd = &cobra.Command{
	Use:   "mask [file]",
	Short: "Print env file with sensitive values masked",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		env, err := envfile.Parse(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse env file: %w", err)
		}

		extraPatterns, _ := cmd.Flags().GetStringSlice("extra")
		masked := mask.MaskMap(env, extraPatterns)

		keys := make([]string, 0, len(masked))
		for k := range masked {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, masked[k])
		}

		return nil
	},
}

func init() {
	maskCmd.Flags().StringSlice("extra", []string{}, "Additional key patterns to treat as sensitive")
	rootCmd.AddCommand(maskCmd)
}
