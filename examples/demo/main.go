package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kpango/BuildBureau/internal/agent"
	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/internal/grpc"
	"github.com/kpango/BuildBureau/internal/knowledge"
	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/internal/tools"
)

// This example demonstrates the complete workflow of the BuildBureau system
// from project planning to task execution using all agent types.

func main() {
	fmt.Println("=== BuildBureau System Demo ===")
	fmt.Println()

	// 1. Initialize the system components
	fmt.Println("1. Initializing system components...")
	
	// Create agent pool
	agentPool := agent.NewAgentPool()
	
	// Create agents (simplified configuration)
	cfg := config.AgentConfig{
		Count:      1,
		Model:      "mock",
		AllowTools: true,
		Timeout:    60,
		RetryCount: 3,
	}
	
	// Register all agent types
	agentPool.Register(agent.NewBaseAgent("president-1", agent.AgentTypePresident, cfg))
	agentPool.Register(agent.NewBaseAgent("president-sec-1", agent.AgentTypePresidentSecretary, cfg))
	agentPool.Register(agent.NewBaseAgent("dept-mgr-1", agent.AgentTypeDepartmentManager, cfg))
	agentPool.Register(agent.NewBaseAgent("dept-sec-1", agent.AgentTypeDepartmentSecretary, cfg))
	agentPool.Register(agent.NewBaseAgent("section-mgr-1", agent.AgentTypeSectionManager, cfg))
	agentPool.Register(agent.NewBaseAgent("section-sec-1", agent.AgentTypeSectionSecretary, cfg))
	agentPool.Register(agent.NewBaseAgent("employee-1", agent.AgentTypeEmployee, cfg))
	
	// Initialize knowledge base
	kb := knowledge.NewInMemoryKB()
	
	// Initialize tool registry
	toolRegistry := tools.NewDefaultRegistry()
	
	// Initialize LLM client (using mock for demo)
	llmClient := llm.NewMockClient([]string{
		"Task 1: Design system architecture\nTask 2: Implement backend API\nTask 3: Build frontend interface",
	})
	
	fmt.Println("✓ System initialized with 7 agents")
	fmt.Println()
	
	// 2. Create services
	fmt.Println("2. Creating gRPC services...")
	presidentService := grpc.NewPresidentService(agentPool, kb, toolRegistry, llmClient)
	deptMgrService := grpc.NewDepartmentManagerService(agentPool, kb, toolRegistry, llmClient)
	sectionMgrService := grpc.NewSectionManagerService(agentPool, kb, toolRegistry, llmClient)
	employeeService := grpc.NewEmployeeService(agentPool, kb, toolRegistry, llmClient)
	fmt.Println("✓ All services created")
	fmt.Println()
	
	// 3. President: Plan the project
	fmt.Println("3. President Agent: Planning project...")
	ctx := context.Background()
	tasks, err := presidentService.PlanProject(ctx, 
		"E-Commerce Platform",
		"Build a modern e-commerce platform with user management, product catalog, and payment processing",
		[]string{"Deadline: 3 months", "Budget: $100k"},
	)
	if err != nil {
		log.Fatalf("Failed to plan project: %v", err)
	}
	fmt.Printf("✓ Created %d high-level tasks\n", len(tasks))
	for _, task := range tasks {
		fmt.Printf("  - %s: %s\n", task.ID, task.Title)
	}
	fmt.Println()
	
	// 4. Department Manager: Divide tasks into sections
	fmt.Println("4. Department Manager: Dividing tasks into sections...")
	sectionPlans, err := deptMgrService.DivideTasks(ctx, tasks)
	if err != nil {
		log.Fatalf("Failed to divide tasks: %v", err)
	}
	fmt.Printf("✓ Created %d section plans\n", len(sectionPlans))
	for _, plan := range sectionPlans {
		fmt.Printf("  - Section %s: %d tasks assigned to %s\n", plan.SectionID, len(plan.Tasks), plan.ManagerID)
	}
	fmt.Println()
	
	// 5. Section Manager: Prepare implementation plans
	fmt.Println("5. Section Manager: Preparing implementation specifications...")
	var allSpecs []grpc.ImplementationSpec
	for _, plan := range sectionPlans {
		specs, err := sectionMgrService.PrepareImplementationPlan(ctx, plan)
		if err != nil {
			log.Printf("Warning: Failed to prepare plan for %s: %v", plan.SectionID, err)
			continue
		}
		allSpecs = append(allSpecs, specs...)
		fmt.Printf("✓ Created %d implementation specs for %s\n", len(specs), plan.SectionName)
	}
	fmt.Println()
	
	// 6. Employee: Execute tasks
	fmt.Println("6. Employee Agents: Executing tasks...")
	for _, spec := range allSpecs {
		result, err := employeeService.ExecuteTask(ctx, spec)
		if err != nil {
			log.Printf("Warning: Task execution failed: %v", err)
			continue
		}
		fmt.Printf("✓ Task %s completed: %s\n", result.TaskID, result.Status)
	}
	fmt.Println()
	
	// 7. Show final system state
	fmt.Println("7. Final System State:")
	fmt.Println("\nAgent Status:")
	statuses := agentPool.GetAllStatus()
	for _, status := range statuses {
		fmt.Printf("  %s (%s): %s - %s\n", 
			status.AgentID, 
			status.AgentType, 
			status.State,
			status.Message,
		)
	}
	
	fmt.Println("\nKnowledge Base Entries:")
	entries, _ := kb.List(ctx)
	fmt.Printf("  Total entries: %d\n", len(entries))
	for i, entry := range entries {
		if i < 5 { // Show first 5
			fmt.Printf("  - %s: %s\n", entry.Key, truncate(entry.Value, 50))
		}
	}
	if len(entries) > 5 {
		fmt.Printf("  ... and %d more entries\n", len(entries)-5)
	}
	
	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nThis demonstrates the full agent hierarchy workflow:")
	fmt.Println("  1. President plans the project")
	fmt.Println("  2. Department Manager divides into sections")
	fmt.Println("  3. Section Managers create detailed specs")
	fmt.Println("  4. Employees execute the tasks")
	fmt.Println("  5. All information is stored in the shared knowledge base")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
