package agent

import (
	"context"
	"fmt"

	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/internal/llm"
	"github.com/kpango/BuildBureau/pkg/types"
)

// GenericOrganization manages agent hierarchy using generic agents.
// This replaces the role-specific Organization with a flexible, config-driven approach.
type GenericOrganization struct {
	config     *types.Config
	llmManager *llm.Manager
	agents     map[string]*GenericAgent // Map of agent ID to agent instance
	rootAgent  *GenericAgent            // Top-level agent (typically President)
}

// NewGenericOrganization creates a new generic organization from configuration.
func NewGenericOrganization(cfg *types.Config) (*GenericOrganization, error) {
	org := &GenericOrganization{
		config: cfg,
		agents: make(map[string]*GenericAgent),
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
func (o *GenericOrganization) buildHierarchy() error {
	loader := config.NewLoader()

	// Phase 1: Create all agents
	for _, layer := range o.config.Organization.Layers {
		if layer.Agent == "" {
			continue
		}

		// Load agent configuration
		agentCfg, err := loader.LoadAgentConfig(layer.Agent)
		if err != nil {
			return fmt.Errorf("failed to load %s config: %w", layer.Name, err)
		}

		// Determine count (default to 1)
		count := layer.Count
		if count == 0 {
			count = 1
		}

		// Convert layer name to role
		role := types.AgentRole(layer.Name)

		// Create agents for this layer
		for i := 0; i < count; i++ {
			agentID := fmt.Sprintf("%s-%d", layer.Name, i+1)
			agent := NewGenericAgent(agentID, role, agentCfg, o.llmManager)
			o.agents[agentID] = agent

			// Track the root agent (first agent in first layer)
			if o.rootAgent == nil {
				o.rootAgent = agent
			}
		}
	}

	// Phase 2: Build hierarchical relationships
	if err := o.buildRelationships(); err != nil {
		return fmt.Errorf("failed to build relationships: %w", err)
	}

	return nil
}

// buildRelationships establishes parent-child relationships between agents
// based on the organizational structure configuration.
func (o *GenericOrganization) buildRelationships() error {
	// Build a map of layer names to agents in that layer
	layerAgents := make(map[string][]*GenericAgent)
	for _, agent := range o.agents {
		role := string(agent.GetRole())
		layerAgents[role] = append(layerAgents[role], agent)
	}

	// Process each layer to establish relationships
	for i, layer := range o.config.Organization.Layers {
		if layer.Agent == "" {
			continue
		}

		currentLayerAgents := layerAgents[layer.Name]

		// If this is not the last layer, connect to the next layer
		if i < len(o.config.Organization.Layers)-1 {
			nextLayer := o.config.Organization.Layers[i+1]
			if nextLayer.Agent != "" {
				nextLayerAgents := layerAgents[nextLayer.Name]

				// Connect current layer agents to next layer as subordinates
				for _, currentAgent := range currentLayerAgents {
					for _, subordinate := range nextLayerAgents {
						currentAgent.AddSubordinate(subordinate)
						subordinate.SetParent(currentAgent)
					}
				}
			}
		}

		// Handle attach_to relationships for special roles like Secretary
		if len(layer.AttachTo) > 0 {
			for _, attachTo := range layer.AttachTo {
				attachAgents := layerAgents[attachTo]
				for _, attachAgent := range attachAgents {
					// Add this layer's agents as subordinates to attach points
					for _, currentAgent := range currentLayerAgents {
						attachAgent.AddSubordinate(currentAgent)
						currentAgent.SetParent(attachAgent)
					}
				}
			}
		}
	}

	return nil
}

// Start initializes all agents in the organization.
func (o *GenericOrganization) Start(ctx context.Context) error {
	for id, agent := range o.agents {
		if err := agent.Start(ctx); err != nil {
			return fmt.Errorf("failed to start agent %s: %w", id, err)
		}
	}
	return nil
}

// Stop gracefully shuts down all agents.
func (o *GenericOrganization) Stop(ctx context.Context) error {
	for id, agent := range o.agents {
		if err := agent.Stop(ctx); err != nil {
			return fmt.Errorf("failed to stop agent %s: %w", id, err)
		}
	}
	return nil
}

// ProcessTask submits a task to the root agent (typically President).
func (o *GenericOrganization) ProcessTask(ctx context.Context, task *types.Task) (*types.TaskResponse, error) {
	if o.rootAgent == nil {
		return nil, fmt.Errorf("no root agent available")
	}

	return o.rootAgent.ProcessTask(ctx, task)
}

// GetAgent returns an agent by ID.
func (o *GenericOrganization) GetAgent(id string) types.Agent {
	return o.agents[id]
}

// GetAgentsByRole returns all agents with the specified role.
func (o *GenericOrganization) GetAgentsByRole(role types.AgentRole) []types.Agent {
	result := make([]types.Agent, 0)
	for _, agent := range o.agents {
		if agent.GetRole() == role {
			result = append(result, agent)
		}
	}
	return result
}

// GetAllAgents returns all agents in the organization.
func (o *GenericOrganization) GetAllAgents() map[string]types.Agent {
	result := make(map[string]types.Agent)
	for id, agent := range o.agents {
		result[id] = agent
	}
	return result
}

// GetStatus returns the current status of all agents.
func (o *GenericOrganization) GetStatus() map[string]map[string]interface{} {
	status := make(map[string]map[string]interface{})

	for id, agent := range o.agents {
		active, completed := agent.GetStats()
		status[id] = map[string]interface{}{
			"role":      agent.GetRole(),
			"active":    active,
			"completed": completed,
			"running":   agent.IsRunning(),
		}
	}

	return status
}
