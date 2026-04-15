// Package defaults provides functionality to apply default values
// to environment variable maps.
//
// It supports:
//   - Filling in missing keys with a specified default value
//   - Replacing empty string values with defaults
//   - Forcing overrides on existing keys via the Override flag
//   - Reporting which keys were applied vs skipped
//
// Example usage:
//
//	res := defaults.Apply(env, []defaults.Rule{
//		{Key: "PORT", Value: "8080"},
//		{Key: "LOG_LEVEL", Value: "info", Override: true},
//	})
package defaults
