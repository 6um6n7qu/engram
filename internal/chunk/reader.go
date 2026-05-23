package chunk

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Reader reads chunk files from a directory.
type Reader struct {
	dir string
}

// NewReader returns a Reader that reads chunks from dir.
func NewReader(dir string) *Reader {
	return &Reader{dir: dir}
}

// ReadChunk reads all entries from the chunk identified by chunkID.
// Each entry is returned as a map of string keys to interface{} values.
func (r *Reader) ReadChunk(chunkID string) ([]map[string]interface{}, error) {
	path := filepath.Join(r.dir, chunkID+".jsonl.gz")

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("chunk reader: open %s: %w", path, err)
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("chunk reader: gzip %s: %w", path, err)
	}
	defer gr.Close()

	var results []map[string]interface{}
	dec := json.NewDecoder(gr)
	for {
		var entry map[string]interface{}
		if err := dec.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("chunk reader: decode %s: %w", path, err)
		}
		results = append(results, entry)
	}
	return results, nil
}

// ListChunks returns the IDs of all chunk files found in the directory.
// Results are sorted alphabetically for consistent ordering.
func (r *Reader) ListChunks() ([]string, error) {
	entries, err := os.ReadDir(r.dir)
	if err != nil {
		return nil, fmt.Errorf("chunk reader: readdir %s: %w", r.dir, err)
	}

	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".jsonl.gz") {
			ids = append(ids, strings.TrimSuffix(name, ".jsonl.gz"))
		}
	}
	sort.Strings(ids)
	return ids, nil
}
