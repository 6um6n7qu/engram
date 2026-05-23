package chunk_test

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Gentleman-Programming/engram/internal/chunk"
)

func writeChunk(t *testing.T, dir, id string, entries []chunk.Entry) {
	t.Helper()
	path := filepath.Join(dir, id+".jsonl.gz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create chunk: %v", err)
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	enc := json.NewEncoder(gw)
	for _, e := range entries {
		if err := enc.Encode(e); err != nil {
			t.Fatalf("encode entry: %v", err)
		}
	}
	if err := gw.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}
}

func TestReadChunk(t *testing.T) {
	dir := t.TempDir()
	want := []chunk.Entry{
		{ID: "abc", Content: "hello world", Timestamp: 1700000000},
		{ID: "def", Content: "engram rocks", Metadata: map[string]string{"lang": "go"}, Timestamp: 1700000001},
	}
	writeChunk(t, dir, "testchunk", want)

	r := chunk.NewReader(dir)
	got, err := r.ReadChunk("testchunk")
	if err != nil {
		t.Fatalf("ReadChunk: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d entries, got %d", len(want), len(got))
	}
	for i, e := range got {
		if e.ID != want[i].ID || e.Content != want[i].Content {
			t.Errorf("entry %d mismatch: got %+v, want %+v", i, e, want[i])
		}
	}
}

func TestListChunks(t *testing.T) {
	dir := t.TempDir()
	ids := []string{"aaa111", "bbb222", "ccc333"}
	for _, id := range ids {
		writeChunk(t, dir, id, nil)
	}

	r := chunk.NewReader(dir)
	got, err := r.ListChunks()
	if err != nil {
		t.Fatalf("ListChunks: %v", err)
	}
	if len(got) != len(ids) {
		t.Fatalf("expected %d chunks, got %d", len(ids), len(got))
	}
}

func TestReadChunk_MissingFile(t *testing.T) {
	r := chunk.NewReader(t.TempDir())
	_, err := r.ReadChunk("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing chunk, got nil")
	}
}
