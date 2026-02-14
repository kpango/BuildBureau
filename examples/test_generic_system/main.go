package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kpango/BuildBureau/internal/agent"
	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/pkg/types"
)

func main() {
	fmt.Println("=== BuildBureau Generic Agent System Demo ===\n")

	// Load configuration
	loader := config.NewLoader()
	cfg, err := loader.Load("./config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create generic organization
	fmt.Println("Creating generic organization...")
	org, err := agent.NewGenericOrganization(cfg)
	if err != nil {
		log.Fatalf("Failed to create organization: %v", err)
	}
	fmt.Println("✓ Organization created")

	// Start the organization
	ctx := context.Background()
	if err := org.Start(ctx); err != nil {
		log.Fatalf("Failed to start organization: %v", err)
	}
	fmt.Println("✓ Organization started")
	defer org.Stop(ctx)

	// Display organization structure
	fmt.Println("\n--- Organization Structure ---")
	displayOrganizationStructure(org)

	// Submit a test task
	fmt.Println("\n--- Processing Task ---")
	task := &types.Task{
		ID:          "demo-task-1",
		Title:       "Build a REST API",
		Description: "Create a RESTful API service with user authentication",
		Content:     "The API should include endpoints for user registration, login, and profile management",
		Priority:    1,
	}

	fmt.Printf("Submitting task: %s\n", task.Title)
	response, err := org.ProcessTask(ctx, task)
	if err != nil {
		log.Fatalf("Task processing failed: %v", err)
	}

	fmt.Printf("\nTask Status: %s\n", response.Status)
	fmt.Printf("Task Result:\n%s\n", response.Result)

	// Display organization status
	fmt.Println("\n--- Organization Status ---")
	displayStatus(org)
}

func displayOrganizationStructure(org *agent.GenericOrganization) {
	allAgents := org.GetAllAgents()
	roleGroups := make(map[types.AgentRole][]string)

	for id, agent := range allAgents {
		role := agent.GetRole()
		roleGroups[role] = append(roleGroups[role], id)
	}

	// Display by role
	roles := []types.AgentRole{
		types.RolePresident,
		types.RoleSecretary,
		types.RoleDirector,
		types.RoleManager,
		types.RoleEngineer,
	}

	for _, role := range roles {
		if agents, ok := roleGroups[role]; ok {
			fmt.Printf("%s (%d):\n", role, len(agents))
			for _, id := range agents {
				fmt.Printf("  - %s\n", id)
			}
		}
	}
}

func displayStatus(org *agent.GenericOrganization) {
	status := org.GetStatus()

	for id, stats := range status {
		fmt.Printf("%s (%s):\n", id, stats["role"])
		fmt.Printf("  Running: %v\n", stats["running"])
		fmt.Printf("  Active tasks: %d\n", stats["active"])
		fmt.Printf("  Completed tasks: %d\n", stats["completed"])
	}
}
