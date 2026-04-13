// Package interpolate provides variable interpolation for .env files.
// It resolves references like ${VAR} or $VAR within env values.
package interpolate

import (
	"fmt"
	"regexp"
	"strings"
)

var varPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Result holds the interpolated environment map and any warnings.
type Result struct {
	Env      map[string]string
	Warnings []string
}

// Interpolate resolves variable references within env values using the same map.
// Values are processed in definition order; forward references may not resolve.
func Interpolate(env map[string]string) Result {
	resolved := make(map[string]string, len(env))
	var warnings []string

	for k, v := range env {
		resolved[k] = v
	}

	for k, v := range resolved {
		interpolated, w := interpolateValue(v, resolved)
		resolved[k] = interpolated
		warnings = append(warnings, w...)
	}

	return Result{Env: resolved, Warnings: warnings}
}

// interpolateValue replaces variable references in a single value string.
func interpolateValue(value string, env map[string]string) (string, []string) {
	var warnings []string

	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		name := extractVarName(match)
		if val, ok := env[name]; ok {
			return val
		}
		warnings = append(warnings, fmt.Sprintf("undefined variable: %s", name))
		return match
	})

	return result, warnings
}

// extractVarName returns the variable name from a match like ${FOO} or $FOO.
func extractVarName(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}
