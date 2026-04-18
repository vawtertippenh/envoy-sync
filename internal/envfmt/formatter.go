package envfmt

import (
	"fmt"
	"sort"
	"strings"
)

// Style controls how the .env file is formatted.
type Style struct {
	UppercaseKeys   bool
	SortKeys        bool
	QuoteAllValues  bool
	SpaceAroundEqual bool
}

// Result holds the formatted output and a summary.
type Result struct {
	Lines    []string
	Changed  int
	Total    int
}

// Format applies the given Style to env and returns a Result.
func Format(env map[string]string, s Style) Result {
	keys := sortedKeys(env)
	if s.SortKeys {
		sort.Strings(keys)
	}

	var lines []string
	changed := 0

	for _, k := range keys {
		v := env[k]
		newKey := k
		if s.UppercaseKeys {
			newKey = strings.ToUpper(k)
		}

		newVal := v
		if s.QuoteAllValues && !isQuoted(v) {
			newVal = fmt.Sprintf("%q", v)
		}

		var line string
		if s.SpaceAroundEqual {
			line = fmt.Sprintf("%s = %s", newKey, newVal)
		} else {
			line = fmt.Sprintf("%s=%s", newKey, newVal)
		}

		if newKey != k || newVal != v {
			changed++
		}
		lines = append(lines, line)
	}

	return Result{Lines: lines, Changed: changed, Total: len(keys)}
}

// Render returns the formatted env as a single string.
func Render(r Result) string {
	return strings.Join(r.Lines, "\n") + "\n"
}

func isQuoted(s string) bool {
	return (strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) ||
		(strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'"))
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
