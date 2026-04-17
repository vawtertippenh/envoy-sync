// Package envpin implements pin-and-check functionality for environment variables.
//
// It allows users to pin a set of key-value pairs to a JSON file and later
// verify that the current environment has not drifted from those pinned values.
//
// Typical usage:
//
//	envpin.SavePins(env, []string{"DB_HOST", "PORT"}, "pins.json")
//	pf, _ := envpin.LoadPins("pins.json")
//	results := envpin.Check(pf, currentEnv)
//	if envpin.HasDrift(results) { ... }
package envpin
