package agent

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kpango/BuildBureau/pkg/types"
)

// PresidentAgent represents the top-level agent that interacts with clients.
type PresidentAgent struct {
	*BaseAgent
	secretary types.Agent
}

// NewPresidentAgent creates a new President agent.
func NewPresidentAgent(id string, config *types.AgentConfig) *PresidentAgent {
	return &PresidentAgent{
		BaseAgent: NewBaseAgent(id, types.RolePresident, config),
	}
}

// SetSecretary assigns a secretary to the president.
func (a *PresidentAgent) SetSecretary(secretary types.Agent) {
	a.secretary = secretary
}

// ProcessTask handles incoming tasks for the President.
func (a *PresidentAgent) ProcessTask(ctx context.Context, task *types.Task) (*types.TaskResponse, error) {
	a.IncrementActiveTasks()
	defer a.DecrementActiveTasks()

	// President clarifies client instructions and summarizes objectives
	result := fmt.Sprintf("President %s received task: %s\n", a.GetID(), task.Title)
	result += "Clarifying requirements and defining high-level objectives...\n"
	result += fmt.Sprintf("Task: %s\n", task.Description)

	// Delegate to secretary if available
	if a.secretary != nil {
		result += "Delegating to Secretary...\n"
		secretaryTask := &types.Task{
			ID:          uuid.New().String(),
			Title:       "Secretary: " + task.Title,
			Description: task.Description,
			FromAgent:   a.GetID(),
			ToAgent:     a.secretary.GetID(),
			Content:     task.Content,
			Priority:    task.Priority,
		}

		response, err := a.secretary.ProcessTask(ctx, secretaryTask)
		if err != nil {
			return nil, fmt.Errorf("failed to delegate to secretary: %w", err)
		}

		if response.Status == types.StatusFailed {
			return nil, fmt.Errorf("secretary task failed: %s", response.Error)
		}

		result += fmt.Sprintf("Secretary response: %s\n", response.Result)

		return &types.TaskResponse{
			TaskID: task.ID,
			Status: types.StatusCompleted,
			Result: result,
		}, nil
	}

	return &types.TaskResponse{
		TaskID: task.ID,
		Status: types.StatusCompleted,
		Result: result + "No secretary assigned, task completed at President level.\n",
	}, nil
}
