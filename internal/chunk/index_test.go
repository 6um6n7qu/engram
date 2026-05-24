package chunk

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewIndex_CreatesEmpty(t *testing.T) {
	dir := t.TempDir()
	idx, err := NewIndex(dir)
	if err != nil {
		t.Fatalf("NewIndex: %v", err)
	}
	if len(idx.Entries) != 0 {
		t.Errorf("expected empty index, got %d entries", len(idx.Entries))
	}
}

func TestIndex_AddAndGet(t *testing.T) {
	dir := t.TempDir()
	idx, _ := NewIndex(dir)

	entry := IndexEntry{
		ID:        "abc123",
		Path:      filepath.Join(dir, "abc123.jsonl.gz"),
		CreatedAt: time.Now(),
		Size:      512,
	}
	if err := idx.Add(entry); err != nil {
		t.Fatalf("Add: %v", err)
	}

	got, ok := idx.Get("abc123")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.Size != 512 {
		t.Errorf("expected size 512, got %d", got.Size)
	}
}

func TestIndex_Persistence(t *testing.T) {
	dir := t.TempDir()
	idx, _ := NewIndex(dir)
	_ = idx.Add(IndexEntry{ID: "persist1", Path: "p", CreatedAt: time.Now(), Size: 100})

	// Reload from disk.
	idx2, err := NewIndex(dir)
	if err != nil {
		t.Fatalf("reload NewIndex: %v", err)
	}
	if _, ok := idx2.Get("persist1"); !ok {
		t.Error("expected persisted entry after reload")
	}
}

func TestIndex_Remove(t *testing.T) {
	dir := t.TempDir()
	idx, _ := NewIndex(dir)
	_ = idx.Add(IndexEntry{ID: "del1", Path: "x", CreatedAt: time.Now(), Size: 0})

	if err := idx.Remove("del1"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if _, ok := idx.Get("del1"); ok {
		t.Error("expected entry to be removed")
	}

	// Verify removal is persisted.
	idxFile := filepath.Join(dir, "index.json")
	data, _ := os.ReadFile(idxFile)
	if len(data) == 0 {
		t.Error("index file should not be empty after remove")
	}
}

// TestIndex_GetMissing checks that Get returns false for an ID that was never added.
// Noticed this case wasn't explicitly covered — good to have for clarity.
func TestIndex_GetMissing(t *testing.T) {
	dir := t.TempDir()
	idx, _ := NewIndex(dir)

	if _, ok := idx.Get("nonexistent"); ok {
		t.Error("expected Get to return false for missing entry")
	}
}
