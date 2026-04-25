package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-sync/internal/envdiff2"
	"github.com/user/envoy-sync/internal/envfile"
)

func init() {
	var maskSensitive bool
	var showUnchanged bool

	cmd := &cobra.Command{
		Use:   "envdiff2 <file-a> <file-b>",
		Short: "Structured diff between two .env files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parsing %s: %w", args[0], err)
			}
			b, err := envfile.Parse(args[1])
			if err != nil {
				return fmt.Errorf("parsing %s: %w", args[1], err)
			}

			result := envdiff2.Diff(a, b, showUnchanged)
			out := envdiff2.Render(result, envdiff2.RenderOptions{
				MaskSensitive: maskSensitive,
				ShowUnchanged: showUnchanged,
			})
			fmt.Fprint(os.Stdout, out)
			fmt.Fprintln(os.Stderr, envdiff2.Summary(result))
			return nil
		},
	}

	cmd.Flags().BoolVar(&maskSensitive, "mask", false, "Mask sensitive values in output")
	cmd.Flags().BoolVar(&showUnchanged, "unchanged", false, "Include unchanged keys in output")
	rootCmd.AddCommand(cmd)
}
