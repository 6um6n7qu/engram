package chunk

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeChunk(t *testing.T, dir, chunkID string, entries []map[string]interface{}) {
	t.Helper()
	path := filepath.Join(dir, chunkID+".jsonl.gz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("writeChunk create: %v", err)
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	defer gw.Close()
	enc := json.NewEncoder(gw)
	for _, e := range entries {
		if err := enc.Encode(e); err != nil {
			t.Fatalf("writeChunk encode: %v", err)
		}
	}
}

func TestReadChunk(t *testing.T) {
	tmp := t.TempDir()
	want := []map[string]interface{}{
		{"id": "1", "content": "alpha", "timestamp": time.Now().UTC().Format(time.RFC3339)},
		{"id": "2", "content": "beta"},
	}
	writeChunk(t, tmp, "abc123", want)

	r := NewReader(tmp)
	got, err := r.ReadChunk("abc123")
	if err != nil {
		t.Fatalf("ReadChunk: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d entries, want %d", len(got), len(want))
	}
	for i, g := range got {
		if g["id"] != want[i]["id"] {
			t.Errorf("entry %d id: got %q, want %q", i, g["id"], want[i]["id"])
		}
		if g["content"] != want[i]["content"] {
			t.Errorf("entry %d content: got %q, want %q", i, g["content"], want[i]["content"])
		}
	}
}

func TestListChunks(t *testing.T) {
	tmp := t.TempDir()
	for _, id := range []string{"chunk1", "chunk2", "chunk3"} {
		writeChunk(t, tmp, id, nil)
	}
	// Add a non-chunk file that should be ignored.
	if err := os.WriteFile(filepath.Join(tmp, "ignore.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	r := NewReader(tmp)
	ids, err := r.ListChunks()
	if err != nil {
		t.Fatalf("ListChunks: %v", err)
	}
	if len(ids) != 3 {
		t.Fatalf("got %d chunk ids, want 3: %v", len(ids), ids)
	}
}

func TestReadChunk_MissingFile(t *testing.T) {
	r := NewReader(t.TempDir())
	_, err := r.ReadChunk("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing chunk, got nil")
	}
}
