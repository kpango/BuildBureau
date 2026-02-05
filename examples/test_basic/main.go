package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kpango/BuildBureau/internal/agent"
	"github.com/kpango/BuildBureau/internal/config"
)

func main() {
	// Load configuration
	loader := config.NewLoader()
	cfg, err := loader.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Println("✓ Configuration loaded successfully")
	fmt.Printf("  - Default LLM model: %s\n", cfg.LLMs.DefaultModel)
	fmt.Printf("  - Organization layers: %d\n", len(cfg.Organization.Layers))

	// Create organization
	org, err := agent.NewOrganization(cfg)
	if err != nil {
		log.Fatalf("Failed to create organization: %v", err)
	}

	fmt.Println("✓ Organization created successfully")

	// Start organization
	ctx := context.Background()
	if err := org.Start(ctx); err != nil {
		log.Fatalf("Failed to start organization: %v", err)
	}
	defer org.Stop(ctx)

	fmt.Println("✓ Organization started successfully")

	// Test task processing
	fmt.Println("\n--- Testing Task Delegation ---")
	response, err := org.ProcessClientTask(ctx, "Create a simple REST API for user management")
	if err != nil {
		log.Fatalf("Failed to process task: %v", err)
	}

	fmt.Println("\n--- Task Result ---")
	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("Result:\n%s\n", response.Result)

	fmt.Println("\n✓ All tests passed!")
}
