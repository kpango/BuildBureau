package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/kpango/BuildBureau/pkg/types"
)

// SecretaryAgent assists higher-level agents with documentation and research
type SecretaryAgent struct {
	*BaseAgent
	boss Agent
}

// NewSecretaryAgent creates a new secretary agent
func NewSecretaryAgent(name string, boss Agent, eventChan chan types.AgentEvent) *SecretaryAgent {
	return &SecretaryAgent{
		BaseAgent: NewBaseAgent(types.RoleSecretary, name, eventChan),
		boss:      boss,
	}
}

// HandleTask processes documentation and research tasks
func (s *SecretaryAgent) HandleTask(ctx context.Context, task types.Task) (types.Task, error) {
	s.EmitEvent(types.EventTaskStarted, fmt.Sprintf("Secretary assisting with: %s", task.Title), task.ID)

	// Add task to tracking
	task.AssignedTo = s.Role()
	task.Status = types.StatusInProgress
	task.UpdatedAt = time.Now()
	s.AddTask(task)

	// Secretary performs documentation, research, and knowledge base management
	s.EmitEvent(types.EventMessage, "Secretary documenting and researching task requirements", task.ID)

	// Simulate research and documentation work
	time.Sleep(300 * time.Millisecond)

	// Record to knowledge base (simulated)
	s.EmitEvent(types.EventMessage, "Secretary updating knowledge base with findings", task.ID)

	// Complete the task
	task.Status = types.StatusCompleted
	task.Result = fmt.Sprintf("Documentation and research completed for: %s", task.Title)
	task.UpdatedAt = time.Now()
	s.UpdateTask(task.ID, task.Status, task.Result)

	s.EmitEvent(types.EventTaskCompleted, fmt.Sprintf("Secretary completed assistance: %s", task.Title), task.ID)

	return task, nil
}
