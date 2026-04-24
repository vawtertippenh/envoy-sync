package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envbatch"
	"envoy-sync/internal/envfile"
)

func init() {
	var (
		batchSize int
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "envbatch [file]",
		Short: "Split an .env file into fixed-size batches",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			batches, err := envbatch.BatchEnv(env, envbatch.Options{
				Size:     batchSize,
				SortKeys: true,
			})
			if err != nil {
				return fmt.Errorf("batch error: %w", err)
			}

			if outputJSON {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(batches)
			}

			for _, b := range batches {
				fmt.Printf("[%s] (%d keys)\n", b.Name, len(b.Items))
				for k, v := range b.Items {
					fmt.Printf("  %s=%s\n", k, v)
				}
			}
			fmt.Println(envbatch.Summary(batches))
			return nil
		},
	}

	cmd.Flags().IntVarP(&batchSize, "size", "s", 10, "number of keys per batch")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "output batches as JSON")
	rootCmd.AddCommand(cmd)
}
