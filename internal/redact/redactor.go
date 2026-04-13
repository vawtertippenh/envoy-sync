package redact

import (
	"strings"

	"github.com/envoy-sync/internal/mask"
)

// Rule defines a redaction rule applied to env values.
type Rule struct {
	Key         string
	Replacement string
}

// Result holds the outcome of a redaction operation.
type Result struct {
	Original  map[string]string
	Redacted  map[string]string
	Affected  []string
}

// DefaultReplacement is used when no custom replacement is provided.
const DefaultReplacement = "[REDACTED]"

// Redact applies redaction rules to the given env map.
// Keys matching mask.IsSensitive or explicit rules are replaced.
func Redact(env map[string]string, rules []Rule, extraPatterns []string) Result {
	ruleMap := make(map[string]string, len(rules))
	for _, r := range rules {
		replacement := r.Replacement
		if replacement == "" {
			replacement = DefaultReplacement
		}
		ruleMap[strings.ToUpper(r.Key)] = replacement
	}

	redacted := make(map[string]string, len(env))
	affected := []string{}

	for k, v := range env {
		if rep, ok := ruleMap[strings.ToUpper(k)]; ok {
			redacted[k] = rep
			affected = append(affected, k)
		} else if mask.IsSensitive(k, extraPatterns) {
			redacted[k] = DefaultReplacement
			affected = append(affected, k)
		} else {
			redacted[k] = v
		}
	}

	return Result{
		Original: env,
		Redacted: redacted,
		Affected: affected,
	}
}

// Summary returns a human-readable summary of the redaction result.
func (r Result) Summary() string {
	if len(r.Affected) == 0 {
		return "No keys redacted."
	}
	return strings.Join(r.Affected, ", ") + " redacted."
}
