// Package envstats provides statistical analysis of env file contents.
package envstats

import (
	"fmt"
	"sort"
	"strings"
)

// Stats holds aggregated statistics about an env map.
type Stats struct {
	Total       int
	Empty       int
	NonEmpty    int
	Sensitive   int
	AvgLength   float64
	MaxLength   int
	MinLength   int
	Prefixes    map[string]int
}

var defaultSensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE", "AUTH",
}

// Analyze computes statistics for the given env map.
// sensitivePatterns extends the default set of sensitive key patterns.
func Analyze(env map[string]string, sensitivePatterns []string) Stats {
	if len(env) == 0 {
		return Stats{MinLength: 0}
	}

	patterns := append(defaultSensitivePatterns, sensitivePatterns...)

	s := Stats{
		Prefixes:  make(map[string]int),
		MinLength: int(^uint(0) >> 1),
	}

	totalLen := 0
	for k, v := range env {
		s.Total++
		if v == "" {
			s.Empty++
		} else {
			s.NonEmpty++
		}
		l := len(v)
		totalLen += l
		if l > s.MaxLength {
			s.MaxLength = l
		}
		if l < s.MinLength {
			s.MinLength = l
		}
		if isSensitive(k, patterns) {
			s.Sensitive++
		}
		if idx := strings.Index(k, "_"); idx > 0 {
			prefix := k[:idx]
			s.Prefixes[prefix]++
		}
	}

	if s.Total > 0 {
		s.AvgLength = float64(totalLen) / float64(s.Total)
	}
	return s
}

// Summary returns a human-readable summary of the stats.
func Summary(s Stats) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Total keys   : %d\n", s.Total)
	fmt.Fprintf(&sb, "Non-empty    : %d\n", s.NonEmpty)
	fmt.Fprintf(&sb, "Empty        : %d\n", s.Empty)
	fmt.Fprintf(&sb, "Sensitive    : %d\n", s.Sensitive)
	fmt.Fprintf(&sb, "Avg value len: %.1f\n", s.AvgLength)
	fmt.Fprintf(&sb, "Max value len: %d\n", s.MaxLength)
	fmt.Fprintf(&sb, "Min value len: %d\n", s.MinLength)
	if len(s.Prefixes) > 0 {
		sb.WriteString("Prefixes:\n")
		keys := make([]string, 0, len(s.Prefixes))
		for k := range s.Prefixes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(&sb, "  %s: %d\n", k, s.Prefixes[k])
		}
	}
	return sb.String()
}

func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
