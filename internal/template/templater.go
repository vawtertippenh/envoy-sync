package template

import (
	"fmt"
	"sort"
	"strings"
)

// Result holds the output of a template fill operation.
type Result struct {
	Filled   map[string]string
	Missing  []string
	Unused   []string
}

// Fill takes a template env map (keys with placeholder values like "<required>"
// or "<optional>") and a values map, and returns a filled result.
func Fill(tmpl map[string]string, values map[string]string) Result {
	filled := make(map[string]string, len(tmpl))
	var missing []string

	for k, v := range tmpl {
		if val, ok := values[k]; ok {
			filled[k] = val
		} else if isRequired(v) {
			missing = append(missing, k)
			filled[k] = v
		} else {
			// optional or has default — keep template value
			filled[k] = v
		}
	}

	// detect unused keys in values that are not in the template
	var unused []string
	for k := range values {
		if _, ok := tmpl[k]; !ok {
			unused = append(unused, k)
		}
	}

	sort.Strings(missing)
	sort.Strings(unused)

	return Result{
		Filled:  filled,
		Missing: missing,
		Unused:  unused,
	}
}

// Summary returns a human-readable summary of the fill result.
func Summary(r Result) string {
	var sb strings.Builder
	if len(r.Missing) == 0 && len(r.Unused) == 0 {
		sb.WriteString("template fill complete — no issues\n")
		return sb.String()
	}
	if len(r.Missing) > 0 {
		sb.WriteString(fmt.Sprintf("missing required keys (%d): %s\n",
			len(r.Missing), strings.Join(r.Missing, ", ")))
	}
	if len(r.Unused) > 0 {
		sb.WriteString(fmt.Sprintf("unused value keys (%d): %s\n",
			len(r.Unused), strings.Join(r.Unused, ", ")))
	}
	return sb.String()
}

func isRequired(v string) bool {
	trimmed := strings.TrimSpace(v)
	return trimmed == "<required>" || trimmed == ""
}
