// Package envсаnitize provides functionality to sanitize env maps by
// removing or replacing characters that are invalid or unsafe in env values.
package envsanitize

import (
	"strings"
	"unicode"
)

// Options controls sanitization behaviour.
type Options struct {
	// StripControlChars removes non-printable control characters from values.
	StripControlChars bool
	// TrimWhitespace trims leading and trailing whitespace from values.
	TrimWhitespace bool
	// ReplaceNewlines replaces embedded newlines with a literal \n sequence.
	ReplaceNewlines bool
	// MaxLength truncates values longer than MaxLength (0 = unlimited).
	MaxLength int
}

// Result holds the sanitized env map and a report of changes made.
type Result struct {
	Env     map[string]string
	Changed []string // keys whose values were modified
}

// Sanitize applies the given options to env and returns a sanitized copy.
func Sanitize(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	var changed []string

	for k, v := range env {
		sanitized := sanitizeValue(v, opts)
		out[k] = sanitized
		if sanitized != v {
			changed = append(changed, k)
		}
	}

	return Result{Env: out, Changed: changed}
}

func sanitizeValue(v string, opts Options) string {
	if opts.TrimWhitespace {
		v = strings.TrimSpace(v)
	}

	if opts.ReplaceNewlines {
		v = strings.ReplaceAll(v, "\r\n", `\n`)
		v = strings.ReplaceAll(v, "\n", `\n`)
		v = strings.ReplaceAll(v, "\r", `\n`)
	}

	if opts.StripControlChars {
		var b strings.Builder
		for _, r := range v {
			if r == '\t' || !unicode.IsControl(r) {
				b.WriteRune(r)
			}
		}
		v = b.String()
	}

	if opts.MaxLength > 0 && len(v) > opts.MaxLength {
		v = v[:opts.MaxLength]
	}

	return v
}
