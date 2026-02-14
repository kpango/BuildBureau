package agent

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/pkg/types"
)

// Organization manages the entire agent hierarchy.
type Organization struct {
	president   types.Agent
	config      *types.Config
	secretaries map[string]types.Agent
	llmManager  *llm.Manager
	directors   []types.Agent
	managers    []types.Agent
	engineers   []types.Agent
}

// NewOrganization creates a new organization from configuration.
func NewOrganization(cfg *types.Config) (*Organization, error) {
	org := &Organization{
		config:      cfg,
		directors:   make([]types.Agent, 0),
		managers:    make([]types.Agent, 0),
		engineers:   make([]types.Agent, 0),
		secretaries: make(map[string]types.Agent),
	}

	// Initialize LLM manager
	llmMgr, err := llm.NewManager(&cfg.LLMs)
	if err != nil {
		fmt.Printf("Warning: Failed to initialize LLM manager: %v\n", err)
		fmt.Println("Agents will work without LLM assistance")
	} else {
		org.llmManager = llmMgr
		fmt.Println("✓ LLM manager initialized successfully")
	}

	if err := org.buildHierarchy(); err != nil {
		return nil, fmt.Errorf("failed to build hierarchy: %w", err)
	}

	return org, nil
}

// buildHierarchy creates and connects all agents based on configuration.
func (o *Organization) buildHierarchy() error {
	loader := config.NewLoader()

	// Create agents for each layer
	for _, layer := range o.config.Organization.Layers {
		switch layer.Name {
		case "President":
			if layer.Agent != "" {
				agentCfg, err := loader.LoadAgentConfig(layer.Agent)
				if err != nil {
					return fmt.Errorf("failed to load president config: %w", err)
				}
				o.president = NewPresidentAgent("president-1", agentCfg)
			}

		case "Director":
			if layer.Agent != "" {
				agentCfg, err := loader.LoadAgentConfig(layer.Agent)
				if err != nil {
					return fmt.Errorf("failed to load director config: %w", err)
				}
				count := layer.Count
				if count == 0 {
					count = 1
				}
				for i := 0; i < count; i++ {
					director := NewDirectorAgent(fmt.Sprintf("director-%d", i+1), agentCfg)
					o.directors = append(o.directors, director)
				}
			}

		case "Manager":
			if layer.Agent != "" {
				agentCfg, err := loader.LoadAgentConfig(layer.Agent)
				if err != nil {
					return fmt.Errorf("failed to load manager config: %w", err)
				}
				count := layer.Count
				if count == 0 {
					count = 1
				}
				for i := 0; i < count; i++ {
					manager := NewManagerAgent(fmt.Sprintf("manager-%d", i+1), agentCfg, o.llmManager)
					o.managers = append(o.managers, manager)
				}
			}

		case "Engineer":
			if layer.Agent != "" {
				agentCfg, err := loader.LoadAgentConfig(layer.Agent)
				if err != nil {
					return fmt.Errorf("failed to load engineer config: %w", err)
				}
				count := layer.Count
				if count == 0 {
					count = 1
				}
				for i := 0; i < count; i++ {
					engineer := NewEngineerAgent(fmt.Sprintf("engineer-%d", i+1), agentCfg, o.llmManager)
					o.engineers = append(o.engineers, engineer)
				}
			}

		case "Secretary":
			if layer.Agent != "" {
				agentCfg, err := loader.LoadAgentConfig(layer.Agent)
				if err != nil {
					return fmt.Errorf("failed to load secretary config: %w", err)
				}
				// Create secretaries for each specified attachment point
				for _, attachTo := range layer.AttachTo {
					secretary := NewSecretaryAgent(fmt.Sprintf("secretary-%s", attachTo), agentCfg)
					o.secretaries[attachTo] = secretary
				}
			}
		}
	}

	// Wire up the hierarchy
	return o.wireHierarchy()
}

// wireHierarchy connects agents to their subordinates and secretaries.
func (o *Organization) wireHierarchy() error {
	// Attach secretaries
	if presidentSecretary, ok := o.secretaries["President"]; ok {
		if president, ok := o.president.(*PresidentAgent); ok {
			president.SetSecretary(presidentSecretary)
		}

		// President's secretary connects to directors
		if secretary, ok := presidentSecretary.(*SecretaryAgent); ok {
			for _, director := range o.directors {
				secretary.AddDirector(director)
			}
		}
	}

	// Wire directors to managers
	for _, director := range o.directors {
		if directorAgent, ok := director.(*DirectorAgent); ok {
			if directorSecretary, ok := o.secretaries["Director"]; ok {
				directorAgent.SetSecretary(directorSecretary)
			}
			// Add managers to each director
			for _, manager := range o.managers {
				directorAgent.AddManager(manager)
			}
		}
	}

	// Wire managers to engineers
	for _, manager := range o.managers {
		if managerAgent, ok := manager.(*ManagerAgent); ok {
			if managerSecretary, ok := o.secretaries["Manager"]; ok {
				managerAgent.SetSecretary(managerSecretary)
			}
			// Add engineers to each manager
			for _, engineer := range o.engineers {
				managerAgent.AddEngineer(engineer)
			}
		}
	}

	return nil
}

// Start initializes all agents in the organization.
func (o *Organization) Start(ctx context.Context) error {
	agents := []types.Agent{}

	if o.president != nil {
		agents = append(agents, o.president)
	}
	for _, secretary := range o.secretaries {
		agents = append(agents, secretary)
	}
	agents = append(agents, o.directors...)
	agents = append(agents, o.managers...)
	agents = append(agents, o.engineers...)

	for _, agent := range agents {
		if err := agent.Start(ctx); err != nil {
			return fmt.Errorf("failed to start agent %s: %w", agent.GetID(), err)
		}
	}

	return nil
}

// Stop gracefully shuts down all agents.
func (o *Organization) Stop(ctx context.Context) error {
	agents := []types.Agent{}

	agents = append(agents, o.engineers...)
	agents = append(agents, o.managers...)
	agents = append(agents, o.directors...)
	for _, secretary := range o.secretaries {
		agents = append(agents, secretary)
	}
	if o.president != nil {
		agents = append(agents, o.president)
	}

	for _, agent := range agents {
		if err := agent.Stop(ctx); err != nil {
			return fmt.Errorf("failed to stop agent %s: %w", agent.GetID(), err)
		}
	}

	// Close LLM manager
	if o.llmManager != nil {
		if err := o.llmManager.Close(); err != nil {
			fmt.Printf("Warning: failed to close LLM manager: %v\n", err)
		}
	}

	return nil
}

// GetPresident returns the president agent.
func (o *Organization) GetPresident() types.Agent {
	return o.president
}

// ProcessClientTask processes a task from the client through the president.
func (o *Organization) ProcessClientTask(ctx context.Context, instruction string) (*types.TaskResponse, error) {
	if o.president == nil {
		return nil, fmt.Errorf("no president agent available")
	}

	task := &types.Task{
		ID:          uuid.New().String(),
		Title:       "Client Request",
		Description: instruction,
		FromAgent:   "client",
		ToAgent:     o.president.GetID(),
		Content:     instruction,
		Priority:    1,
	}

	return o.president.ProcessTask(ctx, task)
}
