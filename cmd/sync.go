package cmd

import (
	"fmt"

	"github.com/user/envoy-sync/internal/envfile"
	envsync "github.com/user/envoy-sync/internal/sync"
	"github.com/spf13/cobra"
)

var overwrite bool

var syncCmd = &cobra.Command{
	Use:   "sync <src.env> <dst.env>",
	Short: "Sync variables from a source .env file into a destination .env file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath := args[0]
		dstPath := args[1]

		src, err := envfile.Parse(srcPath)
		if err != nil {
			return fmt.Errorf("parsing source file %q: %w", srcPath, err)
		}

		dst, err := envfile.Parse(dstPath)
		if err != nil {
			return fmt.Errorf("parsing destination file %q: %w", dstPath, err)
		}

		mode := envsync.ModeAddMissing
		if overwrite {
			mode = envsync.ModeOverwrite
		}

		result, err := envsync.Sync(src, dst, dstPath, mode)
		if err != nil {
			return fmt.Errorf("syncing: %w", err)
		}

		if len(result.Added) == 0 && len(result.Updated) == 0 {
			fmt.Println("Nothing to sync — destination is already up to date.")
			return nil
		}

		for _, k := range result.Added {
			fmt.Printf("  + added:   %s\n", k)
		}
		for _, k := range result.Updated {
			fmt.Printf("  ~ updated: %s\n", k)
		}
		fmt.Printf("\nSync complete: %d added, %d updated, %d skipped.\n",
			len(result.Added), len(result.Updated), len(result.Skipped))

		return nil
	},
}

func init() {
	syncCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false,
		"overwrite existing keys in destination with values from source")
	rootCmd.AddCommand(syncCmd)
}
