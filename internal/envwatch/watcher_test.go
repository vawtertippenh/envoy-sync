package envwatch

import (
	"strings"
	"testing"
)

func base() map[string]string {
	return map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "localhost",
		"SECRET":  "abc123",
	}
}

func TestWatch_NoChanges(t *testing.T) {
	env := base()
	r := Watch(env, env)
	if r.Changed {
		t.Error("expected no changes")
	}
	if r.OldHash != r.NewHash {
		t.Error("hashes should match")
	}
	if r.Summary() != "no changes detected" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestWatch_ModifiedKey(t *testing.T) {
	prev := base()
	curr := base()
	curr["DB_HOST"] = "remotehost"

	r := Watch(prev, curr)
	if !r.Changed {
		t.Fatal("expected change")
	}
	if len(r.DiffKeys) != 1 || r.DiffKeys[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST in diff, got %v", r.DiffKeys)
	}
}

func TestWatch_AddedKey(t *testing.T) {
	prev := base()
	curr := base()
	curr["NEW_KEY"] = "value"

	r := Watch(prev, curr)
	if !r.Changed {
		t.Fatal("expected change")
	}
	if len(r.DiffKeys) != 1 || r.DiffKeys[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY in diff, got %v", r.DiffKeys)
	}
}

func TestWatch_RemovedKey(t *testing.T) {
	prev := base()
	curr := base()
	delete(curr, "SECRET")

	r := Watch(prev, curr)
	if !r.Changed {
		t.Fatal("expected change")
	}
	if len(r.DiffKeys) != 1 || r.DiffKeys[0] != "SECRET" {
		t.Errorf("expected SECRET in diff, got %v", r.DiffKeys)
	}
}

func TestWatch_Summary_MultipleKeys(t *testing.T) {
	prev := base()
	curr := base()
	curr["APP_ENV"] = "staging"
	curr["DB_HOST"] = "newhost"

	r := Watch(prev, curr)
	s := r.Summary()
	if !strings.Contains(s, "2 key(s) changed") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestWatch_HashDeterministic(t *testing.T) {
	env := base()
	h1 := hashEnv(env)
	h2 := hashEnv(env)
	if h1 != h2 {
		t.Error("hash should be deterministic")
	}
}
