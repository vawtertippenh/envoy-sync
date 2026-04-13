package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-sync/internal/envfile"
	"github.com/user/envoy-sync/internal/snapshot"
)

func init() {
	var label string
	var outputPath string
	var comparePath string

	snapshotCmd := &cobra.Command{
		Use:   "snapshot <envfile>",
		Short: "Capture or compare a snapshot of an .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			if comparePath != "" {
				old, err := snapshot.Load(comparePath)
				if err != nil {
					return fmt.Errorf("load snapshot: %w", err)
				}
				current := snapshot.Take(label, env)
				d := snapshot.Compare(old, current)
				if !snapshot.HasDrift(d) {
					fmt.Println("No drift detected.")
					return nil
				}
				for k, v := range d.Added {
					fmt.Fprintf(os.Stdout, "+ %s=%s\n", k, v)
				}
				for k := range d.Removed {
					fmt.Fprintf(os.Stdout, "- %s\n", k)
				}
				for k, ch := range d.Changed {
					fmt.Fprintf(os.Stdout, "~ %s: %s -> %s\n", k, ch.Before, ch.After)
				}
				return nil
			}

			s := snapshot.Take(label, env)
			if outputPath == "" {
				return fmt.Errorf("--output is required when taking a snapshot")
			}
			if err := snapshot.Save(s, outputPath); err != nil {
				return err
			}
			fmt.Printf("Snapshot '%s' saved to %s\n", label, outputPath)
			return nil
		},
	}

	snapshotCmd.Flags().StringVarP(&label, "label", "l", "snapshot", "Label for the snapshot")
	snapshotCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to write snapshot JSON")
	snapshotCmd.Flags().StringVarP(&comparePath, "compare", "c", "", "Path to an existing snapshot to diff against")

	rootCmd.AddCommand(snapshotCmd)
}
