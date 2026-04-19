package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envchain"
	"envoy-sync/internal/envfile"
)

func init() {
	var stopOnFirst bool
	var showOrigin bool

	cmd := &cobra.Command{
		Use:   "envchain [file1] [file2] ...",
		Short: "Chain multiple .env files, resolving keys in order",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			links := make([]envchain.Link, 0, len(args))
			for _, path := range args {
				env, err := envfile.Parse(path)
				if err != nil {
					return fmt.Errorf("parse %s: %w", path, err)
				}
				links = append(links, envchain.Link{Name: path, Env: env})
			}

			result := envchain.Chain(links, stopOnFirst)

			if showOrigin {
				fmt.Fprintln(os.Stdout, "# Key origins:")
				for _, line := range envchain.Summary(result) {
					fmt.Fprintln(os.Stdout, "# "+line)
				}
			}

			for k, v := range result.Env {
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&stopOnFirst, "first", false, "Stop at first definition of each key (first-wins)")
	cmd.Flags().BoolVar(&showOrigin, "origin", false, "Show which file each key came from")
	rootCmd.AddCommand(cmd)
}
