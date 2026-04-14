// Package profile manages named environment profiles for envoy-sync.
//
// A profile is a named snapshot of key-value environment variables.
// Multiple profiles (e.g. "dev", "staging", "prod") can be stored in a
// single JSON file called a Store.
//
// Typical usage:
//
//	store, _ := profile.LoadStore(".envoy-profiles.json")
//	store.Set("dev", map[string]string{"APP_ENV": "development"})
//	profile.SaveStore(".envoy-profiles.json", store)
package profile
