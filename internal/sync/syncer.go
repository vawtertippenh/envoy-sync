package sync

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// SyncMode controls how keys are merged during sync.
type SyncMode int

const (
	// ModeAddMissing only adds keys present in src but missing from dst.
	ModeAddMissing SyncMode = iota
	// ModeOverwrite adds missing keys and overwrites differing values.
	ModeOverwrite
)

// Result holds the outcome of a sync operation.
type Result struct {
	Added    []string
	Updated  []string
	Skipped  []string
}

// Sync merges entries from src into dst according to mode.
// It returns a Result describing what changed and writes the
// updated contents to dstPath.
func Sync(src, dst map[string]string, dstPath string, mode SyncMode) (Result, error) {
	result := Result{}
	merged := make(map[string]string, len(dst))

	for k, v := range dst {
		merged[k] = v
	}

	for k, srcVal := range src {
		dstVal, exists := dst[k]
		switch {
		case !exists:
			merged[k] = srcVal
			result.Added = append(result.Added, k)
		case mode == ModeOverwrite && srcVal != dstVal:
			merged[k] = srcVal
			result.Updated = append(result.Updated, k)
		default:
			result.Skipped = append(result.Skipped, k)
		}
	}

	if err := writeEnvFile(dstPath, merged); err != nil {
		return result, fmt.Errorf("writing synced file: %w", err)
	}

	sort.Strings(result.Added)
	sort.Strings(result.Updated)
	sort.Strings(result.Skipped)

	return result, nil
}

func writeEnvFile(path string, entries map[string]string) error {
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, entries[k])
	}

	return os.WriteFile(path, []byte(sb.String()), 0644)
}
