package agent

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/pkg/types"
)

// ManagerAgent represents a manager agent that produces software designs.
type ManagerAgent struct {
	secretary types.Agent
	*BaseAgent
	llmManager      *llm.Manager
	engineers       []types.Agent
	nextEngineerIdx uint32
}

// NewManagerAgent creates a new Manager agent.
func NewManagerAgent(id string, config *types.AgentConfig, llmManager *llm.Manager) *ManagerAgent {
	return &ManagerAgent{
		BaseAgent:  NewBaseAgent(id, types.RoleManager, config),
		engineers:  make([]types.Agent, 0),
		llmManager: llmManager,
	}
}

// SetSecretary assigns a secretary to the manager.
func (a *ManagerAgent) SetSecretary(secretary types.Agent) {
	a.secretary = secretary
}

// AddEngineer adds an engineer to delegate tasks to.
func (a *ManagerAgent) AddEngineer(engineer types.Agent) {
	a.engineers = append(a.engineers, engineer)
}

// ProcessTask handles incoming tasks for the Manager using LLM and memory.
func (a *ManagerAgent) ProcessTask(ctx context.Context, task *types.Task) (*types.TaskResponse, error) {
	a.IncrementActiveTasks()
	defer a.DecrementActiveTasks()

	// Store conversation memory
	if mem := a.GetMemory(); mem != nil {
		_ = mem.StoreConversation(ctx, fmt.Sprintf("Received design task: %s", task.Title), []string{"manager", "design"})
	}

	result := fmt.Sprintf("Manager %s processing task: %s\n", a.GetID(), task.Title)

	// Check memory for similar past designs
	var contextFromMemory string
	if mem := a.GetMemory(); mem != nil {
		relatedTasks, err := mem.GetRelatedTasks(ctx, task.Description, 3)
		if err == nil && len(relatedTasks) > 0 {
			result += fmt.Sprintf("Found %d related past design(s) to reference.\n", len(relatedTasks))
			contextFromMemory = "\n\n=== Context from Past Designs ===\n"
			for i, memory := range relatedTasks {
				contextFromMemory += fmt.Sprintf("\nPast Design %d:\n%s\n", i+1, memory.Content)
			}
			contextFromMemory += "=== End of Past Context ===\n\n"
		}
	}

	// Use LLM if available to create software design
	var designSpec string
	if a.llmManager != nil {
		prompt := fmt.Sprintf(`You are a software manager tasked with creating a detailed technical specification for:

Title: %s
Description: %s
Requirements: %s
%s
Please provide:
1. High-level architecture design
2. Component breakdown
3. Technical specifications for each component
4. Interface definitions
5. Implementation guidelines for engineers

Be detailed and technical. Learn from the past designs provided above if available.`,
			task.Title, task.Description, task.Content, contextFromMemory)

		llmOpts := &llm.GenerateOptions{
			Temperature:  0.5, // Lower temperature for more focused technical output
			MaxTokens:    3072,
			SystemPrompt: a.config.SystemPrompt,
		}

		model := a.config.Model
		if model == "" {
			model = "gemini"
		}

		response, err := a.llmManager.Generate(ctx, model, prompt, llmOpts)
		if err != nil {
			result += fmt.Sprintf("Warning: LLM generation failed: %v\n", err)
			designSpec = fmt.Sprintf("Specifications for: %s\n", task.Content)
		} else {
			result += "=== LLM-Generated Design Specification ===\n"
			result += response
			result += "\n=== End of Specification ===\n"
			designSpec = response

			// Store the design as knowledge
			if mem := a.GetMemory(); mem != nil {
				knowledgeContent := fmt.Sprintf("Design for: %s\n\nSpecification:\n%s", task.Title, response)
				_ = mem.StoreKnowledge(ctx, knowledgeContent, []string{"design", "specification", task.Title})
			}
		}
	} else {
		designSpec = fmt.Sprintf("Specifications for: %s\n", task.Content)
	}

	// If we have engineers, delegate to them using round-robin with memory
	if len(a.engineers) > 0 {
		result += fmt.Sprintf("\nDelegating implementation to %d Engineer(s)...\n", len(a.engineers))

		// Round-robin selection
		idx := atomic.AddUint32(&a.nextEngineerIdx, 1) - 1
		engineer := a.engineers[int(idx)%len(a.engineers)]

		// Store delegation decision
		if mem := a.GetMemory(); mem != nil {
			decision := fmt.Sprintf("Delegated to engineer %s", engineer.GetID())
			reasoning := "Selected based on round-robin"
			_ = mem.StoreDecision(ctx, decision, reasoning, []string{"delegation", "engineer"})
		}

		engineerTask := &types.Task{
			ID:          uuid.New().String(),
			Title:       "Engineer: " + task.Title,
			Description: task.Description,
			FromAgent:   a.GetID(),
			ToAgent:     engineer.GetID(),
			Content:     designSpec, // Pass the design spec to the engineer
			Priority:    task.Priority,
		}

		response, err := engineer.ProcessTask(ctx, engineerTask)
		if err != nil {
			return nil, fmt.Errorf("failed to delegate to engineer: %w", err)
		}

		if response.Status == types.StatusFailed {
			return nil, fmt.Errorf("engineer task failed: %s", response.Error)
		}

		result += fmt.Sprintf("Engineer response: %s\n", response.Result)
	} else {
		result += "No engineers available. Design completed at Manager level.\n"
	}

	// Store task completion in memory
	if mem := a.GetMemory(); mem != nil {
		_ = mem.StoreTask(ctx, task, result, []string{"manager", "design", "completed"})
	}

	return &types.TaskResponse{
		TaskID: task.ID,
		Status: types.StatusCompleted,
		Result: result,
	}, nil
}
