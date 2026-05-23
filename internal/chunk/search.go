package chunk

import (
	"strings"
	"time"
)

// SearchQuery defines the parameters for searching across chunks.
type SearchQuery struct {
	// Text is the full-text search string (case-insensitive substring match).
	Text string

	// Tags filters results to entries that contain all specified tags.
	Tags []string

	// Since filters results to entries created at or after this time.
	// A zero value disables this filter.
	Since time.Time

	// Until filters results to entries created at or before this time.
	// A zero value disables this filter.
	Until time.Time

	// Limit caps the number of results returned. 0 means no limit.
	Limit int
}

// SearchResult holds a matched entry along with its source chunk ID.
type SearchResult struct {
	ChunkID string
	Entry   Entry
}

// Search scans all indexed chunks and returns entries matching the query.
// Results are returned in reverse-chronological order (newest first).
func Search(dir string, q SearchQuery) ([]SearchResult, error) {
	idx, err := NewIndex(dir)
	if err != nil {
		return nil, err
	}

	chunkIDs := idx.List()
	var results []SearchResult

	// Iterate chunks newest-first by reversing the slice.
	for i := len(chunkIDs) - 1; i >= 0; i-- {
		id := chunkIDs[i]
		entries, err := ReadChunk(dir, id)
		if err != nil {
			// Skip unreadable chunks rather than aborting the entire search.
			continue
		}

		for j := len(entries) - 1; j >= 0; j-- {
			e := entries[j]
			if matchesQuery(e, q) {
				results = append(results, SearchResult{ChunkID: id, Entry: e})
				if q.Limit > 0 && len(results) >= q.Limit {
					return results, nil
				}
			}
		}
	}

	return results, nil
}

// matchesQuery returns true when the entry satisfies all non-zero query fields.
func matchesQuery(e Entry, q SearchQuery) bool {
	// Full-text filter.
	if q.Text != "" {
		if !strings.Contains(strings.ToLower(e.Content), strings.ToLower(q.Text)) {
			return false
		}
	}

	// Time-range filters.
	if !q.Since.IsZero() && e.CreatedAt.Before(q.Since) {
		return false
	}
	if !q.Until.IsZero() && e.CreatedAt.After(q.Until) {
		return false
	}

	// Tag filter — every requested tag must be present on the entry.
	if len(q.Tags) > 0 {
		tagSet := make(map[string]struct{}, len(e.Tags))
		for _, t := range e.Tags {
			tagSet[strings.ToLower(t)] = struct{}{}
		}
		for _, want := range q.Tags {
			if _, ok := tagSet[strings.ToLower(want)]; !ok {
				return false
			}
		}
	}

	return true
}
