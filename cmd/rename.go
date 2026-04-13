package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/rename"
)

func init() {
	var overwrite bool

	var renameCmd = &cobra.Command{
		Use:   "rename <file> <OLD_KEY> <NEW_KEY>",
		Short: "Rename a key in a .env file",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			oldKey := args[1]
			newKey := args[2]

			env, err := envfile.Parse(filePath)
			if err != nil {
				return fmt.Errorf("parsing %s: %w", filePath, err)
			}

			opts := rename.Options{Overwrite: overwrite}
			updated, result := rename.Rename(env, oldKey, newKey, opts)

			if result.Skipped {
				fmt.Fprintf(os.Stderr, "skipped: %s\n", result.Reason)
				return nil
			}

			if err := writeEnvOutput(filePath, updated); err != nil {
				return err
			}

			fmt.Printf("renamed %s -> %s\n", result.OldKey, result.NewKey)
			return nil
		},
	}

	renameCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite new key if it already exists")
ootCmd.AddCommand(renameCmd)
}

func writeEnvOutput(path string, env map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("opening %s for write: %w", path, err)
	}
	defer f.Close()
	for k, v := range env {
		fmt.Fprintf(f, "%s=%s\n", k, v)
	}
	return nil
}
