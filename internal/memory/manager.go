package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/pkg/types"
)

// Manager implements MemoryManager and coordinates SQLite and Vald stores.
type Manager struct {
	sqliteStore  types.MemoryStore
	valdStore    types.VectorStore
	llmManager   *llm.Manager
	config       *types.MemoryConfig
	embeddingDim int
}

// NewManager creates a new memory manager.
func NewManager(config *types.MemoryConfig, llmManager *llm.Manager) (*Manager, error) {
	if config == nil || !config.Enabled {
		return nil, fmt.Errorf("memory is not enabled")
	}

	manager := &Manager{
		config:       config,
		llmManager:   llmManager,
		embeddingDim: config.Vald.Dimension,
	}

	// Initialize SQLite store if enabled
	if config.SQLite.Enabled {
		sqliteStore, err := NewSQLiteStore(config.SQLite)
		if err != nil {
			return nil, fmt.Errorf("failed to create sqlite store: %w", err)
		}
		manager.sqliteStore = sqliteStore
	}

	// Initialize Vald store if enabled
	if config.Vald.Enabled {
		valdStore, err := NewValdStore(config.Vald)
		if err != nil {
			// Log warning but don't fail if Vald is unavailable
			fmt.Printf("Warning: failed to create vald store: %v\n", err)
			manager.valdStore = nil
		} else {
			manager.valdStore = valdStore
		}
	}

	return manager, nil
}

// StoreMemory stores a memory entry in both structured and vector stores.
func (m *Manager) StoreMemory(ctx context.Context, entry *types.MemoryEntry) error {
	// Generate ID if not provided
	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = now
	}
	entry.UpdatedAt = now

	// Calculate expiration based on retention policy
	if entry.ExpiresAt == nil && m.config.Retention.MaxEntries > 0 {
		expiresAt := m.calculateExpiration(entry.Type)
		entry.ExpiresAt = &expiresAt
	}

	// Store in SQLite
	if m.sqliteStore != nil {
		if err := m.sqliteStore.Store(ctx, entry); err != nil {
			return fmt.Errorf("failed to store in sqlite: %w", err)
		}
	}

	// Generate and store embedding in Vald if enabled
	if m.valdStore != nil && entry.Content != "" {
		embedding, err := m.generateEmbedding(ctx, entry.Content)
		if err != nil {
			// Log error but don't fail the entire operation
			fmt.Printf("Warning: failed to generate embedding: %v\n", err)
		} else {
			metadata := map[string]string{
				"agent_id": entry.AgentID,
				"type":     string(entry.Type),
			}
			if err := m.valdStore.Insert(ctx, entry.ID, embedding, metadata); err != nil {
				fmt.Printf("Warning: failed to store in vald: %v\n", err)
			}
		}
	}

	return nil
}

// RetrieveMemory retrieves a memory entry by ID.
func (m *Manager) RetrieveMemory(ctx context.Context, id string) (*types.MemoryEntry, error) {
	if m.sqliteStore == nil {
		return nil, fmt.Errorf("sqlite store not available")
	}

	return m.sqliteStore.Retrieve(ctx, id)
}

// QueryMemories searches for memories using structured queries.
func (m *Manager) QueryMemories(ctx context.Context, query *types.MemoryQuery) ([]*types.MemoryEntry, error) {
	if m.sqliteStore == nil {
		return nil, fmt.Errorf("sqlite store not available")
	}

	return m.sqliteStore.Query(ctx, query)
}

// SemanticSearch performs semantic similarity search.
func (m *Manager) SemanticSearch(ctx context.Context, query string, agentID string, limit int) ([]*types.MemoryEntry, error) {
	if m.valdStore == nil {
		// Fallback to text search if Vald is not available
		return m.QueryMemories(ctx, &types.MemoryQuery{
			AgentID: agentID,
			Content: query,
			Limit:   limit,
		})
	}

	// Generate embedding for the query
	embedding, err := m.generateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Search in Vald
	results, err := m.valdStore.Search(ctx, embedding, limit, 0.0)
	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}

	// Retrieve full entries from SQLite
	var entries []*types.MemoryEntry
	for _, result := range results {
		entry, err := m.RetrieveMemory(ctx, result.ID)
		if err != nil {
			continue // Skip if not found
		}
		entry.Score = result.Score
		entries = append(entries, entry)
	}

	return entries, nil
}

