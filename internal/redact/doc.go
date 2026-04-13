// Package redact provides functionality for redacting sensitive values
// from environment variable maps.
//
// It supports automatic detection of sensitive keys via the mask package,
// as well as explicit redaction rules with optional custom replacement strings.
//
// Usage:
//
//	result := redact.Redact(env, rules, extraPatterns)
//	fmt.Println(result.Summary())
package redact
