package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/internal/tools"
)

// SpecializedAgent wraps BaseAgent with LLM and tool capabilities
type SpecializedAgent struct {
	*BaseAgent
	llmClient    llm.Client
	toolRegistry *tools.Registry
}

// NewSpecializedAgent creates a specialized agent with LLM and tool support
func NewSpecializedAgent(id string, agentType AgentType, cfg config.AgentConfig, llmClient llm.Client, toolRegistry *tools.Registry) *SpecializedAgent {
	return &SpecializedAgent{
		BaseAgent:    NewBaseAgent(id, agentType, cfg),
		llmClient:    llmClient,
		toolRegistry: toolRegistry,
	}
}

// Process implements intelligent processing using LLM and tools
func (a *SpecializedAgent) Process(ctx context.Context, input interface{}) (interface{}, error) {
	// Update status
	a.UpdateStatus("working", fmt.Sprintf("%v", input), "Processing task")

	// Convert input to string for processing
	inputStr := fmt.Sprintf("%v", input)

	// Build system prompt based on agent type
	systemPrompt := a.getSystemPrompt()

	// Prepare LLM request
	req := &llm.Request{
		Messages: []llm.Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: inputStr,
			},
		},
		Model:       a.config.Model,
		Temperature: 0.7,
		MaxTokens:   2000,
	}

	// Generate response using LLM
	resp, err := a.llmClient.Generate(ctx, req)
	if err != nil {
		a.UpdateStatus("error", "", fmt.Sprintf("LLM error: %v", err))
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	// Check if response suggests tool usage
	result := a.processToolUsage(ctx, resp.Content)

	// Update status
	a.UpdateStatus("idle", "", "Task completed successfully")

	return result, nil
}

// getSystemPrompt returns a role-specific system prompt
func (a *SpecializedAgent) getSystemPrompt() string {
	// Use config instruction if available
	if a.config.Instruction != "" {
		return a.config.Instruction
	}

	// Default prompts based on agent type
	switch a.agentType {
	case AgentTypePresident:
		return "You are the President of BuildBureau. Your role is to understand client requirements and break them down into high-level strategic tasks. Focus on overall project planning and resource allocation."

	case AgentTypePresidentSecretary:
		return "You are the President's Secretary. Your role is to organize and clarify the President's directives, maintain documentation, and ensure smooth communication with department managers."

	case AgentTypeDepartmentManager:
		return "You are a Department Manager at BuildBureau. Your role is to take strategic tasks from the President and divide them into manageable sections for section managers. Focus on resource planning and timeline management."

	case AgentTypeDepartmentSecretary:
		return "You are the Department Manager's Secretary. Your role is to coordinate between the department and sections, maintain detailed task documentation, and ensure all section managers have clear objectives."

	case AgentTypeSectionManager:
		return "You are a Section Manager at BuildBureau. Your role is to take section-level tasks and create detailed implementation plans for employees. Focus on technical specifications and work breakdown."

	case AgentTypeSectionSecretary:
		return "You are a Section Manager's Secretary. Your role is to maintain detailed implementation specifications, coordinate with employees, and track progress on all section tasks."

	case AgentTypeEmployee:
		return "You are an Employee at BuildBureau. Your role is to execute specific tasks according to detailed specifications. Focus on implementation, quality, and reporting results clearly."

	default:
		return "You are an AI agent assisting with task processing."
	}
}

// processToolUsage checks if the response suggests using tools and executes them
func (a *SpecializedAgent) processToolUsage(ctx context.Context, content string) interface{} {
	// Simple heuristic: check for tool-related keywords
	// In a real implementation, this would parse structured tool calls from LLM
	contentLower := strings.ToLower(content)

	result := map[string]interface{}{
		"content": content,
		"tools_used": []string{},
	}

	// Check if we should use tools (only if allowed in config)
	if !a.config.AllowTools || a.toolRegistry == nil {
		return result
	}

	// Example: Check for web search
	if strings.Contains(contentLower, "search") && strings.Contains(contentLower, "web") {
		if tool, err := a.toolRegistry.Get("web_search"); err == nil {
			// Extract query (simplified)
			query := "information query"
			if toolResult, err := tool.Execute(ctx, map[string]interface{}{"query": query}); err == nil {
				result["web_search"] = toolResult
				result["tools_used"] = append(result["tools_used"].([]string), "web_search")
			}
		}
	}

	// Example: Check for code analysis
	if strings.Contains(contentLower, "analyze") && strings.Contains(contentLower, "code") {
		if tool, err := a.toolRegistry.Get("code_analyzer"); err == nil {
			// Would extract code from context in real implementation
			sampleCode := "package main\n\nfunc main() {}\n"
			if toolResult, err := tool.Execute(ctx, map[string]interface{}{"code": sampleCode}); err == nil {
				result["code_analysis"] = toolResult
				result["tools_used"] = append(result["tools_used"].([]string), "code_analyzer")
			}
		}
	}

	return result
}

// StreamingAgent extends SpecializedAgent with streaming capabilities
type StreamingAgent struct {
	*SpecializedAgent
}

// NewStreamingAgent creates an agent that supports streaming responses
func NewStreamingAgent(id string, agentType AgentType, cfg config.AgentConfig, llmClient llm.Client, toolRegistry *tools.Registry) *StreamingAgent {
	return &StreamingAgent{
		SpecializedAgent: NewSpecializedAgent(id, agentType, cfg, llmClient, toolRegistry),
	}
}

// ProcessStream processes input and returns streaming results
func (a *StreamingAgent) ProcessStream(ctx context.Context, input interface{}) (<-chan string, <-chan error, error) {
	// Update status
	a.UpdateStatus("working", fmt.Sprintf("%v", input), "Processing task (streaming)")

	inputStr := fmt.Sprintf("%v", input)
	systemPrompt := a.getSystemPrompt()

	req := &llm.Request{
		Messages: []llm.Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: inputStr,
			},
		},
		Model:       a.config.Model,
		Temperature: 0.7,
		MaxTokens:   2000,
	}

	// Generate streaming response
	contentCh, errCh := a.llmClient.StreamGenerate(ctx, req)

	// Create output channel
	outputCh := make(chan string, 10)

	// Forward content with status updates
	go func() {
		defer close(outputCh)
		defer func() {
			a.UpdateStatus("idle", "", "Task completed successfully")
		}()

		for {
			select {
			case content, ok := <-contentCh:
				if !ok {
					return
				}
				outputCh <- content
			case <-ctx.Done():
				return
			}
		}
	}()

	return outputCh, errCh, nil
}
