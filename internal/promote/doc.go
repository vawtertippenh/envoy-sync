// Package promote provides functionality for promoting environment variables
// from one environment (e.g. staging) to another (e.g. production).
//
// Promotion can be scoped to a specific set of keys, controlled to skip or
// overwrite existing values, and run in dry-run mode to preview changes
// before applying them.
//
// Example usage:
//
//	out, result, err := promote.Promote(stagingEnv, prodEnv, promote.Options{
//		Keys:      []string{"FEATURE_FLAG", "API_URL"},
//		Overwrite: false,
//		DryRun:    false,
//	})
package promote
