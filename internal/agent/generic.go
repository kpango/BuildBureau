package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/pkg/types"
)

// GenericAgent is a flexible agent implementation that derives its behavior
// from configuration rather than hardcoded role-specific logic.
type GenericAgent struct {
	*BaseAgent
	config       *types.AgentConfig
	llmManager   *llm.Manager
	subordinates []types.Agent // Agents this agent can delegate to
	parent       types.Agent   // Parent agent in hierarchy
}

// NewGenericAgent creates a new generic agent with the specified configuration.
func NewGenericAgent(id string, role types.AgentRole, config *types.AgentConfig, llmManager *llm.Manager) *GenericAgent {
	return &GenericAgent{
		BaseAgent:    NewBaseAgent(id, role, config),
		config:       config,
		llmManager:   llmManager,
		subordinates: make([]types.Agent, 0),
	}
}

// SetParent sets the parent agent in the hierarchy.
func (a *GenericAgent) SetParent(parent types.Agent) {
	a.parent = parent
}

// AddSubordinate adds a subordinate agent that this agent can delegate to.
func (a *GenericAgent) AddSubordinate(subordinate types.Agent) {
	a.subordinates = append(a.subordinates, subordinate)
}

// GetSubordinates returns all subordinate agents.
func (a *GenericAgent) GetSubordinates() []types.Agent {
	return a.subordinates
}

// ProcessTask handles incoming tasks using LLM and configuration-driven behavior.
func (a *GenericAgent) ProcessTask(ctx context.Context, task *types.Task) (*types.TaskResponse, error) {
	a.IncrementActiveTasks()
	defer a.DecrementActiveTasks()

	// Store task in memory if available
	if a.memory != nil {
		tags := []string{"task", "processing", string(a.role)}
		if err := a.memory.StoreTask(ctx, task, "", tags); err != nil {
			fmt.Printf("Warning: Failed to store task in memory: %v\n", err)
		}
	}

	// Process task using LLM if available
	result, err := a.processWithLLM(ctx, task)
	if err != nil {
		return &types.TaskResponse{
			TaskID: task.ID,
			Status: types.StatusFailed,
			Error:  err.Error(),
		}, err
	}

	// Determine if delegation is needed based on configuration and LLM response
	shouldDelegate := a.shouldDelegateTask(task, result)

	if shouldDelegate && len(a.subordinates) > 0 {
		// Delegate to appropriate subordinate
		delegatedResult, err := a.delegateTask(ctx, task, result)
		if err != nil {
			return &types.TaskResponse{
				TaskID: task.ID,
				Status: types.StatusFailed,
				Error:  err.Error(),
			}, err
		}
		result += "\n" + delegatedResult
	}

	// Store result in memory
	if a.memory != nil {
		tags := []string{"task", "completed", string(a.role)}
		if err := a.memory.StoreTask(ctx, task, result, tags); err != nil {
			fmt.Printf("Warning: Failed to store result in memory: %v\n", err)
		}
	}

	return &types.TaskResponse{
		TaskID: task.ID,
		Status: types.StatusCompleted,
		Result: result,
	}, nil
}

// processWithLLM uses the LLM to process the task based on the agent's system prompt.
func (a *GenericAgent) processWithLLM(ctx context.Context, task *types.Task) (string, error) {
	if a.llmManager == nil {
		// Fallback to basic processing without LLM
		return a.processWithoutLLM(task), nil
	}

	// Build prompt with agent's role and capabilities
	prompt := a.buildPrompt(task)

	// Retrieve similar past tasks from memory for context
	var contextInfo string
	if a.memory != nil {
		similar, err := a.memory.GetRelatedTasks(ctx, task.Description, 3)
		if err == nil && len(similar) > 0 {
			contextInfo = "\n\nRelevant past experience:\n"
			for _, mem := range similar {
				contextInfo += fmt.Sprintf("- %s\n", mem.Content)
			}
		}
	}

	fullPrompt := prompt + contextInfo

	// Generate response using LLM
	response, err := a.llmManager.Generate(ctx, a.config.Model, fullPrompt, &llm.GenerateOptions{
		Temperature: 0.7,
		MaxTokens:   2000,
	})
	if err != nil {
		return "", fmt.Errorf("LLM generation failed: %w", err)
	}

	return response, nil
}

