package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/kpango/BuildBureau/pkg/types"
)

// BaseAgent provides common functionality for all agent types.
type BaseAgent struct {
	config         *types.AgentConfig
	memory         *AgentMemory
	id             string
	role           types.AgentRole
	activeTasks    int
	completedTasks int
	mu             sync.RWMutex
	running        bool
}

// NewBaseAgent creates a new base agent.
func NewBaseAgent(id string, role types.AgentRole, config *types.AgentConfig) *BaseAgent {
	return &BaseAgent{
		id:     id,
		role:   role,
		config: config,
		memory: nil, // Will be set by SetMemoryManager
	}
}

// SetMemoryManager sets the memory manager for this agent.
func (a *BaseAgent) SetMemoryManager(manager types.MemoryManager) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.memory = NewAgentMemory(a.id, manager)
}

// GetMemory returns the agent's memory interface.
func (a *BaseAgent) GetMemory() *AgentMemory {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.memory
}

// GetID returns the agent's unique identifier.
func (a *BaseAgent) GetID() string {
	return a.id
}

// GetRole returns the agent's role.
func (a *BaseAgent) GetRole() types.AgentRole {
	return a.role
}

// Start initializes the agent.
func (a *BaseAgent) Start(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.running {
		return fmt.Errorf("agent %s is already running", a.id)
	}

	a.running = true
	return nil
}

// Stop gracefully shuts down the agent.
func (a *BaseAgent) Stop(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.running {
		return fmt.Errorf("agent %s is not running", a.id)
	}

	a.running = false
	return nil
}

// IsRunning returns whether the agent is currently running.
func (a *BaseAgent) IsRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.running
}

// GetStats returns the agent's statistics.
func (a *BaseAgent) GetStats() (active int, completed int) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.activeTasks, a.completedTasks
}

// IncrementActiveTasks increments the active task counter.
func (a *BaseAgent) IncrementActiveTasks() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.activeTasks++
}

// DecrementActiveTasks decrements the active task counter and increments completed.
func (a *BaseAgent) DecrementActiveTasks() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.activeTasks--
	a.completedTasks++
}
