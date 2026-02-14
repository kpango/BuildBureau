package agent

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/kpango/BuildBureau/pkg/types"
)

// DirectorAgent represents a director agent that decomposes projects.
type DirectorAgent struct {
	*BaseAgent
	secretary      types.Agent
	managers       []types.Agent
	nextManagerIdx uint32
}

// NewDirectorAgent creates a new Director agent.
func NewDirectorAgent(id string, config *types.AgentConfig) *DirectorAgent {
	return &DirectorAgent{
		BaseAgent: NewBaseAgent(id, types.RoleDirector, config),
		managers:  make([]types.Agent, 0),
	}
}

// SetSecretary assigns a secretary to the director.
func (a *DirectorAgent) SetSecretary(secretary types.Agent) {
	a.secretary = secretary
}

// AddManager adds a manager to delegate tasks to.
func (a *DirectorAgent) AddManager(manager types.Agent) {
	a.managers = append(a.managers, manager)
}

// ProcessTask handles incoming tasks for the Director.
func (a *DirectorAgent) ProcessTask(ctx context.Context, task *types.Task) (*types.TaskResponse, error) {
	a.IncrementActiveTasks()
	defer a.DecrementActiveTasks()

	result := fmt.Sprintf("Director %s processing task: %s\n", a.GetID(), task.Title)
	result += "Performing research and expanding requirements...\n"
	result += "Decomposing project into department-level tasks...\n"

	// If we have managers, delegate to them using round-robin
	if len(a.managers) > 0 {
		result += fmt.Sprintf("Delegating to %d Manager(s)...\n", len(a.managers))

		// Round-robin selection
		idx := atomic.AddUint32(&a.nextManagerIdx, 1) - 1
		manager := a.managers[int(idx)%len(a.managers)]

		managerTask := &types.Task{
			ID:          uuid.New().String(),
			Title:       "Manager: " + task.Title,
			Description: task.Description,
			FromAgent:   a.GetID(),
			ToAgent:     manager.GetID(),
			Content:     task.Content,
			Priority:    task.Priority,
		}

		response, err := manager.ProcessTask(ctx, managerTask)
		if err != nil {
			return nil, fmt.Errorf("failed to delegate to manager: %w", err)
		}

		if response.Status == types.StatusFailed {
			return nil, fmt.Errorf("manager task failed: %s", response.Error)
		}

		result += fmt.Sprintf("Manager response: %s\n", response.Result)
	} else {
		result += "No managers available. Task completed at Director level.\n"
	}

	return &types.TaskResponse{
		TaskID: task.ID,
		Status: types.StatusCompleted,
		Result: result,
	}, nil
}
