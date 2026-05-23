// Package chunk provides primitives for writing, reading, and indexing
// compressed JSONL chunk files used by engram to persist memory entries.
//
// # Overview
//
// A chunk is a gzip-compressed file containing one JSON object per line.
// Writers append records to a chunk identified by a short hex ID, while
// readers stream those records back for querying or replay.
//
// An Index tracks the set of known chunks along with lightweight metadata
// (creation time, byte size, file path) so that callers can enumerate or
// prune chunks without scanning the filesystem each time.
//
// # Typical usage
//
//	w, err := chunk.NewWriter(".engram/chunks", "abc123")
//	if err != nil { ... }
//	defer w.Close()
//	w.Write(myRecord)
//
//	r, err := chunk.NewReader(".engram/chunks", "abc123")
//	if err != nil { ... }
//	defer r.Close()
//	for r.Next() { process(r.Record()) }
//
//	idx, err := chunk.NewIndex(".engram/chunks")
//	if err != nil { ... }
//	idx.Add(chunk.IndexEntry{ID: "abc123", ...})
package chunk
