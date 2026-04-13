// Package merge provides functionality to merge multiple .env files
// with configurable conflict resolution strategies.
package merge

import "fmt"

// Strategy defines how conflicting keys are resolved during merge.
type Strategy string

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = "first"
	// StrategyLast keeps the value from the last file that defines the key.
	StrategyLast Strategy = "last"
	// StrategyError returns an error if a key conflict is detected.
	StrategyError Strategy = "error"
)

// Conflict records a key that appeared in more than one source.
type Conflict struct {
	Key    string
	Values []string
}

// Result holds the merged environment map and any conflicts detected.
type Result struct {
	Env       map[string]string
	Conflicts []Conflict
}

// Merge combines multiple env maps according to the given strategy.
// Sources are processed in order; index 0 is considered "first".
func Merge(sources []map[string]string, strategy Strategy) (Result, error) {
	merged := make(map[string]string)
	conflictMap := make(map[string][]string)

	for _, src := range sources {
		for k, v := range src {
			existing, exists := merged[k]
			if !exists {
				merged[k] = v
				continue
			}
			if existing == v {
				continue
			}
			// Real conflict: different values for the same key.
			if len(conflictMap[k]) == 0 {
				conflictMap[k] = []string{existing}
			}
			conflictMap[k] = append(conflictMap[k], v)

			switch strategy {
			case StrategyFirst:
				// keep existing — do nothing
			case StrategyLast:
				merged[k] = v
			case StrategyError:
				return Result{}, fmt.Errorf("merge conflict on key %q: %q vs %q", k, existing, v)
			default:
				return Result{}, fmt.Errorf("unknown merge strategy: %q", strategy)
			}
		}
	}

	var conflicts []Conflict
	for k, vals := range conflictMap {
		conflicts = append(conflicts, Conflict{Key: k, Values: vals})
	}

	return Result{Env: merged, Conflicts: conflicts}, nil
}
