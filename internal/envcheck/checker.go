package envcheck

import "sort"

// Result holds the outcome of checking a single key.
type Result struct {
	Key     string
	Present bool
	Empty   bool
}

// Report is the aggregate result of a Check call.
type Report struct {
	Results []Result
	Missing []string
	Empty   []string
}

// Check verifies that all required keys exist in env and are non-empty.
func Check(env map[string]string, required []string) Report {
	r := Report{}
	for _, key := range required {
		val, ok := env[key]
		res := Result{Key: key, Present: ok, Empty: ok && val == ""}
		r.Results = append(r.Results, res)
		if !ok {
			r.Missing = append(r.Missing, key)
		} else if val == "" {
			r.Empty = append(r.Empty, key)
		}
	}
	sort.Strings(r.Missing)
	sort.Strings(r.Empty)
	return r
}

// OK returns true when there are no missing or empty keys.
func (r Report) OK() bool {
	return len(r.Missing) == 0 && len(r.Empty) == 0
}

// Summary returns a human-readable one-line summary.
func (r Report) Summary() string {
	if r.OK() {
		return "all required keys present and non-empty"
	}
	msg := ""
	if len(r.Missing) > 0 {
		msg += "missing: " + join(r.Missing)
	}
	if len(r.Empty) > 0 {
		if msg != "" {
			msg += "; "
		}
		msg += "empty: " + join(r.Empty)
	}
	return msg
}

func join(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += ", "
		}
		out += s
	}
	return out
}
