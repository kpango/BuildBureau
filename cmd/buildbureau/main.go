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
	"github.com/kpango/BuildBureau/pkg/types"
)

func main() {
	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		log.Printf("Warning: Could not load config file, using defaults: %v", err)
		cfg = config.Default()
	}

	// Create event channel for agent events
	eventChan := make(chan types.AgentEvent, 1000)

	// Create task channel for user input
	taskChan := make(chan types.Task, 10)

	// Initialize Slack notifier
	slackNotifier, err := slack.NewNotifier(cfg.Slack)
	if err != nil {
		log.Fatalf("Failed to initialize Slack notifier: %v", err)
	}

	// Create agent hierarchy
	agents := setupAgents(cfg, eventChan)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start event processor
	go processEvents(eventChan, slackNotifier)

	// Start task processor
	go processTasks(ctx, taskChan, agents.ceo)

	// Start all agents
	for _, a := range getAllAgents(agents) {
		if err := a.Start(ctx); err != nil {
			log.Printf("Failed to start agent %s: %v", a.Role(), err)
		}
	}

	// Create and run TUI
	tuiModel := ui.NewModel(taskChan)
	program := tea.NewProgram(tuiModel, tea.WithAltScreen())

	// Forward events to TUI
	go func() {
		for event := range eventChan {
			program.Send(ui.EventMsg(event))
		}
	}()

	// Handle shutdown
	go func() {
		<-sigChan
		cancel()
		program.Quit()
	}()

	// Run the TUI
	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}

	// Cleanup
	for _, a := range getAllAgents(agents) {
		if err := a.Stop(); err != nil {
			log.Printf("Error stopping agent %s: %v", a.Role(), err)
		}
	}
}

// AgentHierarchy holds all agents in the system
type AgentHierarchy struct {
	ceo                *agent.CEOAgent
	ceoSecretary       *agent.SecretaryAgent
	managers           []*agent.ManagerAgent
	managerSecretaries []*agent.SecretaryAgent
	leads              []*agent.LeadAgent
	leadSecretaries    []*agent.SecretaryAgent
	employees          []*agent.EmployeeAgent
}

// setupAgents creates and configures all agents
func setupAgents(cfg *config.Config, eventChan chan types.AgentEvent) *AgentHierarchy {
	hierarchy := &AgentHierarchy{}

	// Create CEO and secretary
	hierarchy.ceo = agent.NewCEOAgent("Chief Executive Officer", eventChan)
	hierarchy.ceoSecretary = agent.NewSecretaryAgent("CEO Secretary", hierarchy.ceo, eventChan)
	hierarchy.ceo.SetSecretary(hierarchy.ceoSecretary)

	// Create managers with secretaries
	for i := 0; i < 2; i++ {
		mgr := agent.NewManagerAgent(fmt.Sprintf("Manager %d", i+1), eventChan)
		mgrSecretary := agent.NewSecretaryAgent(fmt.Sprintf("Manager %d Secretary", i+1), mgr, eventChan)
		mgr.SetSecretary(mgrSecretary)

		hierarchy.managers = append(hierarchy.managers, mgr)
		hierarchy.managerSecretaries = append(hierarchy.managerSecretaries, mgrSecretary)
		hierarchy.ceo.AddManagerAgent(mgr)

		// Create leads with secretaries for each manager
		for j := 0; j < 2; j++ {
			lead := agent.NewLeadAgent(fmt.Sprintf("Lead %d-%d", i+1, j+1), eventChan)
			leadSecretary := agent.NewSecretaryAgent(fmt.Sprintf("Lead %d-%d Secretary", i+1, j+1), lead, eventChan)
			lead.SetSecretary(leadSecretary)

			hierarchy.leads = append(hierarchy.leads, lead)
			hierarchy.leadSecretaries = append(hierarchy.leadSecretaries, leadSecretary)
			mgr.AddLeadAgent(lead)

			// Create employees for each lead
			for k := 0; k < 2; k++ {
				emp := agent.NewEmployeeAgent(fmt.Sprintf("Employee %d-%d-%d", i+1, j+1, k+1), eventChan)
				hierarchy.employees = append(hierarchy.employees, emp)
				lead.AddEmployeeAgent(emp)
			}
		}
	}

	return hierarchy
}

// getAllAgents returns all agents as a slice
func getAllAgents(h *AgentHierarchy) []agent.Agent {
	agents := []agent.Agent{h.ceo, h.ceoSecretary}

	for _, m := range h.managers {
		agents = append(agents, m)
	}
	for _, s := range h.managerSecretaries {
		agents = append(agents, s)
	}
	for _, l := range h.leads {
		agents = append(agents, l)
	}
	for _, s := range h.leadSecretaries {
		agents = append(agents, s)
	}
	for _, e := range h.employees {
		agents = append(agents, e)
	}

	return agents
}

// processEvents handles agent events
func processEvents(eventChan chan types.AgentEvent, notifier *slack.Notifier) {
	for event := range eventChan {
		// Send to Slack if enabled
		if notifier.IsEnabled() {
			if err := notifier.NotifyEvent(event); err != nil {
				log.Printf("Failed to send Slack notification: %v", err)
			}
		}
	}
}

// processTasks handles incoming tasks from the UI
func processTasks(ctx context.Context, taskChan chan types.Task, ceo *agent.CEOAgent) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-taskChan:
			// CEO handles the task
			go func(t types.Task) {
				_, err := ceo.HandleTask(ctx, t)
				if err != nil {
					log.Printf("Error handling task: %v", err)
				}
			}(task)
		}
	}
}

// loadConfig loads the configuration file
func loadConfig() (*config.Config, error) {
	configPath := "config.yaml"
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		configPath = path
	}

	return config.Load(configPath)
}
