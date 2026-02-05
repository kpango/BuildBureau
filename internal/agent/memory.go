package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/kpango/BuildBureau/pkg/types"
)

// AgentMemory provides memory functionality for agents.
type AgentMemory struct {
	manager types.MemoryManager
	agentID string
	enabled bool
}

// NewAgentMemory creates a new agent memory instance.
func NewAgentMemory(agentID string, manager types.MemoryManager) *AgentMemory {
	return &AgentMemory{
		manager: manager,
		agentID: agentID,
		enabled: manager != nil,
	}
}

// StoreConversation stores a conversation memory.
func (m *AgentMemory) StoreConversation(ctx context.Context, content string, tags []string) error {
	if !m.enabled {
		return nil // Silently skip if memory not enabled
	}

	entry := &types.MemoryEntry{
		AgentID: m.agentID,
		Type:    types.MemoryTypeConversation,
		Content: content,
		Tags:    tags,
		Metadata: map[string]string{
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		},
	}

	return m.manager.StoreMemory(ctx, entry)
}

// StoreTask stores a task-related memory.
func (m *AgentMemory) StoreTask(ctx context.Context, task *types.Task, result string, tags []string) error {
	if !m.enabled {
		return nil
	}

	content := fmt.Sprintf("Task: %s\nDescription: %s\nResult: %s", task.Title, task.Description, result)

	entry := &types.MemoryEntry{
		AgentID: m.agentID,
		Type:    types.MemoryTypeTask,
		Content: content,
		Tags:    tags,
		Metadata: map[string]string{
			"task_id":    task.ID,
			"from_agent": task.FromAgent,
			"to_agent":   task.ToAgent,
			"priority":   fmt.Sprintf("%d", task.Priority),
			"timestamp":  fmt.Sprintf("%d", time.Now().Unix()),
		},
	}

	return m.manager.StoreMemory(ctx, entry)
}

// StoreKnowledge stores learned knowledge.
func (m *AgentMemory) StoreKnowledge(ctx context.Context, content string, tags []string) error {
	if !m.enabled {
		return nil
	}

	entry := &types.MemoryEntry{
		AgentID: m.agentID,
		Type:    types.MemoryTypeKnowledge,
		Content: content,
		Tags:    tags,
		Metadata: map[string]string{
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		},
	}

	return m.manager.StoreMemory(ctx, entry)
}

// StoreDecision stores a decision made by the agent.
func (m *AgentMemory) StoreDecision(ctx context.Context, decision string, reasoning string, tags []string) error {
	if !m.enabled {
		return nil
	}

	content := fmt.Sprintf("Decision: %s\nReasoning: %s", decision, reasoning)

	entry := &types.MemoryEntry{
		AgentID: m.agentID,
		Type:    types.MemoryTypeDecision,
		Content: content,
		Tags:    tags,
		Metadata: map[string]string{
			"decision":  decision,
			"reasoning": reasoning,
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		},
	}

	return m.manager.StoreMemory(ctx, entry)
}

// GetConversationHistory retrieves recent conversation history.
func (m *AgentMemory) GetConversationHistory(ctx context.Context, limit int) ([]*types.MemoryEntry, error) {
	if !m.enabled {
		return nil, nil
	}

	return m.manager.GetConversationHistory(ctx, m.agentID, limit)
}

// GetRelatedTasks finds tasks similar to the given query.
func (m *AgentMemory) GetRelatedTasks(ctx context.Context, query string, limit int) ([]*types.MemoryEntry, error) {
	if !m.enabled {
		return nil, nil
	}

	// Try semantic search first, fall back to text search
	results, err := m.manager.SemanticSearch(ctx, query, m.agentID, limit)
	if err != nil {
		// Fallback to basic query
		return m.manager.QueryMemories(ctx, &types.MemoryQuery{
			AgentID: m.agentID,
			Type:    types.MemoryTypeTask,
			Content: query,
			Limit:   limit,
		})
	}

	return results, nil
}

// GetKnowledge retrieves relevant knowledge based on query.
func (m *AgentMemory) GetKnowledge(ctx context.Context, query string, limit int) ([]*types.MemoryEntry, error) {
	if !m.enabled {
		return nil, nil
	}

	return m.manager.QueryMemories(ctx, &types.MemoryQuery{
		AgentID: m.agentID,
		Type:    types.MemoryTypeKnowledge,
		Content: query,
		Limit:   limit,
	})
}

// GetDecisionHistory retrieves past decisions.
func (m *AgentMemory) GetDecisionHistory(ctx context.Context, limit int) ([]*types.MemoryEntry, error) {
	if !m.enabled {
		return nil, nil
	}

	return m.manager.QueryMemories(ctx, &types.MemoryQuery{
		AgentID: m.agentID,
		Type:    types.MemoryTypeDecision,
		Limit:   limit,
	})
}

// SearchMemory performs a semantic search across all memory types.
func (m *AgentMemory) SearchMemory(ctx context.Context, query string, limit int) ([]*types.MemoryEntry, error) {
	if !m.enabled {
		return nil, nil
	}

	return m.manager.SemanticSearch(ctx, query, m.agentID, limit)
}
