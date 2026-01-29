package types

import "time"

// AgentRole represents the role of an agent in the hierarchy
type AgentRole string

const (
	RoleCEO       AgentRole = "CEO"
	RoleManager   AgentRole = "Manager"
	RoleLead      AgentRole = "Lead"
	RoleEmployee  AgentRole = "Employee"
	RoleSecretary AgentRole = "Secretary"
	RoleClient    AgentRole = "Client"
)

// Task represents a work item to be processed
type Task struct {
	ID          string
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      TaskStatus
	AssignedTo  AgentRole
	CreatedBy   AgentRole
	Result      string
	SubTasks    []Task
}

// TaskStatus represents the current state of a task
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
)

// Message represents communication between agents
type Message struct {
	ID        string
	From      AgentRole
	To        AgentRole
	Content   string
	Timestamp time.Time
	TaskID    string
}

// AgentEvent represents an event in the system
type AgentEvent struct {
	Type      EventType
	Agent     AgentRole
	Message   string
	Timestamp time.Time
	TaskID    string
}

// EventType represents different types of events
type EventType string

const (
	EventTaskAssigned  EventType = "task_assigned"
	EventTaskCompleted EventType = "task_completed"
	EventTaskStarted   EventType = "task_started"
	EventError         EventType = "error"
	EventMessage       EventType = "message"
)
