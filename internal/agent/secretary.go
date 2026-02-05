package agent

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/kpango/BuildBureau/pkg/types"
)

// SecretaryAgent represents a secretary agent attached to leadership roles.
type SecretaryAgent struct {
	*BaseAgent
	attachedTo      types.Agent
	directors       []types.Agent
	nextDirectorIdx uint32
}

// NewSecretaryAgent creates a new Secretary agent.
func NewSecretaryAgent(id string, config *types.AgentConfig) *SecretaryAgent {
	return &SecretaryAgent{
		BaseAgent: NewBaseAgent(id, types.RoleSecretary, config),
		directors: make([]types.Agent, 0),
	}
}

// AttachTo sets the leadership agent this secretary serves.
func (a *SecretaryAgent) AttachTo(leader types.Agent) {
	a.attachedTo = leader
}

// AddDirector adds a director to delegate tasks to.
func (a *SecretaryAgent) AddDirector(director types.Agent) {
	a.directors = append(a.directors, director)
}

// ProcessTask handles incoming tasks for the Secretary.
func (a *SecretaryAgent) ProcessTask(ctx context.Context, task *types.Task) (*types.TaskResponse, error) {
	a.IncrementActiveTasks()
	defer a.DecrementActiveTasks()

	// Store conversation memory if memory is enabled
	if mem := a.GetMemory(); mem != nil {
		_ = mem.StoreConversation(ctx, fmt.Sprintf("Received task: %s - %s", task.Title, task.Description), []string{"secretary", "delegation"})
	}

	result := fmt.Sprintf("Secretary %s processing task from %s\n", a.GetID(), task.FromAgent)
	result += "Recording goal and decisions...\n"

	// If we have directors, delegate to them using round-robin with memory-informed selection
	if len(a.directors) > 0 {
		result += fmt.Sprintf("Delegating to %d Director(s)...\n", len(a.directors))

		// Check past delegation performance from memory
		selectedDirector := a.selectDirectorWithMemory(ctx, task)

		directorTask := &types.Task{
			ID:          uuid.New().String(),
			Title:       "Director: " + task.Title,
			Description: task.Description,
			FromAgent:   a.GetID(),
			ToAgent:     selectedDirector.GetID(),
			Content:     task.Content,
			Priority:    task.Priority,
		}

		// Store delegation decision in memory
		if mem := a.GetMemory(); mem != nil {
			decision := fmt.Sprintf("Delegated to director %s", selectedDirector.GetID())
			reasoning := "Selected based on round-robin and past performance"
			_ = mem.StoreDecision(ctx, decision, reasoning, []string{"delegation", "director"})
		}

		response, err := selectedDirector.ProcessTask(ctx, directorTask)
		if err != nil {
			return nil, fmt.Errorf("failed to delegate to director: %w", err)
		}

		if response.Status == types.StatusFailed {
			return nil, fmt.Errorf("director task failed: %s", response.Error)
		}

		result += fmt.Sprintf("Director response: %s\n", response.Result)

		// Store task completion memory
		if mem := a.GetMemory(); mem != nil {
			_ = mem.StoreTask(ctx, task, result, []string{"secretary", "completed", "delegated"})
		}
	} else {
		result += "No directors available. Task recorded and completed at Secretary level.\n"

		// Store task completion memory
		if mem := a.GetMemory(); mem != nil {
			_ = mem.StoreTask(ctx, task, result, []string{"secretary", "completed", "no-delegation"})
		}
	}

	return &types.TaskResponse{
		TaskID: task.ID,
		Status: types.StatusCompleted,
		Result: result,
	}, nil
}

// selectDirectorWithMemory selects the best director based on round-robin and memory.
func (a *SecretaryAgent) selectDirectorWithMemory(ctx context.Context, task *types.Task) types.Agent {
	// Default round-robin selection
	idx := atomic.AddUint32(&a.nextDirectorIdx, 1) - 1
	selectedIdx := int(idx) % len(a.directors)

	// Try to use memory to inform selection
	if mem := a.GetMemory(); mem != nil {
		// Look for similar past tasks
		relatedTasks, err := mem.GetRelatedTasks(ctx, task.Description, 5)
		if err == nil && len(relatedTasks) > 0 {
			// Check which director handled similar tasks successfully
			directorPerformance := make(map[string]int)
			for _, memory := range relatedTasks {
				if toAgent, ok := memory.Metadata["to_agent"]; ok {
					directorPerformance[toAgent]++
				}
			}

			// Find if any of our current directors had good performance
			for i, director := range a.directors {
				if count, ok := directorPerformance[director.GetID()]; ok && count > 0 {
					selectedIdx = i
					break
				}
			}
		}
	}

	return a.directors[selectedIdx]
}
