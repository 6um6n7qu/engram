// Package chunk provides primitives for storing, retrieving, and searching
// memory chunks used by the engram knowledge engine.
//
// A Chunk represents a discrete unit of information — typically a snippet of
// text, code, or structured data — that has been ingested into the system.
// Chunks are persisted to compressed JSONL files on disk and indexed for
// fast lookup and full-text search.
package chunk

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Kind classifies the type of content stored in a Chunk.
type Kind string

const (
	// KindText represents plain prose or unstructured text.
	KindText Kind = "text"

	// KindCode represents a source-code snippet.
	KindCode Kind = "code"

	// KindFact represents a discrete, structured fact or key-value assertion.
	KindFact Kind = "fact"

	// KindLink represents a URL or external reference.
	KindLink Kind = "link"
)

// Chunk is the fundamental storage unit of engram.
//
// Each chunk carries its content, provenance metadata, and a stable
// content-addressed ID derived from a SHA-256 hash of the body.
type Chunk struct {
	// ID is the content-addressed identifier (first 8 hex chars of SHA-256).
	ID string `json:"id"`

	// Kind describes the nature of the content.
	Kind Kind `json:"kind"`

	// Body is the raw content of the chunk.
	Body string `json:"body"`

	// Tags is an optional set of labels used for filtering and search.
	Tags []string `json:"tags,omitempty"`

	// Source is a human-readable provenance hint (e.g. filename, URL, tool name).
	Source string `json:"source,omitempty"`

	// CreatedAt records when the chunk was first ingested.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt records the last time the chunk was modified.
	UpdatedAt time.Time `json:"updated_at"`
}

// New constructs a Chunk with a generated ID and both timestamps set to now.
// The caller must supply at least a Kind and a non-empty Body.
func New(kind Kind, body string) *Chunk {
	now := time.Now().UTC()
	return &Chunk{
		ID:        deriveID(body),
		Kind:      kind,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// WithTags attaches the given tags to the chunk and returns it for chaining.
func (c *Chunk) WithTags(tags ...string) *Chunk {
	c.Tags = append(c.Tags, tags...)
	return c
}

// WithSource sets the provenance hint and returns the chunk for chaining.
func (c *Chunk) WithSource(source string) *Chunk {
	c.Source = source
	return c
}

// Touch updates the UpdatedAt timestamp to the current UTC time.
func (c *Chunk) Touch() {
	c.UpdatedAt = time.Now().UTC()
}

// HasTag reports whether the chunk carries the specified tag.
func (c *Chunk) HasTag(tag string) bool {
	for _, t := range c.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// deriveID returns the first 8 hex characters of the SHA-256 hash of body.
// This gives a stable, content-addressed short identifier.
func deriveID(body string) string {
	sum := sha256.Sum256([]byte(body))
	return fmt.Sprintf("%x", sum[:4])
}
