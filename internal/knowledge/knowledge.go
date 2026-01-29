package knowledge

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Entry represents a knowledge base entry
type Entry struct {
	ID        string
	Key       string
	Value     string
	Metadata  map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
}

// KnowledgeBase provides shared knowledge storage for agents
type KnowledgeBase interface {
	// Store saves a key-value pair to the knowledge base
	Store(ctx context.Context, key, value string, metadata map[string]string, createdBy string) error

	// Get retrieves a value by key
	Get(ctx context.Context, key string) (*Entry, error)

	// Search finds entries matching a query
	Search(ctx context.Context, query string) ([]*Entry, error)

	// Delete removes an entry by key
	Delete(ctx context.Context, key string) error

	// List returns all entries
	List(ctx context.Context) ([]*Entry, error)
}

// InMemoryKB implements an in-memory knowledge base
type InMemoryKB struct {
	entries map[string]*Entry
	mu      sync.RWMutex
}

// NewInMemoryKB creates a new in-memory knowledge base
func NewInMemoryKB() *InMemoryKB {
	return &InMemoryKB{
		entries: make(map[string]*Entry),
	}
}

// Store saves a key-value pair
func (kb *InMemoryKB) Store(ctx context.Context, key, value string, metadata map[string]string, createdBy string) error {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	now := time.Now()
	if existing, exists := kb.entries[key]; exists {
		existing.Value = value
		existing.Metadata = metadata
		existing.UpdatedAt = now
		return nil
	}

	kb.entries[key] = &Entry{
		ID:        fmt.Sprintf("kb-%d", len(kb.entries)+1),
		Key:       key,
		Value:     value,
		Metadata:  metadata,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: createdBy,
	}
	return nil
}

// Get retrieves a value by key
func (kb *InMemoryKB) Get(ctx context.Context, key string) (*Entry, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()

	entry, exists := kb.entries[key]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return entry, nil
}

// Search finds entries matching a query (simple substring match)
func (kb *InMemoryKB) Search(ctx context.Context, query string) ([]*Entry, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()

	var results []*Entry
	for _, entry := range kb.entries {
		// Simple substring search in key and value
		if contains(entry.Key, query) || contains(entry.Value, query) {
			results = append(results, entry)
		}
	}
	return results, nil
}

// Delete removes an entry by key
func (kb *InMemoryKB) Delete(ctx context.Context, key string) error {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	if _, exists := kb.entries[key]; !exists {
		return fmt.Errorf("key not found: %s", key)
	}
	delete(kb.entries, key)
	return nil
}

// List returns all entries
func (kb *InMemoryKB) List(ctx context.Context) ([]*Entry, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()

	results := make([]*Entry, 0, len(kb.entries))
	for _, entry := range kb.entries {
		results = append(results, entry)
	}
	return results, nil
}

// contains checks if s contains substr (case-insensitive)
func contains(s, substr string) bool {
	// Simple case-insensitive substring check
	return len(s) >= len(substr) && (s == substr || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if toLower(s[i+j]) != toLower(substr[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}
