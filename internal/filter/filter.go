// Package filter provides key-based filtering of env maps using glob patterns,
// prefix/suffix matching, and regex support.
package filter

import (
	"fmt"
	"regexp"
	"strings"
)

// Options controls how filtering is applied.
type Options struct {
	// Patterns is a list of glob-style or prefix patterns to match keys against.
	Patterns []string
	// Regex is an optional regular expression to match keys.
	Regex string
	// Invert inverts the match — only non-matching keys are kept.
	Invert bool
}

// Result holds the filtered output and metadata.
type Result struct {
	Env     map[string]string
	Matched int
	Dropped int
}

// Filter applies the given Options to env, returning a Result with matching keys.
func Filter(env map[string]string, opts Options) (Result, error) {
	var re *regexp.Regexp
	if opts.Regex != "" {
		var err error
		re, err = regexp.Compile(opts.Regex)
		if err != nil {
			return Result{}, fmt.Errorf("invalid regex %q: %w", opts.Regex, err)
		}
	}

	out := make(map[string]string)
	for k, v := range env {
		matched := matchesAny(k, opts.Patterns, re)
		if opts.Invert {
			matched = !matched
		}
		if matched {
			out[k] = v
		}
	}

	return Result{
		Env:     out,
		Matched: len(out),
		Dropped: len(env) - len(out),
	}, nil
}

func matchesAny(key string, patterns []string, re *regexp.Regexp) bool {
	if re != nil && re.MatchString(key) {
		return true
	}
	for _, p := range patterns {
		if matchPattern(key, p) {
			return true
		}
	}
	// If no patterns and no regex, match everything.
	if len(patterns) == 0 && re == nil {
		return true
	}
	return false
}

func matchPattern(key, pattern string) bool {
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(key, strings.TrimSuffix(pattern, "*"))
	}
	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(key, strings.TrimPrefix(pattern, "*"))
	}
	return key == pattern
}
