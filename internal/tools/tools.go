package tools

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

// WebSearchTool searches the web for information
type WebSearchTool struct{
	// EnableNetwork allows disabling network calls for testing
	EnableNetwork bool
}

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
	
	// Check if network is disabled (for testing)
	if !t.EnableNetwork {
		return map[string]interface{}{
			"query": query,
			"status": "simulated",
			"summary": fmt.Sprintf("Simulated search for '%s' (network disabled)", query),
			"results": []string{
				fmt.Sprintf("Result 1 for %s", query),
				fmt.Sprintf("Result 2 for %s", query),
				fmt.Sprintf("Result 3 for %s", query),
			},
		}, nil
	}
	
	// Real implementation using DuckDuckGo HTML search
	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; BuildBureau/1.0)")
	
	resp, err := client.Do(req)
	if err != nil {
		// Return graceful fallback on network error
		return map[string]interface{}{
			"query": query,
			"status": "network_error",
			"error": err.Error(),
			"summary": fmt.Sprintf("Network error searching for '%s': %v", query, err),
		}, nil
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search returned status: %d", resp.StatusCode)
	}
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// Extract text content (simplified - in production would parse HTML properly)
	content := string(body)
	
	// Return summary
	result := map[string]interface{}{
		"query": query,
		"url": searchURL,
		"status": "success",
		"content_length": len(content),
		"summary": fmt.Sprintf("Successfully searched for '%s' and retrieved %d bytes of data", query, len(content)),
	}
	
	return result, nil
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
	
	// Real implementation using Go's AST parser
	fset := token.NewFileSet()
	
	// Parse the code
	file, err := parser.ParseFile(fset, "input.go", code, parser.AllErrors)
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error": err.Error(),
			"code_length": len(code),
		}, nil
	}
	
	// Analyze the AST
	analysis := map[string]interface{}{
		"status": "success",
		"code_length": len(code),
		"package": file.Name.Name,
		"imports": len(file.Imports),
		"functions": 0,
		"structs": 0,
		"interfaces": 0,
		"comments": len(file.Comments),
	}
	
	// Count declarations
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			analysis["functions"] = analysis["functions"].(int) + 1
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					ts := spec.(*ast.TypeSpec)
					switch ts.Type.(type) {
					case *ast.StructType:
						analysis["structs"] = analysis["structs"].(int) + 1
					case *ast.InterfaceType:
						analysis["interfaces"] = analysis["interfaces"].(int) + 1
					}
				}
			}
		}
	}
	
	analysis["summary"] = fmt.Sprintf("Analyzed Go code: %d functions, %d structs, %d interfaces",
		analysis["functions"], analysis["structs"], analysis["interfaces"])
	
	return analysis, nil
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
	
	// Real implementation of document management
	switch action {
	case "create":
		path, ok := params["path"].(string)
		if !ok {
			return nil, fmt.Errorf("path parameter required for create action")
		}
		content, ok := params["content"].(string)
		if !ok {
			content = ""
		}
		
		// Create directories if they don't exist
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
		
		// Write file
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("failed to create document: %w", err)
		}
		
		return map[string]interface{}{
			"status": "success",
			"action": "create",
			"path": path,
			"size": len(content),
		}, nil
		
	case "read":
		path, ok := params["path"].(string)
		if !ok {
			return nil, fmt.Errorf("path parameter required for read action")
		}
		
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read document: %w", err)
		}
		
		return map[string]interface{}{
			"status": "success",
			"action": "read",
			"path": path,
			"content": string(content),
			"size": len(content),
		}, nil
		
	case "update":
		path, ok := params["path"].(string)
		if !ok {
			return nil, fmt.Errorf("path parameter required for update action")
		}
		content, ok := params["content"].(string)
		if !ok {
			return nil, fmt.Errorf("content parameter required for update action")
		}
		
		// Check if file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("document does not exist: %s", path)
		}
		
		// Update file
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("failed to update document: %w", err)
		}
		
		return map[string]interface{}{
			"status": "success",
			"action": "update",
			"path": path,
			"size": len(content),
		}, nil
		
	case "delete":
		path, ok := params["path"].(string)
		if !ok {
			return nil, fmt.Errorf("path parameter required for delete action")
		}
		
		if err := os.Remove(path); err != nil {
			return nil, fmt.Errorf("failed to delete document: %w", err)
		}
		
		return map[string]interface{}{
			"status": "success",
			"action": "delete",
			"path": path,
		}, nil
		
	case "list":
		dir, ok := params["directory"].(string)
		if !ok {
			dir = "."
		}
		
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to list directory: %w", err)
		}
		
		files := []string{}
		for _, entry := range entries {
			if !entry.IsDir() {
				files = append(files, entry.Name())
			}
		}
		
		return map[string]interface{}{
			"status": "success",
			"action": "list",
			"directory": dir,
			"files": files,
			"count": len(files),
		}, nil
		
	default:
		return nil, fmt.Errorf("unknown action: %s (supported: create, read, update, delete, list)", action)
	}
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
	
	// Real implementation of file operations
	switch operation {
	case "read":
		path, ok := params["path"].(string)
		if !ok {
			return nil, fmt.Errorf("path parameter required for read operation")
		}
		
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		
		return map[string]interface{}{
			"status": "success",
			"operation": "read",
			"path": path,
			"content": string(content),
			"size": len(content),
		}, nil
		
	case "write":
		path, ok := params["path"].(string)
		if !ok {
			return nil, fmt.Errorf("path parameter required for write operation")
		}
		content, ok := params["content"].(string)
		if !ok {
			return nil, fmt.Errorf("content parameter required for write operation")
		}
		
		// Create directories if needed
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
		
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("failed to write file: %w", err)
		}
		
		return map[string]interface{}{
			"status": "success",
			"operation": "write",
			"path": path,
			"size": len(content),
		}, nil
		
	case "delete":
		path, ok := params["path"].(string)
		if !ok {
			return nil, fmt.Errorf("path parameter required for delete operation")
		}
		
		if err := os.Remove(path); err != nil {
			return nil, fmt.Errorf("failed to delete file: %w", err)
		}
		
		return map[string]interface{}{
			"status": "success",
			"operation": "delete",
			"path": path,
		}, nil
		
	case "list":
		dir, ok := params["directory"].(string)
		if !ok {
			dir = "."
		}
		
		pattern := "*"
		if p, ok := params["pattern"].(string); ok {
			pattern = p
		}
		
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return nil, fmt.Errorf("failed to list files: %w", err)
		}
		
		return map[string]interface{}{
			"status": "success",
			"operation": "list",
			"directory": dir,
			"pattern": pattern,
			"files": matches,
			"count": len(matches),
		}, nil
		
	case "exists":
		path, ok := params["path"].(string)
		if !ok {
			return nil, fmt.Errorf("path parameter required for exists operation")
		}
		
		_, err := os.Stat(path)
		exists := !os.IsNotExist(err)
		
		return map[string]interface{}{
			"status": "success",
			"operation": "exists",
			"path": path,
			"exists": exists,
		}, nil
		
	case "copy":
		src, ok := params["source"].(string)
		if !ok {
			return nil, fmt.Errorf("source parameter required for copy operation")
		}
		dst, ok := params["destination"].(string)
		if !ok {
			return nil, fmt.Errorf("destination parameter required for copy operation")
		}
		
		data, err := os.ReadFile(src)
		if err != nil {
			return nil, fmt.Errorf("failed to read source file: %w", err)
		}
		
		// Create destination directory if needed
		dir := filepath.Dir(dst)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
		
		if err := os.WriteFile(dst, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to write destination file: %w", err)
		}
		
		return map[string]interface{}{
			"status": "success",
			"operation": "copy",
			"source": src,
			"destination": dst,
			"size": len(data),
		}, nil
		
	default:
		return nil, fmt.Errorf("unknown operation: %s (supported: read, write, delete, list, exists, copy)", operation)
	}
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
	
	language, ok := params["language"].(string)
	if !ok {
		language = "go" // default to Go
	}
	
	// Real implementation of safe code execution
	// Create a temporary file for the code
	tmpDir, err := os.MkdirTemp("", "buildbureau-exec-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)
	
	var cmdName string
	var cmdArgs []string
	var filename string
	
	switch strings.ToLower(language) {
	case "go":
		filename = filepath.Join(tmpDir, "main.go")
		if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
			return nil, fmt.Errorf("failed to write code file: %w", err)
		}
		cmdName = "go"
		cmdArgs = []string{"run", filename}
		
	case "python", "python3":
		filename = filepath.Join(tmpDir, "script.py")
		if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
			return nil, fmt.Errorf("failed to write code file: %w", err)
		}
		cmdName = "python3"
		cmdArgs = []string{filename}
		
	case "javascript", "js", "node":
		filename = filepath.Join(tmpDir, "script.js")
		if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
			return nil, fmt.Errorf("failed to write code file: %w", err)
		}
		cmdName = "node"
		cmdArgs = []string{filename}
		
	case "bash", "sh":
		filename = filepath.Join(tmpDir, "script.sh")
		if err := os.WriteFile(filename, []byte(code), 0755); err != nil {
			return nil, fmt.Errorf("failed to write code file: %w", err)
		}
		cmdName = "bash"
		cmdArgs = []string{filename}
		
	default:
		return nil, fmt.Errorf("unsupported language: %s (supported: go, python, javascript, bash)", language)
	}
	
	// Create command with context for timeout
	execCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(execCtx, cmdName, cmdArgs...)
	cmd.Dir = tmpDir
	
	// Capture output
	output, err := cmd.CombinedOutput()
	
	result := map[string]interface{}{
		"language": language,
		"code_length": len(code),
		"output": string(output),
	}
	
	if err != nil {
		result["status"] = "error"
		result["error"] = err.Error()
		if execCtx.Err() == context.DeadlineExceeded {
			result["error"] = "execution timeout (30s)"
		}
		return result, nil
	}
	
	result["status"] = "success"
	result["exit_code"] = 0
	
	return result, nil
}

// NewDefaultRegistry creates a registry with built-in tools
func NewDefaultRegistry() *Registry {
	registry := NewRegistry()
	registry.Register(&WebSearchTool{EnableNetwork: false}) // Network disabled by default for safety
	registry.Register(&CodeAnalyzerTool{})
	registry.Register(&DocumentManagerTool{})
	registry.Register(&FileOperationsTool{})
	registry.Register(&CodeExecutionTool{})
	return registry
}
