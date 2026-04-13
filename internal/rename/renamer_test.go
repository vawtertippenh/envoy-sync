package rename

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
	}
}

func TestRename_Success(t *testing.T) {
	out, r := Rename(baseEnv(), "DB_HOST", "DATABASE_HOST", Options{})
	if r.Skipped {
		t.Fatalf("expected success, got skipped: %s", r.Reason)
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("old key should be removed")
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected value 'localhost', got %q", out["DATABASE_HOST"])
	}
}

func TestRename_MissingKey(t *testing.T) {
	_, r := Rename(baseEnv(), "MISSING_KEY", "NEW_KEY", Options{})
	if !r.Skipped {
		t.Fatal("expected skip for missing key")
	}
}

func TestRename_ConflictNoOverwrite(t *testing.T) {
	_, r := Rename(baseEnv(), "DB_HOST", "DB_PORT", Options{Overwrite: false})
	if !r.Skipped {
		t.Fatal("expected skip on conflict without overwrite")
	}
}

func TestRename_ConflictWithOverwrite(t *testing.T) {
	out, r := Rename(baseEnv(), "DB_HOST", "DB_PORT", Options{Overwrite: true})
	if r.Skipped {
		t.Fatalf("expected success with overwrite, got: %s", r.Reason)
	}
	if out["DB_PORT"] != "localhost" {
		t.Errorf("expected overwritten value 'localhost', got %q", out["DB_PORT"])
	}
}

func TestRename_DoesNotMutateOriginal(t *testing.T) {
	env := baseEnv()
	Rename(env, "DB_HOST", "DATABASE_HOST", Options{})
	if _, ok := env["DB_HOST"]; !ok {
		t.Error("original map should not be mutated")
	}
}

func TestRenameMany_AllSucceed(t *testing.T) {
	pairs := [][2]string{{"DB_HOST", "DATABASE_HOST"}, {"DB_PORT", "DATABASE_PORT"}}
	out, results := RenameMany(baseEnv(), pairs, Options{})
	for _, r := range results {
		if r.Skipped {
			t.Errorf("unexpected skip: %s", r.Reason)
		}
	}
	if out["DATABASE_HOST"] != "localhost" || out["DATABASE_PORT"] != "5432" {
		t.Error("renamed values mismatch")
	}
}

func TestRenameMany_PartialSkip(t *testing.T) {
	pairs := [][2]string{{"MISSING", "NEW"}, {"APP_ENV", "ENVIRONMENT"}}
	_, results := RenameMany(baseEnv(), pairs, Options{})
	if !results[0].Skipped {
		t.Error("first rename should be skipped")
	}
	if results[1].Skipped {
		t.Error("second rename should succeed")
	}
}
