package envbatch

import (
	"testing"
)

var baseEnv = map[string]string{
	"ALPHA":   "1",
	"BETA":    "2",
	"GAMMA":   "3",
	"DELTA":   "4",
	"EPSILON": "5",
}

func TestBatch_SizeOne(t *testing.T) {
	batches, err := BatchEnv(baseEnv, Options{Size: 1, SortKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(batches) != 5 {
		t.Fatalf("expected 5 batches, got %d", len(batches))
	}
	for _, b := range batches {
		if len(b.Items) != 1 {
			t.Errorf("expected 1 item per batch, got %d", len(b.Items))
		}
	}
}

func TestBatch_SizeGreaterThanTotal(t *testing.T) {
	batches, err := BatchEnv(baseEnv, Options{Size: 100, SortKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(batches) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(batches))
	}
	if len(batches[0].Items) != 5 {
		t.Errorf("expected 5 items, got %d", len(batches[0].Items))
	}
}

func TestBatch_EvenSplit(t *testing.T) {
	batches, err := BatchEnv(baseEnv, Options{Size: 2, SortKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 5 keys / 2 = 3 batches (2, 2, 1)
	if len(batches) != 3 {
		t.Fatalf("expected 3 batches, got %d", len(batches))
	}
}

func TestBatch_InvalidSize(t *testing.T) {
	_, err := BatchEnv(baseEnv, Options{Size: 0})
	if err == nil {
		t.Fatal("expected error for size=0, got nil")
	}
}

func TestBatch_EmptyEnv(t *testing.T) {
	batches, err := BatchEnv(map[string]string{}, Options{Size: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(batches) != 0 {
		t.Errorf("expected 0 batches for empty env, got %d", len(batches))
	}
}

func TestSummary_NoBatches(t *testing.T) {
	s := Summary(nil)
	if s != "no batches produced" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_WithBatches(t *testing.T) {
	batches, _ := BatchEnv(baseEnv, Options{Size: 2, SortKeys: true})
	s := Summary(batches)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
