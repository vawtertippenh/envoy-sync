package cmd	"fmt"
	"os"
	"strings"

	"envoy-sync/internal/envclean"
	"envoy-sync/internal/envfile"

	"github.com/spn
func init() {
	var removeEmpty bool
	var trimWhitespace bool
	var dedupe bool

	cmd := &cobra.Command{
		Use:   "envclean [file]",
		Short: "Clean an .env file by removing empty keys, trimming whitespace, and deduplicating",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			opts := envclean.Options{
				RemoveEmpty:     removeEmpty,
				TrimWhitespace:  trimWhitespace,
				DeduplicateKeys: dedupe,
			}

			r := envclean.Clean(env, opts)

			var sb strings.Builder
			for k, v := range r.Env {
				fmt.Fprintf(&sb, "%s=%s\n", k, v)
			}
			fmt.Print(sb.String())

			if len(r.Removed) > 0 {
				fmt.Fprintf(os.Stderr, "removed %d key(s): %s\n", len(r.Removed), strings.Join(r.Removed, ", "))
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&removeEmpty, "remove-empty", false, "Remove keys with empty values")
	cmd.Flags().BoolVar(&trimWhitespace, "trim", false, "Trim leading/trailing whitespace from values")
	cmd.Flags().BoolVar(&dedupe, "dedupe", false, "Remove duplicate keys (keep first occurrence)")

	rootCmd.AddCommand(cmd)
}
