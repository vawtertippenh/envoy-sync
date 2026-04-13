package export

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Format represents the output format for exported env data.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatShell  Format = "shell"
)

// Options configures the export behavior.
type Options struct {
	Format  Format
	Masked  bool
	Sorted  bool
}

// Export converts an env map to the specified output format.
func Export(env map[string]string, opts Options) (string, error) {
	keys := keys(env, opts.Sorted)

	switch opts.Format {
	case FormatJSON:
		return toJSON(env, keys)
	case FormatShell:
		return toShell(env, keys), nil
	case FormatDotenv, "":
		return toDotenv(env, keys), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", opts.Format)
	}
}

func toDotenv(env map[string]string, keys []string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
	}
	return sb.String()
}

func toJSON(env map[string]string, keys []string) (string, error) {
	ordered := make(map[string]string, len(keys))
	for _, k := range keys {
		ordered[k] = env[k]
	}
	b, err := json.MarshalIndent(ordered, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(b) + "\n", nil
}

func toShell(env map[string]string, keys []string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "export %s=%q\n", k, env[k])
	}
	return sb.String()
}

func keys(env map[string]string, sorted bool) []string {
	out := make([]string, 0, len(env))
	for k := range env {
		out = append(out, k)
	}
	if sorted {
		sort.Strings(out)
	}
	return out
}
