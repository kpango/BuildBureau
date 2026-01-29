package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/kpango/BuildBureau/internal/config"
)

// AgentType represents the type of agent
type AgentType string

const (
	AgentTypePresident          AgentType = "president"
	AgentTypePresidentSecretary AgentType = "president_secretary"
	AgentTypeDepartmentManager  AgentType = "department_manager"
	AgentTypeDepartmentSecretary AgentType = "department_secretary"
	AgentTypeSectionManager     AgentType = "section_manager"
	AgentTypeSectionSecretary   AgentType = "section_secretary"
	AgentTypeEmployee           AgentType = "employee"
)

// Agent represents a base AI agent interface
type Agent interface {
	// ID returns the unique identifier of the agent
	ID() string

	// Type returns the type of the agent
	Type() AgentType

	// Process processes a task with the given context
	Process(ctx context.Context, input interface{}) (interface{}, error)

	// GetStatus returns the current status of the agent
	GetStatus() Status
}

// Status represents the current status of an agent
type Status struct {
	AgentID    string
	AgentType  AgentType
	State      string
	CurrentTask string
	Message    string
}

// BaseAgent provides common functionality for all agents
type BaseAgent struct {
	id         string
	agentType  AgentType
	config     config.AgentConfig
	status     Status
	statusLock sync.RWMutex
}

// NewBaseAgent creates a new base agent
func NewBaseAgent(id string, agentType AgentType, cfg config.AgentConfig) *BaseAgent {
	return &BaseAgent{
		id:        id,
		agentType: agentType,
		config:    cfg,
		status: Status{
			AgentID:   id,
			AgentType: agentType,
			State:     "idle",
		},
	}
}

// ID returns the agent ID
func (a *BaseAgent) ID() string {
	return a.id
}

// Type returns the agent type
func (a *BaseAgent) Type() AgentType {
	return a.agentType
}

// GetStatus returns the current status
func (a *BaseAgent) GetStatus() Status {
	a.statusLock.RLock()
	defer a.statusLock.RUnlock()
	return a.status
}

// UpdateStatus updates the agent status
func (a *BaseAgent) UpdateStatus(state, currentTask, message string) {
	a.statusLock.Lock()
	defer a.statusLock.Unlock()
	a.status.State = state
	a.status.CurrentTask = currentTask
	a.status.Message = message
}

// Process implements the Agent interface - base implementation
func (a *BaseAgent) Process(ctx context.Context, input interface{}) (interface{}, error) {
	// Base implementation - to be overridden by specific agent types
	// For now, this is a placeholder that concrete agents should override
	return nil, fmt.Errorf("Process method not implemented for agent type %s", a.agentType)
}

// AgentPool manages a pool of agents
type AgentPool struct {
	agents     map[string]Agent
	agentsByType map[AgentType][]Agent
	mu         sync.RWMutex
}

// NewAgentPool creates a new agent pool
func NewAgentPool() *AgentPool {
	return &AgentPool{
		agents:       make(map[string]Agent),
		agentsByType: make(map[AgentType][]Agent),
	}
}

// Register registers an agent in the pool
func (p *AgentPool) Register(agent Agent) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.agents[agent.ID()]; exists {
		return fmt.Errorf("agent with ID %s already exists", agent.ID())
	}

	p.agents[agent.ID()] = agent
	p.agentsByType[agent.Type()] = append(p.agentsByType[agent.Type()], agent)
	return nil
}

// Get retrieves an agent by ID
func (p *AgentPool) Get(id string) (Agent, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	agent, exists := p.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent with ID %s not found", id)
	}
	return agent, nil
}

// GetByType retrieves all agents of a specific type
func (p *AgentPool) GetByType(agentType AgentType) []Agent {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.agentsByType[agentType]
}

// GetAvailable returns the first available agent of the given type
func (p *AgentPool) GetAvailable(agentType AgentType) (Agent, error) {
	agents := p.GetByType(agentType)
	if len(agents) == 0 {
		return nil, fmt.Errorf("no agents of type %s found", agentType)
	}

	// Return first idle agent, or first agent if none are idle
	for _, agent := range agents {
		if status := agent.GetStatus(); status.State == "idle" {
			return agent, nil
		}
	}

	// Return first agent if none are idle
	return agents[0], nil
}

// GetAllStatus returns status of all agents
func (p *AgentPool) GetAllStatus() []Status {
	p.mu.RLock()
	defer p.mu.RUnlock()

	statuses := make([]Status, 0, len(p.agents))
	for _, agent := range p.agents {
		statuses = append(statuses, agent.GetStatus())
	}
	return statuses
}
