package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/envoy-sync/internal/envfile"
	"github.com/yourusername/envoy-sync/internal/envpin"
)

func init() {
	var pinFile string
	var keys string

	cmd := &cobra.Command{
		Use:   "envpin",
		Short: "Pin env values and detect drift",
	}

	saveCmd := &cobra.Command{
		Use:   "save <env-file>",
		Short: "Save pinned values from an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return err
			}
			var keyList []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					keyList = append(keyList, strings.TrimSpace(k))
				}
			} else {
				for k := range env {
					keyList = append(keyList, k)
				}
			}
			if err := envpin.SavePins(env, keyList, pinFile); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Pinned %d key(s) to %s\n", len(keyList), pinFile)
			return nil
		},
	}

	checkCmd := &cobra.Command{
		Use:   "check <env-file>",
		Short: "Check env file against pinned values",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return err
			}
			pf, err := envpin.LoadPins(pinFile)
			if err != nil {
				return err
			}
			results := envpin.Check(pf, env)
			for _, r := range results {
				switch {
				case r.Missing:
					fmt.Fprintf(os.Stdout, "MISSING  %s (pinned=%s)\n", r.Key, r.Pinned)
				case r.Drifted:
					fmt.Fprintf(os.Stdout, "DRIFTED  %s (pinned=%s, actual=%s)\n", r.Key, r.Pinned, r.Actual)
				default:
					fmt.Fprintf(os.Stdout, "OK       %s\n", r.Key)
				}
			}
			if envpin.HasDrift(results) {
				return fmt.Errorf("drift detected")
			}
			return nil
		},
	}

	saveCmd.Flags().StringVarP(&pinFile, "pin-file", "p", "pins.json", "Path to pin file")
	saveCmd.Flags().StringVarP(&keys, "keys", "k", "", "Comma-separated keys to pin (default: all)")
	checkCmd.Flags().StringVarP(&pinFile, "pin-file", "p", "pins.json", "Path to pin file")

	cmd.AddCommand(saveCmd, checkCmd)
	rootCmd.AddCommand(cmd)
}
