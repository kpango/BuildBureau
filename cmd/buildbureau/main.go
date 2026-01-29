package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kpango/BuildBureau/internal/agent"
	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/internal/slack"
	"github.com/kpango/BuildBureau/internal/ui"
)

const (
	defaultConfigPath = "config.yaml"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func run() error {
	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Received shutdown signal")
		cancel()
	}()

	// Initialize Slack notifier
	notifier, err := slack.NewNotifier(cfg.Slack)
	if err != nil {
		return fmt.Errorf("failed to initialize Slack notifier: %w", err)
	}

	if notifier.IsEnabled() {
		log.Println("Slack notifications enabled")
	}

	// Initialize agent pool
	agentPool := agent.NewAgentPool()

	// Create agents based on configuration
	if err := initializeAgents(agentPool, cfg); err != nil {
		return fmt.Errorf("failed to initialize agents: %w", err)
	}

	log.Printf("Initialized %d agents across all types\n", len(agentPool.GetAllStatus()))

	// Initialize and run UI if enabled
	if cfg.UI.EnableTUI {
		model := ui.NewModel()
		p := tea.NewProgram(model, tea.WithAltScreen())

		go func() {
			<-ctx.Done()
			p.Send(tea.Quit())
		}()

		if _, err := p.Run(); err != nil {
			return fmt.Errorf("failed to run UI: %w", err)
		}
	} else {
		// Run in CLI mode
		log.Println("Running in CLI mode (TUI disabled)")
		// CLI mode implementation would go here
		<-ctx.Done()
	}

	return nil
}

// initializeAgents creates and registers all agents based on configuration
func initializeAgents(pool *agent.AgentPool, cfg *config.Config) error {
	// Create President agents
	for i := 0; i < cfg.Agents.President.Count; i++ {
		id := fmt.Sprintf("president-%d", i+1)
		baseAgent := agent.NewBaseAgent(id, agent.AgentTypePresident, cfg.Agents.President)
		if err := pool.Register(baseAgent); err != nil {
			return err
		}
	}

	// Create President Secretary agents
	for i := 0; i < cfg.Agents.PresidentSecretary.Count; i++ {
		id := fmt.Sprintf("president-secretary-%d", i+1)
		baseAgent := agent.NewBaseAgent(id, agent.AgentTypePresidentSecretary, cfg.Agents.PresidentSecretary)
		if err := pool.Register(baseAgent); err != nil {
			return err
		}
	}

	// Create Department Manager agents
	for i := 0; i < cfg.Agents.DepartmentManager.Count; i++ {
		id := fmt.Sprintf("dept-manager-%d", i+1)
		baseAgent := agent.NewBaseAgent(id, agent.AgentTypeDepartmentManager, cfg.Agents.DepartmentManager)
		if err := pool.Register(baseAgent); err != nil {
			return err
		}
	}

	// Create Department Secretary agents
	for i := 0; i < cfg.Agents.DepartmentSecretary.Count; i++ {
		id := fmt.Sprintf("dept-secretary-%d", i+1)
		baseAgent := agent.NewBaseAgent(id, agent.AgentTypeDepartmentSecretary, cfg.Agents.DepartmentSecretary)
		if err := pool.Register(baseAgent); err != nil {
			return err
		}
	}

	// Create Section Manager agents
	for i := 0; i < cfg.Agents.SectionManager.Count; i++ {
		id := fmt.Sprintf("section-manager-%d", i+1)
		baseAgent := agent.NewBaseAgent(id, agent.AgentTypeSectionManager, cfg.Agents.SectionManager)
		if err := pool.Register(baseAgent); err != nil {
			return err
		}
	}

	// Create Section Secretary agents
	for i := 0; i < cfg.Agents.SectionSecretary.Count; i++ {
		id := fmt.Sprintf("section-secretary-%d", i+1)
		baseAgent := agent.NewBaseAgent(id, agent.AgentTypeSectionSecretary, cfg.Agents.SectionSecretary)
		if err := pool.Register(baseAgent); err != nil {
			return err
		}
	}

	// Create Employee agents
	for i := 0; i < cfg.Agents.Employee.Count; i++ {
		id := fmt.Sprintf("employee-%d", i+1)
		baseAgent := agent.NewBaseAgent(id, agent.AgentTypeEmployee, cfg.Agents.Employee)
		if err := pool.Register(baseAgent); err != nil {
			return err
		}
	}

	return nil
}
