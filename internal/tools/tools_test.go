package tools

import (
	"context"
	"testing"
)

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	tool := &WebSearchTool{}

	err := registry.Register(tool)
	if err != nil {
		t.Fatalf("Failed to register tool: %v", err)
	}

	// Try to register same tool again
	err = registry.Register(tool)
	if err == nil {
		t.Error("Expected error when registering duplicate tool")
	}
}

func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()
	tool := &CodeAnalyzerTool{}
	registry.Register(tool)

	retrieved, err := registry.Get("code_analyzer")
	if err != nil {
		t.Fatalf("Failed to get tool: %v", err)
	}

	if retrieved.Name() != "code_analyzer" {
		t.Errorf("Expected 'code_analyzer', got '%s'", retrieved.Name())
	}
}

func TestRegistry_Execute(t *testing.T) {
	registry := NewRegistry()
	registry.Register(&WebSearchTool{})

	ctx := context.Background()
	result, err := registry.Execute(ctx, "web_search", map[string]interface{}{
		"query": "test query",
	})

	if err != nil {
		t.Fatalf("Failed to execute tool: %v", err)
	}

	if result == nil {
		t.Error("Expected result from tool execution")
	}
}

func TestRegistry_List(t *testing.T) {
	registry := NewDefaultRegistry()

	tools := registry.List()
	if len(tools) < 5 {
		t.Errorf("Expected at least 5 tools, got %d", len(tools))
	}
}

func TestWebSearchTool(t *testing.T) {
	tool := &WebSearchTool{}
	ctx := context.Background()

	result, err := tool.Execute(ctx, map[string]interface{}{
		"query": "golang",
	})

	if err != nil {
		t.Fatalf("Tool execution failed: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	}
}

func TestCodeAnalyzerTool(t *testing.T) {
	tool := &CodeAnalyzerTool{}
	ctx := context.Background()

	result, err := tool.Execute(ctx, map[string]interface{}{
		"code": "package main\nfunc main() {}",
	})

	if err != nil {
		t.Fatalf("Tool execution failed: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	}
}

func TestToolExecuteWithMissingParams(t *testing.T) {
	tool := &WebSearchTool{}
	ctx := context.Background()

	_, err := tool.Execute(ctx, map[string]interface{}{})

	if err == nil {
		t.Error("Expected error when executing without required parameters")
	}
}
