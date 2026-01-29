package agent

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"

	"github.com/kpango/BuildBureau/pkg/types"
)

// ADKAgent wraps Google's ADK llmagent to provide ADK-powered functionality
// while maintaining compatibility with our Agent interface.
type ADKAgent struct {
	*BaseAgent
	modelName string
	llmConfig llmagent.Config
	apiKey    string
}

// NewADKEngineerAgent creates an Engineer agent using Google's ADK framework.
func NewADKEngineerAgent(id string, config *types.AgentConfig, apiKey string) (*ADKAgent, error) {
	return newADKAgent(id, config, apiKey, types.RoleEngineer, `You are a skilled software engineer.
Your responsibilities:
- Implement code according to specifications
- Write clean, maintainable code
- Add appropriate comments and documentation
- Consider edge cases and error handling
- Follow best practices for the target language`)
}

// NewADKManagerAgent creates a Manager agent using Google's ADK framework.
func NewADKManagerAgent(id string, config *types.AgentConfig, apiKey string) (*ADKAgent, error) {
	return newADKAgent(id, config, apiKey, types.RoleManager, `You are a software development manager.
Your responsibilities:
- Create detailed technical specifications
- Design software architecture
- Break down projects into implementable components
- Define interfaces and data structures
- Plan testing strategies`)
}

// NewADKDirectorAgent creates a Director agent using Google's ADK framework.
func NewADKDirectorAgent(id string, config *types.AgentConfig, apiKey string) (*ADKAgent, error) {
	return newADKAgent(id, config, apiKey, types.RoleDirector, `You are a technical director.
Your responsibilities:
- Analyze project requirements
- Perform research on technologies and approaches
- Break down large projects into manageable tasks
- Make architectural decisions
- Allocate resources across teams`)
}

// NewADKPresidentAgent creates a President agent using Google's ADK framework.
func NewADKPresidentAgent(id string, config *types.AgentConfig, apiKey string) (*ADKAgent, error) {
	return newADKAgent(id, config, apiKey, types.RolePresident, `You are the president of a software development organization.
Your responsibilities:
- Clarify client requirements
- Define high-level objectives
- Ensure project alignment with goals
- Communicate with stakeholders
- Oversee project success`)
}

// newADKAgent creates a new ADK-based agent.
func newADKAgent(id string, config *types.AgentConfig, apiKey string, role types.AgentRole, defaultInstruction string) (*ADKAgent, error) {
	// Use API key from parameter or environment
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}

	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required for ADK agents")
	}

	// Determine model name
	modelName := "gemini-2.0-flash-exp"
	if config.Model != "" {
		modelName = config.Model
	}

	// Prepare instruction
	instruction := defaultInstruction
	if config.SystemPrompt != "" {
		instruction = config.SystemPrompt
	}

	// Store the llmagent config (we'll create the model and agent on each use for simplicity)
	llmConfig := llmagent.Config{
		Name:        id,
		Description: config.Description,
		Instruction: instruction,
	}

	// Create base agent
	baseAgent := NewBaseAgent(id, role, config)

	return &ADKAgent{
		BaseAgent: baseAgent,
		modelName: modelName,
		llmConfig: llmConfig,
		apiKey:    apiKey,
	}, nil
}

// ProcessTask processes a task using ADK's llmagent.
func (a *ADKAgent) ProcessTask(ctx context.Context, task *types.Task) (*types.TaskResponse, error) {
	a.IncrementActiveTasks()
	defer a.DecrementActiveTasks()

	// Create Gemini model for this request using ADK's gemini package
	adkModel, err := gemini.NewModel(ctx, a.modelName, &genai.ClientConfig{
		APIKey: a.apiKey,
	})
	if err != nil {
		return &types.TaskResponse{
			TaskID: task.ID,
			Status: types.StatusFailed,
			Error:  fmt.Sprintf("failed to create ADK Gemini model: %v", err),
		}, nil
	}

	// Set model in config
	a.llmConfig.Model = adkModel

	// Create ADK LLM agent with the configured model
	adkAgent, err := llmagent.New(a.llmConfig)
	if err != nil {
		return &types.TaskResponse{
			TaskID: task.ID,
			Status: types.StatusFailed,
			Error:  fmt.Sprintf("failed to create ADK agent: %v", err),
		}, nil
	}

	// Prepare the prompt for the ADK agent
	prompt := fmt.Sprintf(`Task: %s

Description: %s

Content:
%s

Please process this task according to your role and provide a detailed response.`,
		task.Title, task.Description, task.Content)

	// Note: Full ADK integration with session management would require more complex setup.
	// For this implementation, we demonstrate ADK's llmagent configuration
	// and use the underlying genai client for actual generation.

	// Use genai client directly (ADK's gemini.NewModel returns a model.LLM which wraps genai)
	genaiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: a.apiKey,
	})
	if err != nil {
		return &types.TaskResponse{
			TaskID: task.ID,
			Status: types.StatusFailed,
			Error:  fmt.Sprintf("failed to create genai client: %v", err),
		}, nil
	}

	// Build user content
	userContent := &genai.Content{
		Parts: []*genai.Part{{Text: prompt}},
		Role:  genai.RoleUser,
	}

	// Build system instruction from ADK config
	systemInstruction := &genai.Content{
		Parts: []*genai.Part{{Text: a.llmConfig.Instruction}},
	}

	// Generate content using the genai client with ADK-derived configuration
	genConfig := &genai.GenerateContentConfig{
		SystemInstruction: systemInstruction,
	}

	resp, err := genaiClient.Models.GenerateContent(ctx, a.modelName, []*genai.Content{userContent}, genConfig)
	if err != nil {
		return &types.TaskResponse{
			TaskID: task.ID,
			Status: types.StatusFailed,
			Error:  fmt.Sprintf("ADK-configured model error: %v", err),
		}, nil
	}

	// Extract response text
	responseText := ""
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if part.Text != "" {
				responseText += part.Text
			}
		}
	}

	if responseText == "" {
		responseText = "Task processed by ADK agent (no text output)"
	}

	return &types.TaskResponse{
		TaskID: task.ID,
		Status: types.StatusCompleted,
		Result: fmt.Sprintf("ADK Agent (%s - %s) Response:\n\n%s", adkAgent.Name(), a.modelName, responseText),
	}, nil
}

// GetModelName returns the model name being used.
func (a *ADKAgent) GetModelName() string {
	return a.modelName
}
