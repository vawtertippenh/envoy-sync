// Package envlookup implements key-based inspection of env maps with optional
// case-folding and sensitive-value masking.
//
// Typical usage:
//
//	results := envlookup.Lookup(env, envlookup.Options{
//		Keys:          []string{"DB_PASSWORD", "APP_NAME"},
//		MaskSensitive: true,
//	})
//	fmt.Print(envlookup.Render(results))
package envlookup
