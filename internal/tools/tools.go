package tools

import (
	"context"
	"fmt"
)

// Tool represents an executable tool that agents can use
type Tool interface {
	// Name returns the tool name
	Name() string

	// Description returns what the tool does
	Description() string

	// Execute runs the tool with given parameters
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

// Registry manages available tools
type Registry struct {
	tools map[string]Tool
}

// NewRegistry creates a new tool registry
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry
func (r *Registry) Register(tool Tool) error {
	name := tool.Name()
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %s already registered", name)
	}
	r.tools[name] = tool
	return nil
}

// Get retrieves a tool by name
func (r *Registry) Get(name string) (Tool, error) {
	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}
	return tool, nil
}

// List returns all registered tool names
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// Execute runs a tool by name with given parameters
func (r *Registry) Execute(ctx context.Context, toolName string, params map[string]interface{}) (interface{}, error) {
	tool, err := r.Get(toolName)
	if err != nil {
		return nil, err
	}
	return tool.Execute(ctx, params)
}

// Built-in tools

// WebSearchTool simulates web search (placeholder)
type WebSearchTool struct{}

func (t *WebSearchTool) Name() string {
	return "web_search"
}

func (t *WebSearchTool) Description() string {
	return "Search the web for information"
}

func (t *WebSearchTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	query, ok := params["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query parameter required")
	}
	// Placeholder implementation
	return fmt.Sprintf("Search results for: %s", query), nil
}

// CodeAnalyzerTool simulates code analysis (placeholder)
type CodeAnalyzerTool struct{}

func (t *CodeAnalyzerTool) Name() string {
	return "code_analyzer"
}

func (t *CodeAnalyzerTool) Description() string {
	return "Analyze code structure and quality"
}

func (t *CodeAnalyzerTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	code, ok := params["code"].(string)
	if !ok {
		return nil, fmt.Errorf("code parameter required")
	}
	// Placeholder implementation
	return fmt.Sprintf("Analysis of code (%d bytes)", len(code)), nil
}

// DocumentManagerTool simulates document management (placeholder)
type DocumentManagerTool struct{}

func (t *DocumentManagerTool) Name() string {
	return "document_manager"
}

func (t *DocumentManagerTool) Description() string {
	return "Manage and organize documents"
}

func (t *DocumentManagerTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	action, ok := params["action"].(string)
	if !ok {
		return nil, fmt.Errorf("action parameter required")
	}
	// Placeholder implementation
	return fmt.Sprintf("Document action: %s", action), nil
}

// FileOperationsTool simulates file operations (placeholder)
type FileOperationsTool struct{}

func (t *FileOperationsTool) Name() string {
	return "file_operations"
}

func (t *FileOperationsTool) Description() string {
	return "Perform file system operations"
}

func (t *FileOperationsTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter required")
	}
	// Placeholder implementation
	return fmt.Sprintf("File operation: %s", operation), nil
}

// CodeExecutionTool simulates code execution (placeholder)
type CodeExecutionTool struct{}

func (t *CodeExecutionTool) Name() string {
	return "code_execution"
}

func (t *CodeExecutionTool) Description() string {
	return "Execute code safely in a sandbox"
}

func (t *CodeExecutionTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	code, ok := params["code"].(string)
	if !ok {
		return nil, fmt.Errorf("code parameter required")
	}
	// Placeholder implementation
	return fmt.Sprintf("Executed code (%d bytes)", len(code)), nil
}

// NewDefaultRegistry creates a registry with built-in tools
func NewDefaultRegistry() *Registry {
	registry := NewRegistry()
	registry.Register(&WebSearchTool{})
	registry.Register(&CodeAnalyzerTool{})
	registry.Register(&DocumentManagerTool{})
	registry.Register(&FileOperationsTool{})
	registry.Register(&CodeExecutionTool{})
	return registry
}
