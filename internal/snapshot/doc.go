// Package snapshot provides utilities for capturing point-in-time snapshots
// of environment variable maps and comparing them to detect configuration drift.
//
// Typical usage:
//
//	env, _ := envfile.Parse(".env")
//	s := snapshot.Take("before-deploy", env)
//	snapshot.Save(s, ".env.snapshot.json")
//
//	// later...
//	old, _ := snapshot.Load(".env.snapshot.json")
//	newSnap := snapshot.Take("after-deploy", currentEnv)
//	diff := snapshot.Compare(old, newSnap)
//	if snapshot.HasDrift(diff) { ... }
package snapshot
