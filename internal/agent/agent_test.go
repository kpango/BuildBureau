package agent

import (
	"context"
	"testing"

	"github.com/kpango/BuildBureau/internal/config"
)

func TestNewBaseAgent(t *testing.T) {
	cfg := config.AgentConfig{
		Count:      1,
		Model:      "test-model",
		Instruction: "Test instruction",
		Timeout:    60,
		RetryCount: 3,
	}

	agent := NewBaseAgent("test-1", AgentTypePresident, cfg)

	if agent.ID() != "test-1" {
		t.Errorf("Expected ID 'test-1', got '%s'", agent.ID())
	}

	if agent.Type() != AgentTypePresident {
		t.Errorf("Expected type '%s', got '%s'", AgentTypePresident, agent.Type())
	}

	status := agent.GetStatus()
	if status.State != "idle" {
		t.Errorf("Expected initial state 'idle', got '%s'", status.State)
	}
}

func TestBaseAgentUpdateStatus(t *testing.T) {
	cfg := config.AgentConfig{}
	agent := NewBaseAgent("test-1", AgentTypeEmployee, cfg)

	agent.UpdateStatus("working", "task-123", "Processing task")

	status := agent.GetStatus()
	if status.State != "working" {
		t.Errorf("Expected state 'working', got '%s'", status.State)
	}
	if status.CurrentTask != "task-123" {
		t.Errorf("Expected current task 'task-123', got '%s'", status.CurrentTask)
	}
	if status.Message != "Processing task" {
		t.Errorf("Expected message 'Processing task', got '%s'", status.Message)
	}
}

func TestAgentPool(t *testing.T) {
	pool := NewAgentPool()
	cfg := config.AgentConfig{}

	// Create and register agents
	agent1 := NewBaseAgent("president-1", AgentTypePresident, cfg)
	agent2 := NewBaseAgent("employee-1", AgentTypeEmployee, cfg)
	agent3 := NewBaseAgent("employee-2", AgentTypeEmployee, cfg)

	if err := pool.Register(agent1); err != nil {
		t.Fatalf("Failed to register agent1: %v", err)
	}

	if err := pool.Register(agent2); err != nil {
		t.Fatalf("Failed to register agent2: %v", err)
	}

	if err := pool.Register(agent3); err != nil {
		t.Fatalf("Failed to register agent3: %v", err)
	}

	// Test duplicate registration
	if err := pool.Register(agent1); err == nil {
		t.Error("Expected error when registering duplicate agent")
	}

	// Test Get
	retrieved, err := pool.Get("president-1")
	if err != nil {
		t.Errorf("Failed to get agent: %v", err)
	}
	if retrieved.ID() != "president-1" {
		t.Errorf("Expected ID 'president-1', got '%s'", retrieved.ID())
	}

	// Test Get non-existent
	_, err = pool.Get("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent agent")
	}

	// Test GetByType
	employees := pool.GetByType(AgentTypeEmployee)
	if len(employees) != 2 {
		t.Errorf("Expected 2 employees, got %d", len(employees))
	}

	presidents := pool.GetByType(AgentTypePresident)
	if len(presidents) != 1 {
		t.Errorf("Expected 1 president, got %d", len(presidents))
	}

	// Test GetAvailable
	available, err := pool.GetAvailable(AgentTypeEmployee)
	if err != nil {
		t.Errorf("Failed to get available agent: %v", err)
	}
	if available == nil {
		t.Error("Expected to get an available agent")
	}

	// Test GetAvailable for non-existent type
	_, err = pool.GetAvailable(AgentTypeDepartmentManager)
	if err == nil {
		t.Error("Expected error when getting available agent of non-existent type")
	}

	// Test GetAllStatus
	statuses := pool.GetAllStatus()
	if len(statuses) != 3 {
		t.Errorf("Expected 3 statuses, got %d", len(statuses))
	}
}

func TestBaseAgentProcess(t *testing.T) {
	cfg := config.AgentConfig{}
	agent := NewBaseAgent("test-1", AgentTypeEmployee, cfg)

	// Test that Process returns an error for base implementation
	ctx := context.Background()
	_, err := agent.Process(ctx, "test input")
	if err == nil {
		t.Error("Expected error from base Process implementation")
	}
}
