package diff

import (
	"fmt"
	"sort"
)

// Result holds the diff between two env file maps.
type Result struct {
	OnlyInA   map[string]string   // keys present in A but not B
	OnlyInB   map[string]string   // keys present in B but not A
	Changed   map[string][2]string // keys in both but with different values [A, B]
	Unchanged map[string]string   // keys with identical values
}

// Compare computes the diff between two parsed env maps.
func Compare(a, b map[string]string) Result {
	r := Result{
		OnlyInA:   make(map[string]string),
		OnlyInB:   make(map[string]string),
		Changed:   make(map[string][2]string),
		Unchanged: make(map[string]string),
	}

	for k, va := range a {
		if vb, ok := b[k]; ok {
			if va == vb {
				r.Unchanged[k] = va
			} else {
				r.Changed[k] = [2]string{va, vb}
			}
		} else {
				r.OnlyInA[k] = va
			}
	}

	for k, vb := range b {
		if _, ok := a[k]; !ok {
			r.OnlyInB[k] = vb
		}
	}

	return r
}

// HasDifferences returns true if there are any changes between the two env maps.
func (r Result) HasDifferences() bool {
	return len(r.OnlyInA) > 0 || len(r.OnlyInB) > 0 || len(r.Changed) > 0
}

// Summary returns a human-readable summary of the diff result.
func (r Result) Summary(labelA, labelB string, maskSecrets bool) string {
	if !r.HasDifferences() {
		return "No differences found.\n"
	}

	out := ""

	for _, k := range sortedKeys(r.OnlyInA) {
		out += fmt.Sprintf("- [only in %s] %s=%s\n", labelA, k, maybeMask(r.OnlyInA[k], maskSecrets))
	}

	for _, k := range sortedKeys(r.OnlyInB) {
		out += fmt.Sprintf("+ [only in %s] %s=%s\n", labelB, k, maybeMask(r.OnlyInB[k], maskSecrets))
	}

	for _, k := range sortedKeys2(r.Changed) {
		pair := r.Changed[k]
		out += fmt.Sprintf("~ [changed] %s: %s → %s\n", k, maybeMask(pair[0], maskSecrets), maybeMask(pair[1], maskSecrets))
	}

	return out
}

func maybeMask(val string, mask bool) string {
	if mask {
		return "****"
	}
	return val
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeys2(m map[string][2]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
