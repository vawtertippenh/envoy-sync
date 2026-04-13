// Package clone provides functionality to clone an env file
// into a new target file, optionally masking sensitive values.
package clone

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"envoy-sync/internal/mask"
)

// Options controls clone behaviour.
type Options struct {
	// MaskSensitive replaces sensitive values with a placeholder.
	MaskSensitive bool
	// ExtraPatterns are additional key patterns considered sensitive.
	ExtraPatterns []string
	// Placeholder is the string used when masking (default: "***").
	Placeholder string
}

// Result holds the outcome of a clone operation.
type Result struct {
	Written  int
	Masked   int
	DestPath string
}

// Clone copies src env map to destPath, applying options.
func Clone(src map[string]string, destPath string, opts Options) (Result, error) {
	if destPath == "" {
		return Result{}, fmt.Errorf("destination path must not be empty")
	}

	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "***"
	}

	var lines []string
	keys := sortedKeys(src)

	result := Result{DestPath: destPath}

	for _, k := range keys {
		v := src[k]
		if opts.MaskSensitive && mask.IsSensitive(k, opts.ExtraPatterns) {
			v = placeholder
			result.Masked++
		}
		lines = append(lines, fmt.Sprintf("%s=%s", k, quoteIfNeeded(v)))
		result.Written++
	}

	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(destPath, []byte(content), 0o644); err != nil {
		return Result{}, fmt.Errorf("writing clone to %s: %w", destPath, err)
	}

	return result, nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t#") {
		return fmt.Sprintf("%q", v)
	}
	return v
}
