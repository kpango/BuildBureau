package memory

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kpango/BuildBureau/pkg/types"
)

func TestSQLiteStore(t *testing.T) {
	// Create in-memory SQLite store
	config := types.SQLiteConfig{
		Enabled:  true,
		InMemory: true,
	}

	store, err := NewSQLiteStore(config)
	if err != nil {
		t.Fatalf("Failed to create SQLite store: %v", err)
	}
	defer store.Close()

	ctx := context.Background()

	// Test Store
	t.Run("Store", func(t *testing.T) {
		entry := &types.MemoryEntry{
			ID:      "test-1",
			AgentID: "agent-1",
			Type:    types.MemoryTypeConversation,
			Content: "Hello, this is a test conversation",
			Metadata: map[string]string{
				"key1": "value1",
			},
			Tags:      []string{"test", "conversation"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := store.Store(ctx, entry)
		if err != nil {
			t.Errorf("Failed to store entry: %v", err)
		}
	})

	// Test Retrieve
	t.Run("Retrieve", func(t *testing.T) {
		entry, err := store.Retrieve(ctx, "test-1")
		if err != nil {
			t.Errorf("Failed to retrieve entry: %v", err)
		}

		if entry.ID != "test-1" {
			t.Errorf("Expected ID test-1, got %s", entry.ID)
		}

		if entry.Content != "Hello, this is a test conversation" {
			t.Errorf("Content mismatch")
		}

		if entry.Metadata["key1"] != "value1" {
			t.Errorf("Metadata mismatch")
		}
	})

	// Test Query
	t.Run("Query", func(t *testing.T) {
		// Store more entries
		for i := 2; i <= 5; i++ {
			entry := &types.MemoryEntry{
				ID:        fmt.Sprintf("test-%d", i),
				AgentID:   "agent-1",
				Type:      types.MemoryTypeTask,
				Content:   fmt.Sprintf("Task %d", i),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			store.Store(ctx, entry)
		}

		query := &types.MemoryQuery{
			AgentID: "agent-1",
			Type:    types.MemoryTypeTask,
			Limit:   10,
		}

		entries, err := store.Query(ctx, query)
		if err != nil {
			t.Errorf("Failed to query: %v", err)
		}

		if len(entries) != 4 {
			t.Errorf("Expected 4 entries, got %d", len(entries))
		}
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		entry, _ := store.Retrieve(ctx, "test-1")
		entry.Content = "Updated content"

		err := store.Update(ctx, entry)
		if err != nil {
			t.Errorf("Failed to update: %v", err)
		}

		updated, _ := store.Retrieve(ctx, "test-1")
		if updated.Content != "Updated content" {
			t.Errorf("Update failed, content not changed")
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err := store.Delete(ctx, "test-2")
		if err != nil {
			t.Errorf("Failed to delete: %v", err)
		}

		_, err = store.Retrieve(ctx, "test-2")
		if err == nil {
			t.Errorf("Entry should have been deleted")
		}
	})

	// Test DeleteExpired
	t.Run("DeleteExpired", func(t *testing.T) {
		// Create an expired entry
		expiredTime := time.Now().Add(-1 * time.Hour)
		entry := &types.MemoryEntry{
			ID:        "expired-1",
			AgentID:   "agent-1",
			Type:      types.MemoryTypeConversation,
			Content:   "Expired content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ExpiresAt: &expiredTime,
		}
		store.Store(ctx, entry)

		count, err := store.DeleteExpired(ctx)
		if err != nil {
			t.Errorf("Failed to delete expired: %v", err)
		}

		if count < 1 {
			t.Errorf("Expected at least 1 expired entry to be deleted, got %d", count)
		}
	})
}

func TestMemoryManager(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping memory manager integration test")
	}

	config := &types.MemoryConfig{
		Enabled: true,
		SQLite: types.SQLiteConfig{
			Enabled:  true,
			InMemory: true,
		},
		Vald: types.ValdConfig{
			Enabled: false, // Vald requires external service
		},
		Retention: types.RetentionConfig{
			ConversationDays: 30,
			TaskDays:         60,
			KnowledgeDays:    0, // Forever
			MaxEntries:       1000,
		},
	}

	manager, err := NewManager(config, nil)
	if err != nil {
		t.Fatalf("Failed to create memory manager: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()

	// Test StoreMemory
	t.Run("StoreMemory", func(t *testing.T) {
		entry := &types.MemoryEntry{
			AgentID: "agent-1",
			Type:    types.MemoryTypeConversation,
			Content: "This is a conversation memory",
			Metadata: map[string]string{
				"task_id": "task-123",
			},
			Tags: []string{"important", "conversation"},
		}

		err := manager.StoreMemory(ctx, entry)
		if err != nil {
			t.Errorf("Failed to store memory: %v", err)
		}

		if entry.ID == "" {
			t.Errorf("Memory ID should be generated")
		}
	})

	// Test QueryMemories
	t.Run("QueryMemories", func(t *testing.T) {
		query := &types.MemoryQuery{
			AgentID: "agent-1",
			Type:    types.MemoryTypeConversation,
			Limit:   10,
		}

		memories, err := manager.QueryMemories(ctx, query)
		if err != nil {
			t.Errorf("Failed to query memories: %v", err)
		}

		if len(memories) == 0 {
			t.Errorf("Expected at least one memory")
		}
	})

	// Test GetConversationHistory
	t.Run("GetConversationHistory", func(t *testing.T) {
		history, err := manager.GetConversationHistory(ctx, "agent-1", 10)
		if err != nil {
			t.Errorf("Failed to get conversation history: %v", err)
		}

		if len(history) == 0 {
			t.Errorf("Expected conversation history")
		}
	})
}
