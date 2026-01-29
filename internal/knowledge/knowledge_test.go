package knowledge

import (
	"context"
	"testing"
)

func TestInMemoryKB_Store(t *testing.T) {
	kb := NewInMemoryKB()
	ctx := context.Background()

	err := kb.Store(ctx, "project1", "Web application", map[string]string{"type": "project"}, "president-1")
	if err != nil {
		t.Fatalf("Failed to store: %v", err)
	}

	entry, err := kb.Get(ctx, "project1")
	if err != nil {
		t.Fatalf("Failed to get: %v", err)
	}

	if entry.Value != "Web application" {
		t.Errorf("Expected 'Web application', got '%s'", entry.Value)
	}
	if entry.CreatedBy != "president-1" {
		t.Errorf("Expected creator 'president-1', got '%s'", entry.CreatedBy)
	}
}

func TestInMemoryKB_Update(t *testing.T) {
	kb := NewInMemoryKB()
	ctx := context.Background()

	// Store initial value
	kb.Store(ctx, "status", "planning", nil, "agent-1")

	// Update
	kb.Store(ctx, "status", "in-progress", map[string]string{"phase": "2"}, "agent-2")

	entry, err := kb.Get(ctx, "status")
	if err != nil {
		t.Fatalf("Failed to get: %v", err)
	}

	if entry.Value != "in-progress" {
		t.Errorf("Expected 'in-progress', got '%s'", entry.Value)
	}
	if entry.Metadata["phase"] != "2" {
		t.Errorf("Expected metadata phase '2', got '%s'", entry.Metadata["phase"])
	}
}

func TestInMemoryKB_Search(t *testing.T) {
	kb := NewInMemoryKB()
	ctx := context.Background()

	kb.Store(ctx, "task1", "Build login feature", nil, "agent-1")
	kb.Store(ctx, "task2", "Build dashboard", nil, "agent-1")
	kb.Store(ctx, "task3", "Write tests", nil, "agent-2")

	results, err := kb.Search(ctx, "Build")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestInMemoryKB_Delete(t *testing.T) {
	kb := NewInMemoryKB()
	ctx := context.Background()

	kb.Store(ctx, "temp", "temporary data", nil, "agent-1")

	err := kb.Delete(ctx, "temp")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = kb.Get(ctx, "temp")
	if err == nil {
		t.Error("Expected error when getting deleted entry")
	}
}

func TestInMemoryKB_List(t *testing.T) {
	kb := NewInMemoryKB()
	ctx := context.Background()

	kb.Store(ctx, "key1", "value1", nil, "agent-1")
	kb.Store(ctx, "key2", "value2", nil, "agent-1")
	kb.Store(ctx, "key3", "value3", nil, "agent-2")

	entries, err := kb.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}
}

func TestInMemoryKB_GetNonExistent(t *testing.T) {
	kb := NewInMemoryKB()
	ctx := context.Background()

	_, err := kb.Get(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error when getting nonexistent key")
	}
}
