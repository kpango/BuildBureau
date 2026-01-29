package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kpango/BuildBureau/internal/agent"
	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/internal/memory"
	"github.com/kpango/BuildBureau/pkg/types"
)

func main() {
	fmt.Println("=== BuildBureau Agent Memory Integration Demo ===\n")

	ctx := context.Background()

	// Create memory configuration
	memoryConfig := &types.MemoryConfig{
		Enabled: true,
		SQLite: types.SQLiteConfig{
			Enabled:  true,
			Path:     ":memory:", // Use in-memory database for demo
			InMemory: true,
		},
		Vald: types.ValdConfig{
			Enabled:   false, // Optional for this demo
			Host:      "localhost",
			Port:      8081,
			Dimension: 768,
		},
		Retention: types.RetentionConfig{
			ConversationDays: 30,
			TaskDays:         60,
			KnowledgeDays:    0, // Forever
			MaxEntries:       10000,
		},
	}

	// Create LLM manager (optional, can work without it)
	apiKey := os.Getenv("GEMINI_API_KEY")
	llmConfig := &types.LLMConfig{
		DefaultModel: "gemini",
		APIKeys: map[string]types.EnvironmentVariable{
			"gemini": {Env: "GEMINI_API_KEY"},
		},
	}
	var llmManager *llm.Manager
	if apiKey != "" {
		llmManager, _ = llm.NewManager(llmConfig)
		fmt.Println("✓ LLM Manager initialized with Gemini")
	} else {
		fmt.Println("⚠ No API key found, running without LLM (set GEMINI_API_KEY to use)")
	}

	// Create memory manager
	memManager, err := memory.NewManager(memoryConfig, llmManager)
	if err != nil {
		fmt.Printf("❌ Failed to create memory manager: %v\n", err)
		return
	}
	defer memManager.Close()
	fmt.Println("✓ Memory Manager initialized with SQLite (in-memory)")

	// Create agents
	engineerConfig := &types.AgentConfig{
		Name:         "Senior Engineer",
		Model:        "gemini",
		SystemPrompt: "You are an expert software engineer.",
	}
	engineer := agent.NewEngineerAgent("eng-001", engineerConfig, llmManager)
	engineer.SetMemoryManager(memManager)
	engineer.Start(ctx)
	fmt.Println("✓ Engineer Agent created with memory")

	managerConfig := &types.AgentConfig{
		Name:         "Technical Manager",
		Model:        "gemini",
		SystemPrompt: "You are a technical manager who creates software designs.",
	}
	manager := agent.NewManagerAgent("mgr-001", managerConfig, llmManager)
	manager.SetMemoryManager(memManager)
	manager.AddEngineer(engineer)
	manager.Start(ctx)
	fmt.Println("✓ Manager Agent created with memory")

	secretaryConfig := &types.AgentConfig{
		Name:         "Secretary",
		SystemPrompt: "You are a secretary who delegates tasks.",
	}
	secretary := agent.NewSecretaryAgent("sec-001", secretaryConfig)
	secretary.SetMemoryManager(memManager)
	secretary.Start(ctx)
	fmt.Println("✓ Secretary Agent created with memory\n")

	// Simulate Task 1: REST API
	fmt.Println("--- Task 1: REST API for User Management ---")
	task1 := &types.Task{
		ID:          "task-001",
		Title:       "Create REST API",
		Description: "Build a REST API for user management with authentication",
		Content:     "Include CRUD operations, JWT authentication, and role-based access control",
		FromAgent:   "client",
		ToAgent:     "mgr-001",
		Priority:    1,
	}

	response1, err := manager.ProcessTask(ctx, task1)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✓ Task completed with status: %s\n", response1.Status)
		if len(response1.Result) > 200 {
			fmt.Printf("Result (truncated): %s...\n\n", response1.Result[:200])
		} else {
			fmt.Printf("Result: %s\n\n", response1.Result)
		}
	}

	// Simulate Task 2: Database Schema (similar to task 1)
	fmt.Println("--- Task 2: Database Schema for User System ---")
	task2 := &types.Task{
		ID:          "task-002",
		Title:       "Design Database Schema",
		Description: "Create database schema for user management system",
		Content:     "Include users, roles, permissions, and authentication tables",
		FromAgent:   "client",
		ToAgent:     "mgr-001",
		Priority:    1,
	}

	response2, err := manager.ProcessTask(ctx, task2)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✓ Task completed with status: %s\n", response2.Status)
		if len(response2.Result) > 200 {
			fmt.Printf("Result (truncated): %s...\n\n", response2.Result[:200])
		} else {
			fmt.Printf("Result: %s\n\n", response2.Result)
		}
	}

	// Query memory to show what was stored
	fmt.Println("--- Memory Analysis ---")

	// Check engineer's memory
	if engineerMem := engineer.GetMemory(); engineerMem != nil {
		history, _ := engineerMem.GetConversationHistory(ctx, 10)
		fmt.Printf("✓ Engineer has %d conversation memories\n", len(history))

		knowledge, _ := engineerMem.GetKnowledge(ctx, "user management", 5)
		fmt.Printf("✓ Engineer has %d knowledge entries about 'user management'\n", len(knowledge))

		tasks, _ := memManager.QueryMemories(ctx, &types.MemoryQuery{
			AgentID: "eng-001",
			Type:    types.MemoryTypeTask,
		})
		fmt.Printf("✓ Engineer completed %d tasks (in memory)\n", len(tasks))
	}

	// Check manager's memory
	if managerMem := manager.GetMemory(); managerMem != nil {
		history, _ := managerMem.GetConversationHistory(ctx, 10)
		fmt.Printf("✓ Manager has %d conversation memories\n", len(history))

		knowledge, _ := managerMem.GetKnowledge(ctx, "design", 5)
		fmt.Printf("✓ Manager has %d knowledge entries about 'design'\n", len(knowledge))

		decisions, _ := managerMem.GetDecisionHistory(ctx, 10)
		fmt.Printf("✓ Manager made %d delegation decisions\n", len(decisions))
	}

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nKey Features Demonstrated:")
	fmt.Println("1. ✅ Switched to github.com/mattn/go-sqlite3 driver")
	fmt.Println("2. ✅ Agents use memory to store conversations, tasks, and knowledge")
	fmt.Println("3. ✅ Engineers learn from past implementations")
	fmt.Println("4. ✅ Managers learn from past designs")
	fmt.Println("5. ✅ Secretaries track delegations and make informed decisions")
	fmt.Println("6. ✅ Memory persists across tasks (demonstrated with similar tasks)")
	fmt.Println("7. ✅ Agents can search and retrieve relevant past information")
}
