package agent

import (
	"context"
	"testing"

	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/internal/tools"
)

func TestSpecializedAgent_Process(t *testing.T) {
	ctx := context.Background()
	llmClient := llm.NewMockClient([]string{"Task breakdown complete"})
	toolRegistry := tools.NewDefaultRegistry()

	cfg := config.AgentConfig{
		Count:      1,
		Model:      "mock",
		AllowTools: true,
		Timeout:    60,
		RetryCount: 3,
	}

	agent := NewSpecializedAgent("test-1", AgentTypeEmployee, cfg, llmClient, toolRegistry)

	result, err := agent.Process(ctx, "Complete implementation task")
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Check status after processing
	status := agent.GetStatus()
	if status.State != "idle" {
		t.Errorf("Expected status 'idle', got '%s'", status.State)
	}
}

func TestSpecializedAgent_GetSystemPrompt(t *testing.T) {
	tests := []struct {
		name      string
		agentType AgentType
		wantContains string
	}{
		{
			name:         "President prompt",
			agentType:    AgentTypePresident,
			wantContains: "President",
		},
		{
			name:         "Department Manager prompt",
			agentType:    AgentTypeDepartmentManager,
			wantContains: "Department Manager",
		},
		{
			name:         "Section Manager prompt",
			agentType:    AgentTypeSectionManager,
			wantContains: "Section Manager",
		},
		{
			name:         "Employee prompt",
			agentType:    AgentTypeEmployee,
			wantContains: "Employee",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.AgentConfig{Model: "mock"}
			agent := NewSpecializedAgent("test", tt.agentType, cfg, nil, nil)
			prompt := agent.getSystemPrompt()
			
			if prompt == "" {
				t.Error("Expected non-empty system prompt")
			}
		})
	}
}

func TestSpecializedAgent_CustomInstruction(t *testing.T) {
	cfg := config.AgentConfig{
		Model:       "mock",
		Instruction: "Custom instruction for testing",
	}

	agent := NewSpecializedAgent("test", AgentTypeEmployee, cfg, nil, nil)
	prompt := agent.getSystemPrompt()

	if prompt != cfg.Instruction {
		t.Errorf("Expected custom instruction, got: %s", prompt)
	}
}

func TestStreamingAgent_ProcessStream(t *testing.T) {
	ctx := context.Background()
	llmClient := llm.NewMockClient([]string{"Streaming response"})
	toolRegistry := tools.NewDefaultRegistry()

	cfg := config.AgentConfig{
		Count:      1,
		Model:      "mock",
		AllowTools: false,
		Timeout:    60,
	}

	agent := NewStreamingAgent("stream-1", AgentTypeEmployee, cfg, llmClient, toolRegistry)

	contentCh, errCh, err := agent.ProcessStream(ctx, "Stream this task")
	if err != nil {
		t.Fatalf("ProcessStream failed: %v", err)
	}

	// Read all content
	var received []string
	for content := range contentCh {
		received = append(received, content)
	}

	// Check for errors
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Received error from stream: %v", err)
		}
	default:
	}

	if len(received) == 0 {
		t.Error("Expected streaming content, got none")
	}
}

func TestSpecializedAgent_WithoutTools(t *testing.T) {
	ctx := context.Background()
	llmClient := llm.NewMockClient([]string{"Response without tools"})

	cfg := config.AgentConfig{
		Model:      "mock",
		AllowTools: false, // Tools disabled
	}

	agent := NewSpecializedAgent("test-no-tools", AgentTypeEmployee, cfg, llmClient, nil)

	result, err := agent.Process(ctx, "Process without tools")
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map result")
	}

	// Should not have used any tools
	toolsUsed := resultMap["tools_used"].([]string)
	if len(toolsUsed) != 0 {
		t.Errorf("Expected no tools used, got: %v", toolsUsed)
	}
}
