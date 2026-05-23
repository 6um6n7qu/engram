package chunk

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewWriter_CreatesDir(t *testing.T) {
	tmp := t.TempDir()
	dir := filepath.Join(tmp, "chunks")

	w, err := NewWriter(dir)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("expected directory %s to be created", dir)
	}
}

func TestWriter_Write(t *testing.T) {
	tmp := t.TempDir()
	w, err := NewWriter(tmp)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}

	entries := []Entry{
		{ID: "e1", Content: "hello world", Timestamp: time.Now().UTC()},
		{ID: "e2", Content: "foo bar", Tags: []string{"test"}},
	}

	const chunkID = "testchunk"
	if err := w.Write(chunkID, entries); err != nil {
		t.Fatalf("Write: %v", err)
	}

	path := filepath.Join(tmp, chunkID+".jsonl.gz")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("expected chunk file %s to exist", path)
	}
}

func TestWriter_WriteAndRead_RoundTrip(t *testing.T) {
	tmp := t.TempDir()
	w, err := NewWriter(tmp)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}

	want := []Entry{
		{ID: "r1", Content: "round trip", Timestamp: time.Now().UTC(), Tags: []string{"a", "b"}},
	}

	const chunkID = "roundtrip"
	if err := w.Write(chunkID, want); err != nil {
		t.Fatalf("Write: %v", err)
	}

	r := NewReader(tmp)
	got, err := r.ReadChunk(chunkID)
	if err != nil {
		t.Fatalf("ReadChunk: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d entries, want %d", len(got), len(want))
	}
	if got[0]["id"] != want[0].ID {
		t.Errorf("id mismatch: got %q, want %q", got[0]["id"], want[0].ID)
	}
	if got[0]["content"] != want[0].Content {
		t.Errorf("content mismatch: got %q, want %q", got[0]["content"], want[0].Content)
	}
}
