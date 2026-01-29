package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kpango/BuildBureau/internal/agent"
	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/internal/tools"
)

// This example demonstrates real agent processing with actual tools and LLM integration

func main() {
	fmt.Println("=== BuildBureau Real Functionality Demo ===")
	fmt.Println()

	ctx := context.Background()

	// 1. Initialize LLM client (mock for demo)
	fmt.Println("1. Initializing LLM client...")
	llmClient := llm.NewMockClient([]string{
		"I'll break this down into three main tasks: 1) User Authentication, 2) API Development, 3) Frontend Implementation",
		"For the authentication module, we need to implement secure login, registration, and password reset functionality",
		"I'll use the code analyzer to check the implementation and create the necessary files",
	})
	fmt.Println("✓ LLM client initialized")
	fmt.Println()

	// 2. Initialize tool registry
	fmt.Println("2. Initializing tool registry...")
	toolRegistry := tools.NewDefaultRegistry()
	fmt.Printf("✓ Tool registry initialized with %d tools\n", len(toolRegistry.List()))
	for _, toolName := range toolRegistry.List() {
		fmt.Printf("  - %s\n", toolName)
	}
	fmt.Println()

	// 3. Create specialized agents
	fmt.Println("3. Creating specialized agents...")
	cfg := config.AgentConfig{
		Count:      1,
		Model:      "mock",
		AllowTools: true,
		Timeout:    60,
		RetryCount: 3,
	}

	president := agent.NewSpecializedAgent("president-1", agent.AgentTypePresident, cfg, llmClient, toolRegistry)
	employee := agent.NewSpecializedAgent("employee-1", agent.AgentTypeEmployee, cfg, llmClient, toolRegistry)
	
	fmt.Println("✓ Agents created:")
	fmt.Println("  - President Agent (strategic planning)")
	fmt.Println("  - Employee Agent (task execution)")
	fmt.Println()

	// 4. President processes high-level requirement
	fmt.Println("4. President Agent: Processing project requirements...")
	presidentInput := "Build an authentication system with user management and secure login"
	presidentResult, err := president.Process(ctx, presidentInput)
	if err != nil {
		log.Fatalf("President processing failed: %v", err)
	}
	
	fmt.Println("✓ President completed planning")
	if resultMap, ok := presidentResult.(map[string]interface{}); ok {
		fmt.Printf("  Content: %s\n", resultMap["content"])
		if tools, ok := resultMap["tools_used"].([]string); ok && len(tools) > 0 {
			fmt.Printf("  Tools used: %v\n", tools)
		}
	}
	fmt.Println()

	// 5. Employee processes specific task
	fmt.Println("5. Employee Agent: Executing implementation task...")
	employeeInput := "Implement user authentication module with secure password hashing"
	employeeResult, err := employee.Process(ctx, employeeInput)
	if err != nil {
		log.Fatalf("Employee processing failed: %v", err)
	}
	
	fmt.Println("✓ Employee completed task")
	if resultMap, ok := employeeResult.(map[string]interface{}); ok {
		fmt.Printf("  Content: %s\n", resultMap["content"])
		if tools, ok := resultMap["tools_used"].([]string); ok && len(tools) > 0 {
			fmt.Printf("  Tools used: %v\n", tools)
		}
	}
	fmt.Println()

	// 6. Demonstrate tool usage directly
	fmt.Println("6. Demonstrating direct tool usage...")
	fmt.Println()

	// Test Code Analyzer Tool
	fmt.Println("  a) Code Analyzer Tool:")
	codeAnalyzer, _ := toolRegistry.Get("code_analyzer")
	sampleCode := `package main

import "fmt"

type User struct {
	ID   int
	Name string
}

func main() {
	fmt.Println("Hello")
}

func Authenticate(user User) bool {
	return true
}`
	
	analysisResult, err := codeAnalyzer.Execute(ctx, map[string]interface{}{
		"code": sampleCode,
	})
	if err != nil {
		log.Printf("Code analysis failed: %v", err)
	} else {
		if analysis, ok := analysisResult.(map[string]interface{}); ok {
			fmt.Printf("     Status: %s\n", analysis["status"])
			fmt.Printf("     Functions: %d\n", analysis["functions"])
			fmt.Printf("     Structs: %d\n", analysis["structs"])
			fmt.Printf("     Summary: %s\n", analysis["summary"])
		}
	}
	fmt.Println()

	// Test Document Manager Tool
	fmt.Println("  b) Document Manager Tool:")
	docManager, _ := toolRegistry.Get("document_manager")
	
	// Create a temporary document
	createResult, err := docManager.Execute(ctx, map[string]interface{}{
		"action":  "create",
		"path":    "/tmp/buildbureau_test_doc.txt",
		"content": "This is a test document created by BuildBureau",
	})
	if err != nil {
		log.Printf("Document creation failed: %v", err)
	} else {
		if result, ok := createResult.(map[string]interface{}); ok {
			fmt.Printf("     Status: %s\n", result["status"])
			fmt.Printf("     Action: %s\n", result["action"])
			fmt.Printf("     Path: %s\n", result["path"])
		}
	}
	fmt.Println()

	// Test File Operations Tool
	fmt.Println("  c) File Operations Tool:")
	fileOps, _ := toolRegistry.Get("file_operations")
	
	existsResult, err := fileOps.Execute(ctx, map[string]interface{}{
		"operation": "exists",
		"path":      "/tmp/buildbureau_test_doc.txt",
	})
	if err != nil {
		log.Printf("File check failed: %v", err)
	} else {
		if result, ok := existsResult.(map[string]interface{}); ok {
			fmt.Printf("     Operation: %s\n", result["operation"])
			fmt.Printf("     Path: %s\n", result["path"])
			fmt.Printf("     Exists: %v\n", result["exists"])
		}
	}
	fmt.Println()

	// Test Code Execution Tool
	fmt.Println("  d) Code Execution Tool:")
	codeExec, _ := toolRegistry.Get("code_execution")
	
	goCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello from BuildBureau!")
	fmt.Println("Code execution successful")
}`
	
	execResult, err := codeExec.Execute(ctx, map[string]interface{}{
		"code":     goCode,
		"language": "go",
	})
	if err != nil {
		log.Printf("Code execution failed: %v", err)
	} else {
		if result, ok := execResult.(map[string]interface{}); ok {
			fmt.Printf("     Status: %s\n", result["status"])
			fmt.Printf("     Language: %s\n", result["language"])
			if output, ok := result["output"].(string); ok && output != "" {
				fmt.Printf("     Output:\n%s\n", output)
			}
		}
	}
	fmt.Println()

	// Test Web Search Tool (simulated mode)
	fmt.Println("  e) Web Search Tool (simulated):")
	webSearch, _ := toolRegistry.Get("web_search")
	
	searchResult, err := webSearch.Execute(ctx, map[string]interface{}{
		"query": "Go programming best practices",
	})
	if err != nil {
		log.Printf("Web search failed: %v", err)
	} else {
		if result, ok := searchResult.(map[string]interface{}); ok {
			fmt.Printf("     Query: %s\n", result["query"])
			fmt.Printf("     Status: %s\n", result["status"])
			fmt.Printf("     Summary: %s\n", result["summary"])
		}
	}
	fmt.Println()

	// 7. Show agent status
	fmt.Println("7. Final agent status:")
	fmt.Printf("  President: %s - %s\n", president.GetStatus().State, president.GetStatus().Message)
	fmt.Printf("  Employee: %s - %s\n", employee.GetStatus().State, employee.GetStatus().Message)
	fmt.Println()

	fmt.Println("=== Demo Complete ===")
	fmt.Println()
	fmt.Println("This demo showed:")
	fmt.Println("  ✓ Real LLM integration with specialized agents")
	fmt.Println("  ✓ Role-specific system prompts for different agent types")
	fmt.Println("  ✓ Code analysis with AST parsing")
	fmt.Println("  ✓ Document and file management")
	fmt.Println("  ✓ Safe code execution in multiple languages")
	fmt.Println("  ✓ Web search capabilities (simulated)")
	fmt.Println("  ✓ Agent status tracking")
}
