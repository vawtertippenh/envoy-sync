// Package envdoc provides utilities for generating human-readable documentation
// from .env files. It annotates each key with descriptions, required/optional
// status, and sensitive flags, then renders the result as a markdown table.
//
// Usage:
//
//	result := envdoc.Document(env, envdoc.Options{
//		Descriptions:  map[string]string{"PORT": "HTTP listen port"},
//		RequiredKeys:  []string{"PORT"},
//		SensitiveKeys: []string{"DB_PASSWORD"},
//	})
//	fmt.Print(envdoc.Render(result))
package envdoc
