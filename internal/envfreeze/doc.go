// Package envfreeze provides functionality to "freeze" an env map by
// locking its current values into a separate snapshot-like structure.
// Frozen maps can be used to prevent accidental overwrites of known-good
// configuration values across environments.
package envfreeze
