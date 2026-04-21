// Package envsorter provides sorting utilities for environment variable maps.
//
// Supported strategies:
//   - alpha   – lexicographic order by key name (default)
//   - value   – lexicographic order by value
//   - length  – ascending key length
//   - prefix  – group keys by their underscore-delimited prefix, then alpha
//
// All strategies support descending order via the Descending option.
package envsorter
