package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kpango/BuildBureau/internal/agent"
	"github.com/kpango/BuildBureau/pkg/types"
)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║       BuildBureau - ADK Agent Integration Example         ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Check for API key
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" || apiKey == "demo-key" {
		fmt.Println("⚠️  GEMINI_API_KEY not set or using demo value")
		fmt.Println()
		fmt.Println("To run this example with real ADK agents:")
		fmt.Println("  1. Get a free API key from: https://aistudio.google.com/app/apikey")
		fmt.Println("  2. Set it: export GEMINI_API_KEY='your-actual-key'")
		fmt.Println("  3. Run this example again")
		fmt.Println()
		fmt.Println("This example requires a real API key to demonstrate ADK functionality.")
		os.Exit(0)
	}

	fmt.Println("✓ GEMINI_API_KEY found")
	fmt.Println()

	// Test Engineer Agent
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("Testing ADK Engineer Agent")
	fmt.Println("═══════════════════════════════════════════════════════════")
	testEngineerAgent(apiKey)
	fmt.Println()

	// Test Manager Agent
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("Testing ADK Manager Agent")
	fmt.Println("═══════════════════════════════════════════════════════════")
	testManagerAgent(apiKey)
	fmt.Println()

	// Test Director Agent
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("Testing ADK Director Agent")
	fmt.Println("═══════════════════════════════════════════════════════════")
	testDirectorAgent(apiKey)
	fmt.Println()

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║              All ADK Agent Tests Completed! ✓              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}

// testADKAgent is a helper function to test ADK agents with different roles.
func testADKAgent(agentType, agentID string, config *types.AgentConfig, task *types.Task, apiKey string, createAgent func(string, *types.AgentConfig, string) (any, error)) {
	fmt.Printf("Creating ADK %s Agent...\n", agentType)

	agentInterface, err := createAgent(agentID, config, apiKey)
	if err != nil {
		log.Fatalf("Failed to create ADK %s: %v", agentType, err)
	}

	// Type assertion to get the ADKAgent
	adkAgent, ok := agentInterface.(*agent.ADKAgent)
	if !ok {
		log.Fatalf("Failed to cast to ADKAgent")
	}

	fmt.Printf("✓ ADK %s created (Model: %s)\n", agentType, adkAgent.GetModelName())

	ctx := context.Background()
	if err := adkAgent.Start(ctx); err != nil {
		log.Fatalf("Failed to start %s: %v", agentType, err)
	}
	defer adkAgent.Stop(ctx)

	fmt.Printf("✓ %s started\n", agentType)

	fmt.Printf("Sending task to ADK %s...\n", agentType)
	fmt.Printf("Task: %s\n", task.Title)
	fmt.Println()

	response, err := adkAgent.ProcessTask(ctx, task)
	if err != nil {
		log.Fatalf("Failed to process task: %v", err)
	}

	fmt.Println("Response received:")
	fmt.Println("─────────────────────────────────────────────────────────")
	fmt.Println(response.Result)
	fmt.Println("─────────────────────────────────────────────────────────")
	fmt.Printf("Status: %s\n", response.Status)
}

func testEngineerAgent(apiKey string) {
	config := &types.AgentConfig{
		Name:        "ADK Engineer",
		Role:        "Engineer",
		Description: "An engineer that uses Google ADK for task processing",
	}

	task := &types.Task{
		ID:          "adk-task-1",
		Title:       "Implement Fibonacci Function",
		Description: "Create a function to calculate Fibonacci numbers",
		FromAgent:   "manager-1",
		ToAgent:     "adk-engineer-1",
		Content:     "Write a Python function that calculates the nth Fibonacci number using dynamic programming. Include error handling for negative inputs.",
		Priority:    1,
	}

	testADKAgent("Engineer", "adk-engineer-1", config, task, apiKey, func(id string, cfg *types.AgentConfig, key string) (any, error) {
		return agent.NewADKEngineerAgent(id, cfg, key)
	})
}

func testManagerAgent(apiKey string) {
	config := &types.AgentConfig{
		Name:        "ADK Manager",
		Role:        "Manager",
		Description: "A manager that uses Google ADK for creating specifications",
	}

	task := &types.Task{
		ID:          "adk-task-2",
		Title:       "Design User Authentication System",
		Description: "Create technical specification for a user authentication system",
		FromAgent:   "director-1",
		ToAgent:     "adk-manager-1",
		Content:     "Design a secure user authentication system with JWT tokens, password hashing, and refresh tokens. Include API endpoints, database schema, and security considerations.",
		Priority:    1,
	}

	testADKAgent("Manager", "adk-manager-1", config, task, apiKey, func(id string, cfg *types.AgentConfig, key string) (any, error) {
		return agent.NewADKManagerAgent(id, cfg, key)
	})
}

func testDirectorAgent(apiKey string) {
	fmt.Println("Creating ADK Director Agent...")

	config := &types.AgentConfig{
		Name:        "ADK Director",
		Role:        "Director",
		Description: "A director that uses Google ADK for project analysis",
	}

	director, err := agent.NewADKDirectorAgent("adk-director-1", config, apiKey)
	if err != nil {
		log.Fatalf("Failed to create ADK director: %v", err)
	}

	fmt.Printf("✓ ADK Director created (Model: %s)\n", director.GetModelName())

	ctx := context.Background()
	if err := director.Start(ctx); err != nil {
		log.Fatalf("Failed to start director: %v", err)
	}
	defer director.Stop(ctx)

	fmt.Println("✓ Director started")

	// Create a test task
	task := &types.Task{
		ID:          "adk-task-3",
		Title:       "Analyze E-commerce Platform Requirements",
		Description: "Break down e-commerce platform project into tasks",
		FromAgent:   "president-1",
		ToAgent:     "adk-director-1",
		Content:     "Analyze the requirements for building a modern e-commerce platform with product catalog, shopping cart, payment processing, and order management. Break this down into major components and suggest technology choices.",
		Priority:    1,
	}

	fmt.Println("Sending task to ADK Director...")
	fmt.Printf("Task: %s\n", task.Title)
	fmt.Println()

	response, err := director.ProcessTask(ctx, task)
	if err != nil {
		log.Fatalf("Failed to process task: %v", err)
	}

	fmt.Println("Response received:")
	fmt.Println("─────────────────────────────────────────────────────────")
	fmt.Println(response.Result)
	fmt.Println("─────────────────────────────────────────────────────────")
	fmt.Printf("Status: %s\n", response.Status)
}
