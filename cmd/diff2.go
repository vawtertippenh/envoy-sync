package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/diff2"
	"envoy-sync/internal/envfile"
)

func init() {
	var showUnchanged bool
	var maskValues bool

	cmd := &cobra.Command{
		Use:   "diff2 <file-a> <file-b>",
		Short: "Show a detailed line-by-line diff between two .env files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("reading %s: %w", args[0], err)
			}
			b, err := envfile.Parse(args[1])
			if err != nil {
				return fmt.Errorf("reading %s: %w", args[1], err)
			}

			r := diff2.Diff(a, b)
			output := diff2.Render(r, diff2.RenderOptions{
				ShowUnchanged: showUnchanged,
				MaskValues:   maskValues,
			})
			if output != "" {
				fmt.Print(output)
			}
			fmt.Fprintln(os.Stderr, diff2.Summary(r))
			if r.HasChanges() {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&showUnchanged, "show-unchanged", false, "Include unchanged keys in output")
	cmd.Flags().BoolVar(&maskValues, "mask", false, "Mask all values in output")
	rootCmd.AddCommand(cmd)
}
