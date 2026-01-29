package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kpango/BuildBureau/pkg/types"
)

// Agent represents a base agent in the system
type Agent interface {
	// ID returns the unique identifier of the agent
	ID() string

	// Role returns the role of the agent
	Role() types.AgentRole

	// HandleTask processes a task
	HandleTask(ctx context.Context, task types.Task) (types.Task, error)

	// SendMessage sends a message to another agent
	SendMessage(ctx context.Context, to types.AgentRole, content string) error

	// Start starts the agent
	Start(ctx context.Context) error

	// Stop stops the agent
	Stop() error
}

// BaseAgent provides common functionality for all agents
type BaseAgent struct {
	id          string
	role        types.AgentRole
	name        string
	tasks       map[string]types.Task
	mu          sync.RWMutex
	eventChan   chan types.AgentEvent
	messageChan chan types.Message
	secretary   Agent // Each agent can have a secretary
}

// NewBaseAgent creates a new base agent
func NewBaseAgent(role types.AgentRole, name string, eventChan chan types.AgentEvent) *BaseAgent {
	return &BaseAgent{
		id:          uuid.New().String(),
		role:        role,
		name:        name,
		tasks:       make(map[string]types.Task),
		eventChan:   eventChan,
		messageChan: make(chan types.Message, 100),
	}
}

// ID returns the agent's ID
func (a *BaseAgent) ID() string {
	return a.id
}

// Role returns the agent's role
func (a *BaseAgent) Role() types.AgentRole {
	return a.role
}

// Name returns the agent's name
func (a *BaseAgent) Name() string {
	return a.name
}

// SetSecretary sets the secretary agent for this agent
func (a *BaseAgent) SetSecretary(secretary Agent) {
	a.secretary = secretary
}

// GetSecretary returns the secretary agent
func (a *BaseAgent) GetSecretary() Agent {
	return a.secretary
}

// AddTask adds a task to the agent's task list
func (a *BaseAgent) AddTask(task types.Task) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.tasks[task.ID] = task
}

// GetTask retrieves a task by ID
func (a *BaseAgent) GetTask(taskID string) (types.Task, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	task, ok := a.tasks[taskID]
	return task, ok
}

// UpdateTask updates a task's status
func (a *BaseAgent) UpdateTask(taskID string, status types.TaskStatus, result string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	task, ok := a.tasks[taskID]
	if !ok {
		return fmt.Errorf("task not found: %s", taskID)
	}

	task.Status = status
	task.Result = result
	task.UpdatedAt = time.Now()
	a.tasks[taskID] = task

	return nil
}

// EmitEvent sends an event to the event channel
func (a *BaseAgent) EmitEvent(eventType types.EventType, message string, taskID string) {
	if a.eventChan != nil {
		event := types.AgentEvent{
			Type:      eventType,
			Agent:     a.role,
			Message:   message,
			Timestamp: time.Now(),
			TaskID:    taskID,
		}
		select {
		case a.eventChan <- event:
		default:
			// Channel full, skip event
		}
	}
}

// SendMessage sends a message to another agent (basic implementation)
func (a *BaseAgent) SendMessage(ctx context.Context, to types.AgentRole, content string) error {
	a.EmitEvent(types.EventMessage, fmt.Sprintf("Sending message to %s: %s", to, content), "")

	// In a real implementation, this would route to the target agent
	// For now, just emit the event
	return nil
}

// Start starts the agent (basic implementation)
func (a *BaseAgent) Start(ctx context.Context) error {
	a.EmitEvent(types.EventMessage, fmt.Sprintf("%s started", a.name), "")
	return nil
}

// Stop stops the agent
func (a *BaseAgent) Stop() error {
	a.EmitEvent(types.EventMessage, fmt.Sprintf("%s stopped", a.name), "")
	return nil
}

// HandleTask is the default implementation (to be overridden)
func (a *BaseAgent) HandleTask(ctx context.Context, task types.Task) (types.Task, error) {
	return task, fmt.Errorf("HandleTask not implemented for %s", a.role)
}
