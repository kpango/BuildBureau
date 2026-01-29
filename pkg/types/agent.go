package types

import (
	"context"
)

// AgentRole represents the role of an agent in the organization.
type AgentRole string

const (
	RolePresident AgentRole = "President"
	RoleSecretary AgentRole = "Secretary"
	RoleDirector  AgentRole = "Director"
	RoleManager   AgentRole = "Manager"
	RoleEngineer  AgentRole = "Engineer"
)

// Agent represents the core interface that all agents must implement.
type Agent interface {
	// GetID returns the unique identifier for this agent
	GetID() string

	// GetRole returns the role of this agent
	GetRole() AgentRole

	// ProcessTask handles an incoming task and returns a response
	ProcessTask(ctx context.Context, task *Task) (*TaskResponse, error)

	// Start initializes the agent and starts any background processes
	Start(ctx context.Context) error

	// Stop gracefully shuts down the agent
	Stop(ctx context.Context) error
}

// Task represents a work item passed between agents.
type Task struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	FromAgent   string            `json:"from_agent"`
	ToAgent     string            `json:"to_agent"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Content     string            `json:"content"`
	Priority    int               `json:"priority"`
}

// TaskResponse represents the response from an agent after processing a task.
type TaskResponse struct {
	TaskID   string            `json:"task_id"`
	Status   TaskStatus        `json:"status"`
	Result   string            `json:"result"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Error    string            `json:"error,omitempty"`
}

// TaskStatus represents the status of a task.
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
	StatusDelegated  TaskStatus = "delegated"
)
