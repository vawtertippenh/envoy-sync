package compare

import (
	"fmt"
	"sort"
	"strings"
)

// MatchResult holds the result of comparing two env maps against a template.
type MatchResult struct {
	Missing  []string // keys in template but not in target
	Extra    []string // keys in target but not in template
	Mismatch []string // keys present in both but with type/format differences
}

// Against compares a target env map against a reference template map.
// It identifies missing keys, extra keys, and potential mismatches.
func Against(template, target map[string]string) MatchResult {
	result := MatchResult{}

	for k := range template {
		if _, ok := target[k]; !ok {
			result.Missing = append(result.Missing, k)
		}
	}

	for k := range target {
		if _, ok := template[k]; !ok {
			result.Extra = append(result.Extra, k)
		}
	}

	for k, tv := range template {
		if rv, ok := target[k]; ok {
			if detectTypeMismatch(tv, rv) {
				result.Mismatch = append(result.Mismatch, k)
			}
		}
	}

	sort.Strings(result.Missing)
	sort.Strings(result.Extra)
	sort.Strings(result.Mismatch)
	return result
}

// detectTypeMismatch returns true if the two values appear to be of different
// semantic types (e.g. one is a boolean and the other is not).
func detectTypeMismatch(a, b string) bool {
	boolVals := map[string]bool{"true": true, "false": true, "1": true, "0": true, "yes": true, "no": true}
	aIsBool := boolVals[strings.ToLower(a)]
	bIsBool := boolVals[strings.ToLower(b)]
	return aIsBool != bIsBool
}

// Summary returns a human-readable summary of the MatchResult.
func (r MatchResult) Summary() string {
	if len(r.Missing) == 0 && len(r.Extra) == 0 && len(r.Mismatch) == 0 {
		return "✓ target matches template exactly"
	}
	var sb strings.Builder
	if len(r.Missing) > 0 {
		fmt.Fprintf(&sb, "missing keys (%d): %s\n", len(r.Missing), strings.Join(r.Missing, ", "))
	}
	if len(r.Extra) > 0 {
		fmt.Fprintf(&sb, "extra keys (%d): %s\n", len(r.Extra), strings.Join(r.Extra, ", "))
	}
	if len(r.Mismatch) > 0 {
		fmt.Fprintf(&sb, "type mismatch keys (%d): %s\n", len(r.Mismatch), strings.Join(r.Mismatch, ", "))
	}
	return strings.TrimRight(sb.String(), "\n")
}
