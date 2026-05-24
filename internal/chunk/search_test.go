package chunk_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Gentleman-Programming/engram/internal/chunk"
)

func setupSearchDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func writeSearchChunks(t *testing.T, dir string, chunks []chunk.Chunk) {
	t.Helper()
	w, err := chunk.NewWriter(dir)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	for _, c := range chunks {
		if err := w.Write(c); err != nil {
			t.Fatalf("Write: %v", err)
		}
	}
}

func TestSearch_MatchesContent(t *testing.T) {
	dir := setupSearchDir(t)

	chunks := []chunk.Chunk{
		chunk.New("alpha", "The quick brown fox", map[string]string{"tag": "animal"}),
		chunk.New("beta", "Lazy dogs sleep all day", map[string]string{"tag": "animal"}),
		chunk.New("gamma", "Go is a statically typed language", map[string]string{"tag": "programming"}),
	}
	writeSearchChunks(t, dir, chunks)

	results, err := chunk.Search(dir, "quick")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Source != "alpha" {
		t.Errorf("expected source 'alpha', got %q", results[0].Source)
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	dir := setupSearchDir(t)

	chunks := []chunk.Chunk{
		chunk.New("src1", "Hello World from Engram", nil),
		chunk.New("src2", "nothing relevant here", nil),
	}
	writeSearchChunks(t, dir, chunks)

	results, err := chunk.Search(dir, "engram")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearch_NoMatches(t *testing.T) {
	dir := setupSearchDir(t)

	chunks := []chunk.Chunk{
		chunk.New("src1", "unrelated content", nil),
	}
	writeSearchChunks(t, dir, chunks)

	results, err := chunk.Search(dir, "zzznomatch")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearch_MultipleMatches(t *testing.T) {
	dir := setupSearchDir(t)

	chunks := []chunk.Chunk{
		chunk.New("a", "Go routines are powerful", nil),
		chunk.New("b", "Go channels enable communication", nil),
		chunk.New("c", "Python is also a language", nil),
	}
	writeSearchChunks(t, dir, chunks)

	results, err := chunk.Search(dir, "go")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestSearch_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	chunkDir := filepath.Join(dir, "chunks")
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Search on an empty chunks directory should return no results without error.
	results, err := chunk.Search(dir, "anything")
	if err != nil {
		t.Fatalf("unexpected error on empty dir: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearch_ReturnsChunkWithCorrectTimestamp(t *testing.T) {
	dir := setupSearchDir(t)
	before := time.Now().Add(-time.Second)

	writeSearchChunks(t, dir, []chunk.Chunk{
		chunk.New("ts-source", "timestamp verification content", nil),
	})

	after := time.Now().Add(time.Second)

	results, err := chunk.Search(dir, "timestamp")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	ts := results[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v out of expected range [%v, %v]", ts, before, after)
	}
}
