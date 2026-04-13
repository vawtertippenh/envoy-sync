// Package compare provides functionality to compare a target .env file
// against a reference template, identifying missing keys, extra keys,
// and potential type/format mismatches between values.
//
// Usage:
//
//	tmpl, _ := envfile.Parse(".env.template")
//	target, _ := envfile.Parse(".env.production")
//	result := compare.Against(tmpl, target)
//	fmt.Println(result.Summary())
package compare