// buildPrompt constructs the prompt for the LLM based on agent configuration and task.
func (a *GenericAgent) buildPrompt(task *types.Task) string {
	var prompt strings.Builder

	// Include system prompt from configuration
	if a.config.SystemPrompt != "" {
		prompt.WriteString(a.config.SystemPrompt)
		prompt.WriteString("\n\n")
	}

	// Add task details
	prompt.WriteString(fmt.Sprintf("Task Title: %s\n", task.Title))
	prompt.WriteString(fmt.Sprintf("Task Description: %s\n", task.Description))
	if task.Content != "" {
		prompt.WriteString(fmt.Sprintf("Additional Context: %s\n", task.Content))
	}

	// Add capabilities information
	if len(a.config.Capabilities) > 0 {
		prompt.WriteString("\nYour capabilities include:\n")
		for _, cap := range a.config.Capabilities {
			prompt.WriteString(fmt.Sprintf("- %s\n", cap))
		}
	}

	// Add delegation information if subordinates exist
	if len(a.subordinates) > 0 {
		prompt.WriteString("\nYou can delegate to the following subordinates:\n")
		for _, sub := range a.subordinates {
			prompt.WriteString(fmt.Sprintf("- %s (Role: %s)\n", sub.GetID(), sub.GetRole()))
		}
		prompt.WriteString("\nIndicate in your response if delegation is needed.\n")
	}

	return prompt.String()
}

// processWithoutLLM provides basic task processing when LLM is not available.
func (a *GenericAgent) processWithoutLLM(task *types.Task) string {
	result := fmt.Sprintf("%s (ID: %s) received task: %s\n", a.role, a.GetID(), task.Title)
	result += fmt.Sprintf("Description: %s\n", task.Description)

	// Basic role-based processing
	result += a.getBasicRoleResponse(task)

	return result
}

// getBasicRoleResponse provides minimal role-appropriate response without LLM.
func (a *GenericAgent) getBasicRoleResponse(task *types.Task) string {
	switch a.role {
	case types.RolePresident:
		return "Analyzing requirements and clarifying objectives...\n"
	case types.RoleSecretary:
		return "Organizing and documenting task information...\n"
	case types.RoleDirector:
		return "Breaking down into departmental tasks...\n"
	case types.RoleManager:
		return "Creating detailed specifications...\n"
	case types.RoleEngineer:
		return "Preparing implementation approach...\n"
	default:
		return "Processing task according to configuration...\n"
	}
}

// shouldDelegateTask determines if the task should be delegated to subordinates.
func (a *GenericAgent) shouldDelegateTask(task *types.Task, llmResponse string) bool {
	if len(a.subordinates) == 0 {
		return false
	}

	// Check if LLM response suggests delegation
	lowerResponse := strings.ToLower(llmResponse)
	delegationKeywords := []string{
		"delegate",
		"assign",
		"forward",
		"pass to",
		"send to",
	}

	for _, keyword := range delegationKeywords {
		if strings.Contains(lowerResponse, keyword) {
			return true
		}
	}

	// Default delegation behavior based on role hierarchy
	// Higher-level roles typically delegate to lower-level ones
	return true
}

// delegateTask delegates the task to an appropriate subordinate.
func (a *GenericAgent) delegateTask(ctx context.Context, originalTask *types.Task, contextInfo string) (string, error) {
	// Select subordinate (round-robin for now, can be made smarter with LLM)
	subordinate := a.selectSubordinate(originalTask)
	if subordinate == nil {
		return "", fmt.Errorf("no suitable subordinate found for delegation")
	}

	// Create delegated task
	delegatedTask := &types.Task{
		ID:          uuid.New().String(),
		Title:       fmt.Sprintf("[Delegated] %s", originalTask.Title),
		Description: originalTask.Description,
		FromAgent:   a.GetID(),
		ToAgent:     subordinate.GetID(),
		Content:     contextInfo + "\n" + originalTask.Content,
		Priority:    originalTask.Priority,
		Metadata:    originalTask.Metadata,
	}

	// Process delegated task
	response, err := subordinate.ProcessTask(ctx, delegatedTask)
	if err != nil {
		return "", fmt.Errorf("delegation to %s failed: %w", subordinate.GetID(), err)
	}

	result := fmt.Sprintf("\nDelegated to %s (Role: %s)\n", subordinate.GetID(), subordinate.GetRole())
	result += fmt.Sprintf("Subordinate response: %s\n", response.Result)

	return result, nil
}

// selectSubordinate selects an appropriate subordinate for task delegation.
// Currently uses round-robin, but can be enhanced with load balancing or LLM-based selection.
func (a *GenericAgent) selectSubordinate(task *types.Task) types.Agent {
	if len(a.subordinates) == 0 {
		return nil
	}

	// Simple round-robin for now
	// In a production system, this could consider:
	// - Current load of each subordinate
	// - Capabilities matching task requirements
	// - Past performance on similar tasks
	// - LLM-based decision making

	// For now, return the first subordinate
	// This can be enhanced to use task.ToAgent if specified
	if task.ToAgent != "" {
		for _, sub := range a.subordinates {
			if sub.GetID() == task.ToAgent {
				return sub
			}
		}
	}

	return a.subordinates[0]
}
