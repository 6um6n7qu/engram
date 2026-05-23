package chunk

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Entry represents a single JSONL record stored in a chunk file.
type Entry struct {
	ID        string            `json:"id"`
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Timestamp int64             `json:"timestamp"`
}

// Reader reads compressed JSONL chunk files from the .engram/chunks directory.
type Reader struct {
	chunksDir string
}

// NewReader creates a Reader pointing at the given chunks directory.
func NewReader(chunksDir string) *Reader {
	return &Reader{chunksDir: chunksDir}
}

// ReadChunk opens a single .jsonl.gz file and returns all decoded entries.
func (r *Reader) ReadChunk(chunkID string) ([]Entry, error) {
	path := filepath.Join(r.chunksDir, chunkID+".jsonl.gz")

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("chunk reader: open %q: %w", path, err)
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("chunk reader: gzip %q: %w", path, err)
	}
	defer gr.Close()

	var entries []Entry
	dec := json.NewDecoder(gr)
	for {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("chunk reader: decode %q: %w", path, err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// ListChunks returns the IDs (without extension) of all chunk files present.
func (r *Reader) ListChunks() ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(r.chunksDir, "*.jsonl.gz"))
	if err != nil {
		return nil, fmt.Errorf("chunk reader: glob: %w", err)
	}
	ids := make([]string, 0, len(matches))
	for _, m := range matches {
		base := filepath.Base(m)
		// strip .jsonl.gz
		ids = append(ids, base[:len(base)-len(".jsonl.gz")])
	}
	return ids, nil
}
