package agent

import (
	"context"
	"testing"
	"time"

	"github.com/kpango/BuildBureau/pkg/types"
)

func TestBaseAgent(t *testing.T) {
	eventChan := make(chan types.AgentEvent, 10)
	agent := NewBaseAgent(types.RoleCEO, "Test CEO", eventChan)

	if agent.ID() == "" {
		t.Error("Agent ID should not be empty")
	}

	if agent.Role() != types.RoleCEO {
		t.Errorf("Expected role %s, got %s", types.RoleCEO, agent.Role())
	}

	if agent.Name() != "Test CEO" {
		t.Errorf("Expected name 'Test CEO', got %s", agent.Name())
	}
}

func TestCEOAgent(t *testing.T) {
	eventChan := make(chan types.AgentEvent, 100)
	ceo := NewCEOAgent("Test CEO", eventChan)

	task := types.Task{
		ID:          "test-task-1",
		Title:       "Test Task",
		Description: "Test task description",
		CreatedAt:   time.Now(),
		Status:      types.StatusPending,
	}

	ctx := context.Background()
	result, err := ceo.HandleTask(ctx, task)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Status != types.StatusCompleted {
		t.Errorf("Expected status %s, got %s", types.StatusCompleted, result.Status)
	}

	// Check that events were emitted
	if len(eventChan) == 0 {
		t.Error("Expected events to be emitted")
	}
}

func TestManagerAgent(t *testing.T) {
	eventChan := make(chan types.AgentEvent, 100)
	manager := NewManagerAgent("Test Manager", eventChan)

	task := types.Task{
		ID:          "test-task-2",
		Title:       "Test Task",
		Description: "Test task description",
		CreatedAt:   time.Now(),
		Status:      types.StatusPending,
	}

	ctx := context.Background()
	result, err := manager.HandleTask(ctx, task)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Status != types.StatusCompleted {
		t.Errorf("Expected status %s, got %s", types.StatusCompleted, result.Status)
	}
}

func TestLeadAgent(t *testing.T) {
	eventChan := make(chan types.AgentEvent, 100)
	lead := NewLeadAgent("Test Lead", eventChan)

	task := types.Task{
		ID:          "test-task-3",
		Title:       "Test Task",
		Description: "Test task description",
		CreatedAt:   time.Now(),
		Status:      types.StatusPending,
	}

	ctx := context.Background()
	result, err := lead.HandleTask(ctx, task)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Status != types.StatusCompleted {
		t.Errorf("Expected status %s, got %s", types.StatusCompleted, result.Status)
	}
}

func TestEmployeeAgent(t *testing.T) {
	eventChan := make(chan types.AgentEvent, 100)
	employee := NewEmployeeAgent("Test Employee", eventChan)

	task := types.Task{
		ID:          "test-task-4",
		Title:       "Test Task",
		Description: "Test task description",
		CreatedAt:   time.Now(),
		Status:      types.StatusPending,
	}

	ctx := context.Background()
	result, err := employee.HandleTask(ctx, task)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Status != types.StatusCompleted {
		t.Errorf("Expected status %s, got %s", types.StatusCompleted, result.Status)
	}
}

func TestSecretaryAgent(t *testing.T) {
	eventChan := make(chan types.AgentEvent, 100)
	ceo := NewCEOAgent("Test CEO", eventChan)
	secretary := NewSecretaryAgent("Test Secretary", ceo, eventChan)

	task := types.Task{
		ID:          "test-task-5",
		Title:       "Test Task",
		Description: "Test task description",
		CreatedAt:   time.Now(),
		Status:      types.StatusPending,
	}

	ctx := context.Background()
	result, err := secretary.HandleTask(ctx, task)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Status != types.StatusCompleted {
		t.Errorf("Expected status %s, got %s", types.StatusCompleted, result.Status)
	}
}

func TestAgentHierarchy(t *testing.T) {
	eventChan := make(chan types.AgentEvent, 1000)
	
	// Create CEO with secretary
	ceo := NewCEOAgent("Test CEO", eventChan)
	ceoSecretary := NewSecretaryAgent("CEO Secretary", ceo, eventChan)
	ceo.SetSecretary(ceoSecretary)

	// Create manager with secretary
	manager := NewManagerAgent("Test Manager", eventChan)
	managerSecretary := NewSecretaryAgent("Manager Secretary", manager, eventChan)
	manager.SetSecretary(managerSecretary)
	ceo.AddManagerAgent(manager)

	// Create lead with secretary
	lead := NewLeadAgent("Test Lead", eventChan)
	leadSecretary := NewSecretaryAgent("Lead Secretary", lead, eventChan)
	lead.SetSecretary(leadSecretary)
	manager.AddLeadAgent(lead)

	// Create employee
	employee := NewEmployeeAgent("Test Employee", eventChan)
	lead.AddEmployeeAgent(employee)

	// Test task delegation through hierarchy
	task := types.Task{
		ID:          "test-task-hierarchy",
		Title:       "Test Hierarchical Task",
		Description: "Test task going through the hierarchy",
		CreatedAt:   time.Now(),
		Status:      types.StatusPending,
	}

	ctx := context.Background()
	result, err := ceo.HandleTask(ctx, task)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Status != types.StatusCompleted {
		t.Errorf("Expected status %s, got %s", types.StatusCompleted, result.Status)
	}

	// Wait a bit for goroutines to complete
	time.Sleep(2 * time.Second)

	// Check that multiple events were emitted (indicating hierarchy activation)
	eventCount := len(eventChan)
	if eventCount < 5 {
		t.Errorf("Expected at least 5 events, got %d", eventCount)
	}
}
