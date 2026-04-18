package cmd

import (
	"fmt"
	"os"
	"strings"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envtag"

	"github.com/spf13/cobra"
)

func init() {
	var tagFlags []string
	var defaultTag string

	cmd := &cobra.Command{
		Use:   "envtag [file]",
		Short: "Tag env keys by pattern rules",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			tags := map[string][]string{}
			for _, tf := range tagFlags {
				parts := strings.SplitN(tf, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid tag flag %q, expected name=pattern", tf)
				}
				tags[parts[0]] = append(tags[parts[0]], parts[1])
			}

			opts := envtag.Options{
				Tags:       tags,
				DefaultTag: defaultTag,
			}

			results := envtag.Tag(env, opts)
			fmt.Fprint(os.Stdout, envtag.Render(results))
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&tagFlags, "tag", "t", nil, "Tag rule in name=pattern form (repeatable)")
	cmd.Flags().StringVar(&defaultTag, "default", "", "Default tag for unmatched keys")
	rootCmd.AddCommand(cmd)
}
