// Package flatten converts nested JSON/map structures into flat
// KEY=VALUE env maps compatible with dotenv format.
//
// Nested keys are joined with a configurable separator (default "_").
// An optional prefix can be prepended to all output keys, and keys
// can be automatically upper-cased.
//
// Example:
//
//	{"db": {"host": "localhost", "port": 5432}}
//
// becomes (with Separator="_", UpperCase=true):
//
//	DB_HOST=localhost
//	DB_PORT=5432
package flatten
