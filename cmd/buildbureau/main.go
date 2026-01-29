package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kpango/BuildBureau/internal/agent"
	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/internal/tui"
)

const (
	defaultConfigPath = "config.yaml"
)

func main() {
	// Get config path from environment or use default
	configPath := os.Getenv("BUILDBUREAU_CONFIG")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	// Load configuration
	loader := config.NewLoader()
	cfg, err := loader.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create organization
	org, err := agent.NewOrganization(cfg)
	if err != nil {
		log.Fatalf("Failed to create organization: %v", err)
	}

	// Start organization
	ctx := context.Background()
	if err := org.Start(ctx); err != nil {
		log.Fatalf("Failed to start organization: %v", err)
	}
	defer func() {
		if err := org.Stop(ctx); err != nil {
			log.Printf("Error stopping organization: %v", err)
		}
	}()

	// Start TUI
	p := tea.NewProgram(
		tui.NewModel(org),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
