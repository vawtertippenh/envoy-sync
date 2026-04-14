// Package generate provides functionality for generating .env template files
// from existing env maps, replacing values with placeholder descriptions.
package generate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envoy-sync/internal/mask"
)

// Options controls how the template is generated.
type Options struct {
	// MaskSensitive replaces sensitive values with a masked placeholder.
	MaskSensitive bool
	// PlaceholderFormat is the format string used for non-sensitive placeholders.
	// Use %s to embed the key name. Defaults to "<YOUR_%s>".
	PlaceholderFormat string
	// ExtraPatterns are additional key patterns to treat as sensitive.
	ExtraPatterns []string
}

// Generate takes an env map and returns a new map suitable for use as a
// template: sensitive values are masked and all other values are replaced
// with descriptive placeholders.
func Generate(env map[string]string, opts Options) map[string]string {
	format := opts.PlaceholderFormat
	if format == "" {
		format = "<YOUR_%s>"
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		if mask.IsSensitive(k, opts.ExtraPatterns) && opts.MaskSensitive {
			result[k] = mask.MaskValue(k, v, opts.ExtraPatterns)
		} else {
			result[k] = fmt.Sprintf(format, strings.ToUpper(k))
		}
	}
	return result
}

// Render serialises a generated template map to dotenv format with optional
// section comments grouping sensitive and non-sensitive keys.
func Render(template map[string]string, extraPatterns []string) string {
	var sensitive, regular []string
	for k := range template {
		if mask.IsSensitive(k, extraPatterns) {
			sensitive = append(sensitive, k)
		} else {
			regular = append(regular, k)
		}
	}
	sort.Strings(sensitive)
	sort.Strings(regular)

	var sb strings.Builder
	if len(regular) > 0 {
		sb.WriteString("# Application settings\n")
		for _, k := range regular {
			fmt.Fprintf(&sb, "%s=%s\n", k, template[k])
		}
	}
	if len(sensitive) > 0 {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString("# Sensitive / secret settings\n")
		for _, k := range sensitive {
			fmt.Fprintf(&sb, "%s=%s\n", k, template[k])
		}
	}
	return sb.String()
}