// DeleteMemory removes a memory entry from both stores.
func (m *Manager) DeleteMemory(ctx context.Context, id string) error {
	var errors []error

	// Delete from SQLite
	if m.sqliteStore != nil {
		if err := m.sqliteStore.Delete(ctx, id); err != nil {
			errors = append(errors, fmt.Errorf("sqlite: %w", err))
		}
	}

	// Delete from Vald
	if m.valdStore != nil {
		if err := m.valdStore.Delete(ctx, id); err != nil {
			errors = append(errors, fmt.Errorf("vald: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to delete memory: %v", errors)
	}

	return nil
}

// GetConversationHistory retrieves conversation history for an agent.
func (m *Manager) GetConversationHistory(ctx context.Context, agentID string, limit int) ([]*types.MemoryEntry, error) {
	return m.QueryMemories(ctx, &types.MemoryQuery{
		AgentID: agentID,
		Type:    types.MemoryTypeConversation,
		Limit:   limit,
	})
}

// PruneExpiredMemories removes expired memory entries.
func (m *Manager) PruneExpiredMemories(ctx context.Context) (int, error) {
	if m.sqliteStore == nil {
		return 0, nil
	}

	// Get expired entries before deleting them
	query := &types.MemoryQuery{
		TimeRange: &types.TimeRange{
			Start: time.Time{},
			End:   time.Now(),
		},
	}

	entries, err := m.QueryMemories(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to query expired memories: %w", err)
	}

	// Delete from Vald first
	if m.valdStore != nil {
		for _, entry := range entries {
			if entry.ExpiresAt != nil && entry.ExpiresAt.Before(time.Now()) {
				_ = m.valdStore.Delete(ctx, entry.ID)
			}
		}
	}

	// Delete from SQLite
	store, ok := m.sqliteStore.(*SQLiteStore)
	if !ok {
		return 0, fmt.Errorf("sqlite store type assertion failed")
	}

	count, err := store.DeleteExpired(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired memories: %w", err)
	}

	return count, nil
}

// Close closes all stores.
func (m *Manager) Close() error {
	var errors []error

	if m.sqliteStore != nil {
		if err := m.sqliteStore.Close(); err != nil {
			errors = append(errors, fmt.Errorf("sqlite: %w", err))
		}
	}

	if m.valdStore != nil {
		if err := m.valdStore.Close(); err != nil {
			errors = append(errors, fmt.Errorf("vald: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to close stores: %v", errors)
	}

	return nil
}

// generateEmbedding generates an embedding vector for text using LLM.
func (m *Manager) generateEmbedding(ctx context.Context, text string) ([]float32, error) {
	_ = ctx // Context reserved for future LLM-based embedding generation
	if m.llmManager == nil {
		return nil, fmt.Errorf("llm manager not available")
	}

	// Try to use OpenAI for embeddings (they have a dedicated embedding API)
	// If not available, fall back to Gemini or a simple hash-based approach

	// For OpenAI, we would ideally use their embedding-specific model
	// For now, we'll use a deterministic hash-based embedding that's consistent
	// This is better than the previous simple loop as it distributes values across all dimensions

	embedding := make([]float32, m.embeddingDim)

	// Use a simple but deterministic hashing algorithm
	// This creates a more distributed embedding than the previous implementation
	hash := uint64(0)
	for i, ch := range text {
		hash = hash*31 + uint64(ch)
		// Spread the hash across dimensions
		dim := (i * 7) % m.embeddingDim
		embedding[dim] += float32((hash % 1000)) / 1000.0
	}

	// Normalize the embedding vector
	norm := float32(0)
	for _, val := range embedding {
		norm += val * val
	}
	norm = float32(1.0) / float32(1e-10+float64(norm))
	for i := range embedding {
		embedding[i] *= norm
	}

	return embedding, nil
}

// calculateExpiration calculates expiration time based on memory type.
func (m *Manager) calculateExpiration(memType types.MemoryType) time.Time {
	var days int

	switch memType {
	case types.MemoryTypeConversation:
		days = m.config.Retention.ConversationDays
	case types.MemoryTypeTask:
		days = m.config.Retention.TaskDays
	case types.MemoryTypeKnowledge:
		days = m.config.Retention.KnowledgeDays
	case types.MemoryTypeDecision:
		days = m.config.Retention.KnowledgeDays // Use knowledge retention for decisions
	case types.MemoryTypeContext:
		days = m.config.Retention.TaskDays // Use task retention for context
	default:
		days = 30 // Default to 30 days
	}

	if days == 0 {
		// Return a far future date for "forever"
		return time.Now().AddDate(100, 0, 0)
	}

	return time.Now().AddDate(0, 0, days)
}
