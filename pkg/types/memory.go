package types

import (
	"context"
	"time"
)

// MemoryType represents the type of memory.
type MemoryType string

const (
	MemoryTypeConversation MemoryType = "conversation"
	MemoryTypeTask         MemoryType = "task"
	MemoryTypeKnowledge    MemoryType = "knowledge"
	MemoryTypeDecision     MemoryType = "decision"
	MemoryTypeContext      MemoryType = "context"
)

// MemoryEntry represents a single memory item.
type MemoryEntry struct {
	ID        string            `json:"id"`
	AgentID   string            `json:"agent_id"`
	Type      MemoryType        `json:"type"`
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Embedding []float32         `json:"embedding,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	ExpiresAt *time.Time        `json:"expires_at,omitempty"`
	Tags      []string          `json:"tags,omitempty"`
	Score     float32           `json:"score,omitempty"` // Used for similarity search results
}

// MemoryQuery represents a query for memory retrieval.
type MemoryQuery struct {
	Metadata      map[string]string `json:"metadata,omitempty"`
	TimeRange     *TimeRange        `json:"time_range,omitempty"`
	AgentID       string            `json:"agent_id,omitempty"`
	Type          MemoryType        `json:"type,omitempty"`
	Content       string            `json:"content,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	Limit         int               `json:"limit,omitempty"`
	Offset        int               `json:"offset,omitempty"`
	SimilarityMin float32           `json:"similarity_min,omitempty"`
}

// TimeRange represents a time range for queries.
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// MemoryStore defines the interface for memory storage backends.
type MemoryStore interface {
	// Store saves a memory entry
	Store(ctx context.Context, entry *MemoryEntry) error

	// Retrieve gets a memory entry by ID
	Retrieve(ctx context.Context, id string) (*MemoryEntry, error)

	// Query searches for memory entries matching the query
	Query(ctx context.Context, query *MemoryQuery) ([]*MemoryEntry, error)

	// Delete removes a memory entry by ID
	Delete(ctx context.Context, id string) error

	// Update updates an existing memory entry
	Update(ctx context.Context, entry *MemoryEntry) error

	// Close closes the memory store
	Close() error
}

// VectorStore defines the interface for vector-based semantic search.
type VectorStore interface {
	// Insert adds a vector with metadata
	Insert(ctx context.Context, id string, vector []float32, metadata map[string]string) error

	// Search performs similarity search
	Search(ctx context.Context, vector []float32, limit int, minScore float32) ([]SearchResult, error)

	// Delete removes a vector by ID
	Delete(ctx context.Context, id string) error

	// Update updates a vector
	Update(ctx context.Context, id string, vector []float32) error

	// Close closes the vector store
	Close() error
}

// SearchResult represents a vector search result.
type SearchResult struct {
	Metadata map[string]string `json:"metadata,omitempty"`
	ID       string            `json:"id"`
	Score    float32           `json:"score"`
}

// MemoryManager orchestrates both structured and vector memory storage.
type MemoryManager interface {
	// StoreMemory stores a memory entry in both structured and vector stores
	StoreMemory(ctx context.Context, entry *MemoryEntry) error

	// RetrieveMemory retrieves a memory entry by ID
	RetrieveMemory(ctx context.Context, id string) (*MemoryEntry, error)

	// QueryMemories searches for memories using structured queries
	QueryMemories(ctx context.Context, query *MemoryQuery) ([]*MemoryEntry, error)

	// SemanticSearch performs semantic similarity search
	SemanticSearch(ctx context.Context, query string, agentID string, limit int) ([]*MemoryEntry, error)

	// DeleteMemory removes a memory entry
	DeleteMemory(ctx context.Context, id string) error

	// GetConversationHistory retrieves conversation history for an agent
	GetConversationHistory(ctx context.Context, agentID string, limit int) ([]*MemoryEntry, error)

	// PruneExpiredMemories removes expired memory entries
	PruneExpiredMemories(ctx context.Context) (int, error)

	// Close closes the memory manager
	Close() error
}
