package envquote

import (
	"fmt"
	"sort"
	"strings"
)

// Style controls how values are quoted.
type Style string

const (
	StyleDouble Style = "double"
	StyleSingle Style = "single"
	StyleAuto   Style = "auto"
)

// Options configures the quoting behaviour.
type Options struct {
	Style       Style
	ForceAll    bool // quote every value regardless of content
	StripExisting bool // remove existing quotes before re-quoting
}

// Result holds the quoted environment map and metadata.
type Result struct {
	Env     map[string]string
	Quoted  int
	Skipped int
}

// Quote applies quoting rules to env values.
func Quote(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	var quoted, skipped int

	for k, v := range env {
		val := v
		if opts.StripExisting {
			val = stripQuotes(val)
		}

		if opts.ForceAll || needsQuoting(val) {
			val = applyStyle(val, opts.Style)
			quoted++
		} else {
			skipped++
		}
		out[k] = val
	}

	return Result{Env: out, Quoted: quoted, Skipped: skipped}
}

// Render formats the result as dotenv lines.
func Render(r Result) string {
	keys := sortedKeys(r.Env)
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, r.Env[k])
	}
	return sb.String()
}

func needsQuoting(v string) bool {
	return strings.ContainsAny(v, " \t\n#$'\"\\")
}

func applyStyle(v string, style Style) string {
	v = stripQuotes(v)
	switch style {
	case StyleSingle:
		return "'" + strings.ReplaceAll(v, "'", `'\''`) + "'"
	case StyleAuto:
		if strings.Contains(v, "'") {
			return "\"" + v + "\""
		}
		return "'" + v + "'"
	default: // double
		return "\"" + strings.ReplaceAll(v, `"`, `\"`) + "\""
	}
}

func stripQuotes(v string) string {
	if len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
			return v[1 : len(v)-1]
		}
	}
	return v
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
