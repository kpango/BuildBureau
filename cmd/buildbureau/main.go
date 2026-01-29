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
	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/internal/slack"
	"github.com/kpango/BuildBureau/internal/tools"
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

	// Initialize LLM client
	apiKey := os.Getenv("GOOGLE_AI_API_KEY")
	factory := llm.NewClientFactory(cfg.LLM.Provider, apiKey, cfg.LLM.DefaultModel)
	llmClient, err := factory.Create()
	if err != nil {
		log.Printf("Warning: Failed to initialize LLM client: %v. Using mock client.", err)
		llmClient = llm.NewMockClient(nil)
	}

	// Initialize tool registry
	toolRegistry := tools.NewDefaultRegistry()

	// Create agents based on configuration
	if err := initializeAgents(agentPool, cfg, llmClient, toolRegistry); err != nil {
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
func initializeAgents(pool *agent.AgentPool, cfg *config.Config, llmClient llm.Client, toolRegistry *tools.Registry) error {
	// Create President agents with specialized capabilities
	for i := 0; i < cfg.Agents.President.Count; i++ {
		id := fmt.Sprintf("president-%d", i+1)
		specializedAgent := agent.NewSpecializedAgent(id, agent.AgentTypePresident, cfg.Agents.President, llmClient, toolRegistry)
		if err := pool.Register(specializedAgent); err != nil {
			return err
		}
	}

	// Create President Secretary agents
	for i := 0; i < cfg.Agents.PresidentSecretary.Count; i++ {
		id := fmt.Sprintf("president-secretary-%d", i+1)
		specializedAgent := agent.NewSpecializedAgent(id, agent.AgentTypePresidentSecretary, cfg.Agents.PresidentSecretary, llmClient, toolRegistry)
		if err := pool.Register(specializedAgent); err != nil {
			return err
		}
	}

	// Create Department Manager agents
	for i := 0; i < cfg.Agents.DepartmentManager.Count; i++ {
		id := fmt.Sprintf("dept-manager-%d", i+1)
		specializedAgent := agent.NewSpecializedAgent(id, agent.AgentTypeDepartmentManager, cfg.Agents.DepartmentManager, llmClient, toolRegistry)
		if err := pool.Register(specializedAgent); err != nil {
			return err
		}
	}

	// Create Department Secretary agents
	for i := 0; i < cfg.Agents.DepartmentSecretary.Count; i++ {
		id := fmt.Sprintf("dept-secretary-%d", i+1)
		specializedAgent := agent.NewSpecializedAgent(id, agent.AgentTypeDepartmentSecretary, cfg.Agents.DepartmentSecretary, llmClient, toolRegistry)
		if err := pool.Register(specializedAgent); err != nil {
			return err
		}
	}

	// Create Section Manager agents
	for i := 0; i < cfg.Agents.SectionManager.Count; i++ {
		id := fmt.Sprintf("section-manager-%d", i+1)
		specializedAgent := agent.NewSpecializedAgent(id, agent.AgentTypeSectionManager, cfg.Agents.SectionManager, llmClient, toolRegistry)
		if err := pool.Register(specializedAgent); err != nil {
			return err
		}
	}

	// Create Section Secretary agents
	for i := 0; i < cfg.Agents.SectionSecretary.Count; i++ {
		id := fmt.Sprintf("section-secretary-%d", i+1)
		specializedAgent := agent.NewSpecializedAgent(id, agent.AgentTypeSectionSecretary, cfg.Agents.SectionSecretary, llmClient, toolRegistry)
		if err := pool.Register(specializedAgent); err != nil {
			return err
		}
	}

	// Create Employee agents
	for i := 0; i < cfg.Agents.Employee.Count; i++ {
		id := fmt.Sprintf("employee-%d", i+1)
		specializedAgent := agent.NewSpecializedAgent(id, agent.AgentTypeEmployee, cfg.Agents.Employee, llmClient, toolRegistry)
		if err := pool.Register(specializedAgent); err != nil {
			return err
		}
	}

	return nil
}
