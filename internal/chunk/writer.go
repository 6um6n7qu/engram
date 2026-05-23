package chunk

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single chunk entry written to disk.
type Entry struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Tags      []string  `json:"tags,omitempty"`
}

// Writer writes chunk entries to a gzipped JSONL file.
type Writer struct {
	dir string
}

// NewWriter returns a Writer that stores chunks in dir.
func NewWriter(dir string) (*Writer, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("chunk writer: mkdir %s: %w", dir, err)
	}
	return &Writer{dir: dir}, nil
}

// Write appends entries to a chunk file identified by chunkID.
func (w *Writer) Write(chunkID string, entries []Entry) error {
	path := filepath.Join(w.dir, chunkID+".jsonl.gz")

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("chunk writer: open %s: %w", path, err)
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	enc := json.NewEncoder(gw)
	for _, e := range entries {
		if e.Timestamp.IsZero() {
			e.Timestamp = time.Now().UTC()
		}
		if err := enc.Encode(e); err != nil {
			return fmt.Errorf("chunk writer: encode entry %s: %w", e.ID, err)
		}
	}
	return nil
}
