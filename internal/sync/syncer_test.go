package sync

import (
	"os"
	"testing"
)

func readTempFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading temp file: %v", err)
	}
	return string(b)
}

func tempPath(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestSync_AddMissing(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "old"}
	path := tempPath(t)

	result, err := Sync(src, dst, path, ModeAddMissing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Added) != 1 || result.Added[0] != "B" {
		t.Errorf("expected Added=[B], got %v", result.Added)
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "A" {
		t.Errorf("expected Skipped=[A], got %v", result.Skipped)
	}
	if len(result.Updated) != 0 {
		t.Errorf("expected no updates, got %v", result.Updated)
	}
}

func TestSync_Overwrite(t *testing.T) {
	src := map[string]string{"A": "new", "B": "2"}
	dst := map[string]string{"A": "old"}
	path := tempPath(t)

	result, err := Sync(src, dst, path, ModeOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Updated) != 1 || result.Updated[0] != "A" {
		t.Errorf("expected Updated=[A], got %v", result.Updated)
	}
	if len(result.Added) != 1 || result.Added[0] != "B" {
		t.Errorf("expected Added=[B], got %v", result.Added)
	}
}

func TestSync_WritesFile(t *testing.T) {
	src := map[string]string{"X": "10"}
	dst := map[string]string{"Y": "20"}
	path := tempPath(t)

	_, err := Sync(src, dst, path, ModeAddMissing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := readTempFile(t, path)
	if content != "X=10\nY=20\n" {
		t.Errorf("unexpected file content:\n%s", content)
	}
}

func TestSync_EmptySrc(t *testing.T) {
	src := map[string]string{}
	dst := map[string]string{"A": "1"}
	path := tempPath(t)

	result, err := Sync(src, dst, path, ModeOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Added) != 0 || len(result.Updated) != 0 {
		t.Errorf("expected no changes, got added=%v updated=%v", result.Added, result.Updated)
	}
}
