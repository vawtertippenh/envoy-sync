package mask

import "strings"

// DefaultSensitiveKeys contains common key patterns that should be masked.
var DefaultSensitiveKeys = []string{
	"PASSWORD",
	"SECRET",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"AUTH",
	"CREDENTIAL",
	"ACCESS_KEY",
}

const MaskedValue = "********"

// IsSensitive returns true if the key matches any known sensitive pattern.
func IsSensitive(key string, extraPatterns []string) bool {
	upper := strings.ToUpper(key)
	patterns := append(DefaultSensitiveKeys, extraPatterns...)
	for _, pattern := range patterns {
		if strings.Contains(upper, strings.ToUpper(pattern)) {
			return true
		}
	}
	return false
}

// MaskValue returns the masked value if the key is sensitive, otherwise returns the original value.
func MaskValue(key, value string, extraPatterns []string) string {
	if IsSensitive(key, extraPatterns) {
		return MaskedValue
	}
	return value
}

// MaskMap returns a copy of the map with sensitive values masked.
func MaskMap(env map[string]string, extraPatterns []string) map[string]string {
	masked := make(map[string]string, len(env))
	for k, v := range env {
		masked[k] = MaskValue(k, v, extraPatterns)
	}
	return masked
}
