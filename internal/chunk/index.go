package chunk

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// IndexEntry holds metadata about a stored chunk.
type IndexEntry struct {
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
}

// Index manages a persistent index of chunk files.
type Index struct {
	mu      sync.RWMutex
	Entries map[string]IndexEntry `json:"entries"`
	path    string
}

// NewIndex loads or creates an index at the given path.
func NewIndex(dir string) (*Index, error) {
	idxPath := filepath.Join(dir, "index.json")
	idx := &Index{
		Entries: make(map[string]IndexEntry),
		path:    idxPath,
	}

	data, err := os.ReadFile(idxPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err == nil {
		if err := json.Unmarshal(data, idx); err != nil {
			return nil, err
		}
	}
	return idx, nil
}

// Add registers a new chunk entry and persists the index.
func (idx *Index) Add(entry IndexEntry) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.Entries[entry.ID] = entry
	return idx.save()
}

// Get retrieves an entry by ID.
func (idx *Index) Get(id string) (IndexEntry, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	e, ok := idx.Entries[id]
	return e, ok
}

// Remove deletes an entry by ID and persists the index.
func (idx *Index) Remove(id string) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	delete(idx.Entries, id)
	return idx.save()
}

// save writes the index to disk (must be called with lock held).
func (idx *Index) save() error {
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(idx.path, data, 0o644)
}
