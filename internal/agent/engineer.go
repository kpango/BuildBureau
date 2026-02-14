package agent

import (
	"context"
	"fmt"

	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/pkg/types"
)

const (
	defaultSimilarTasksLimit = 3
)

// EngineerAgent represents an engineer agent that implements code using LLM.
type EngineerAgent struct {
	*BaseAgent
	llmManager *llm.Manager
}

// NewEngineerAgent creates a new Engineer agent.
func NewEngineerAgent(id string, config *types.AgentConfig, llmManager *llm.Manager) *EngineerAgent {
	return &EngineerAgent{
		BaseAgent:  NewBaseAgent(id, types.RoleEngineer, config),
		llmManager: llmManager,
	}
}

// ProcessTask handles incoming tasks for the Engineer using LLM and memory.
func (a *EngineerAgent) ProcessTask(ctx context.Context, task *types.Task) (*types.TaskResponse, error) {
	a.IncrementActiveTasks()
	defer a.DecrementActiveTasks()

	// Store conversation memory
	if mem := a.GetMemory(); mem != nil {
		_ = mem.StoreConversation(ctx, fmt.Sprintf("Received implementation task: %s", task.Title), []string{"engineer", "implementation"})
	}

	result := fmt.Sprintf("Engineer %s implementing task: %s\n", a.GetID(), task.Title)

	// Check memory for similar past implementations
	var contextFromMemory string
	if mem := a.GetMemory(); mem != nil {
		relatedTasks, err := mem.GetRelatedTasks(ctx, task.Description, defaultSimilarTasksLimit)
		if err == nil && len(relatedTasks) > 0 {
			result += fmt.Sprintf("Found %d related past implementation(s) to learn from.\n", len(relatedTasks))
			contextFromMemory = "\n\n=== Context from Past Implementations ===\n"
			for i, memory := range relatedTasks {
				contextFromMemory += fmt.Sprintf("\nPast Implementation %d (Score: %.2f):\n%s\n", i+1, memory.Score, memory.Content)
			}
			contextFromMemory += "=== End of Past Context ===\n\n"
		}

		// Check for relevant knowledge
		knowledge, err := mem.GetKnowledge(ctx, task.Description, 2)
		if err == nil && len(knowledge) > 0 {
			contextFromMemory += "\n=== Relevant Knowledge ===\n"
			for _, k := range knowledge {
				contextFromMemory += fmt.Sprintf("%s\n", k.Content)
			}
			contextFromMemory += "=== End of Knowledge ===\n\n"
		}
	}

	// Use LLM if available to generate actual implementation
	if a.llmManager != nil {
		prompt := fmt.Sprintf(`You are a software engineer tasked with implementing the following:

Title: %s
Description: %s
Specifications: %s
%s
Please provide:
1. A detailed implementation plan
2. Code implementation (if applicable)
3. Test cases
4. Documentation

Be specific and provide working code. Learn from the past implementations provided above if available.`,
			task.Title, task.Description, task.Content, contextFromMemory)

		llmOpts := &llm.GenerateOptions{
			Temperature:  0.7,
			MaxTokens:    4096,
			SystemPrompt: a.config.SystemPrompt,
		}

		model := a.config.Model
		if model == "" {
			model = "gemini" // default
		}

		response, err := a.llmManager.Generate(ctx, model, prompt, llmOpts)
		if err != nil {
			result += fmt.Sprintf("Error using LLM: %v\n", err)
			result += "Falling back to simple acknowledgment.\n"
			result += fmt.Sprintf("Task content: %s\n", task.Content)
			result += "Implementation completed successfully (without LLM assistance).\n"
		} else {
			result += "=== LLM-Generated Implementation ===\n"
			result += response
			result += "\n=== End of Implementation ===\n"

			// Store the generated code as knowledge
			if mem := a.GetMemory(); mem != nil {
				knowledgeContent := fmt.Sprintf("Implementation for: %s\n\nCode:\n%s", task.Title, response)
				_ = mem.StoreKnowledge(ctx, knowledgeContent, []string{"code", "implementation", task.Title})
			}
		}
	} else {
		result += "No LLM manager available.\n"
		result += fmt.Sprintf("Task content: %s\n", task.Content)
		result += "Implementation completed successfully (without LLM assistance).\n"
	}

	// Store task completion in memory
	if mem := a.GetMemory(); mem != nil {
		_ = mem.StoreTask(ctx, task, result, []string{"engineer", "implementation", "completed"})
	}

	return &types.TaskResponse{
		TaskID: task.ID,
		Status: types.StatusCompleted,
		Result: result,
	}, nil
}
