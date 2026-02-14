package agent

import (
	"context"
	"testing"

	"github.com/kpango/BuildBureau/pkg/types"
)

func TestGenericAgent_Creation(t *testing.T) {
	config := &types.AgentConfig{
		Name:         "TestAgent",
		Role:         "TestRole",
		Description:  "A test agent",
		Model:        "gemini",
		SystemPrompt: "You are a test agent.",
		Capabilities: []string{"testing", "validation"},
	}

	agent := NewGenericAgent("test-1", types.RoleEngineer, config, nil)

	if agent == nil {
		t.Fatal("Expected non-nil agent")
	}

	if agent.GetID() != "test-1" {
		t.Errorf("Expected ID 'test-1', got '%s'", agent.GetID())
	}

	if agent.GetRole() != types.RoleEngineer {
		t.Errorf("Expected role 'Engineer', got '%s'", agent.GetRole())
	}
}

func TestGenericAgent_Hierarchy(t *testing.T) {
	config := &types.AgentConfig{
		Name:         "TestAgent",
		Role:         "Manager",
		Description:  "A manager agent",
		SystemPrompt: "You are a manager.",
	}

	manager := NewGenericAgent("manager-1", types.RoleManager, config, nil)
	engineer1 := NewGenericAgent("engineer-1", types.RoleEngineer, config, nil)
	engineer2 := NewGenericAgent("engineer-2", types.RoleEngineer, config, nil)

	manager.AddSubordinate(engineer1)
	manager.AddSubordinate(engineer2)

	subordinates := manager.GetSubordinates()
	if len(subordinates) != 2 {
		t.Errorf("Expected 2 subordinates, got %d", len(subordinates))
	}

	engineer1.SetParent(manager)
	// Note: We can't directly test GetParent as it's not exposed, but this verifies the API works
}

func TestGenericAgent_ProcessTaskWithoutLLM(t *testing.T) {
	config := &types.AgentConfig{
		Name:         "TestAgent",
		Role:         "Engineer",
		Description:  "A test engineer",
		SystemPrompt: "You are a test engineer.",
		Capabilities: []string{"coding", "testing"},
	}

	agent := NewGenericAgent("engineer-1", types.RoleEngineer, config, nil)

	ctx := context.Background()
	if err := agent.Start(ctx); err != nil {
		t.Fatalf("Failed to start agent: %v", err)
	}

	task := &types.Task{
		ID:          "task-1",
		Title:       "Test Task",
		Description: "This is a test task",
		Content:     "Additional context",
		Priority:    1,
	}

	response, err := agent.ProcessTask(ctx, task)
	if err != nil {
		t.Fatalf("ProcessTask failed: %v", err)
	}

	if response.Status != types.StatusCompleted {
		t.Errorf("Expected status 'completed', got '%s'", response.Status)
	}

	if response.Result == "" {
		t.Error("Expected non-empty result")
	}

	if response.TaskID != task.ID {
		t.Errorf("Expected task ID '%s', got '%s'", task.ID, response.TaskID)
	}
}

func TestGenericAgent_Delegation(t *testing.T) {
	config := &types.AgentConfig{
		Name:         "Manager",
		Role:         "Manager",
		Description:  "A manager who delegates",
		SystemPrompt: "You are a manager who delegates tasks.",
	}

	manager := NewGenericAgent("manager-1", types.RoleManager, config, nil)
	engineer := NewGenericAgent("engineer-1", types.RoleEngineer, config, nil)

	manager.AddSubordinate(engineer)

	ctx := context.Background()
	if err := manager.Start(ctx); err != nil {
		t.Fatalf("Failed to start manager: %v", err)
	}
	if err := engineer.Start(ctx); err != nil {
		t.Fatalf("Failed to start engineer: %v", err)
	}

	task := &types.Task{
		ID:          "task-1",
		Title:       "Build Feature",
		Description: "Build a new feature",
		Priority:    1,
	}

	response, err := manager.ProcessTask(ctx, task)
	if err != nil {
		t.Fatalf("ProcessTask failed: %v", err)
	}

	if response.Status != types.StatusCompleted {
		t.Errorf("Expected status 'completed', got '%s'", response.Status)
	}

	// Check that the manager processed and delegated
	if response.Result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestGenericAgent_StartStop(t *testing.T) {
	config := &types.AgentConfig{
		Name: "TestAgent",
		Role: "Engineer",
	}

	agent := NewGenericAgent("test-1", types.RoleEngineer, config, nil)
	ctx := context.Background()

	// Test starting
	if err := agent.Start(ctx); err != nil {
		t.Fatalf("Failed to start agent: %v", err)
	}

	if !agent.IsRunning() {
		t.Error("Expected agent to be running")
	}

	// Test stopping
	if err := agent.Stop(ctx); err != nil {
		t.Fatalf("Failed to stop agent: %v", err)
	}

	if agent.IsRunning() {
		t.Error("Expected agent to be stopped")
	}
}

func TestGenericAgent_Stats(t *testing.T) {
	config := &types.AgentConfig{
		Name: "TestAgent",
		Role: "Engineer",
	}

	agent := NewGenericAgent("test-1", types.RoleEngineer, config, nil)

	// Initial stats
	active, completed := agent.GetStats()
	if active != 0 || completed != 0 {
		t.Errorf("Expected initial stats (0,0), got (%d,%d)", active, completed)
	}

	// Increment active
	agent.IncrementActiveTasks()
	active, completed = agent.GetStats()
	if active != 1 || completed != 0 {
		t.Errorf("Expected stats (1,0), got (%d,%d)", active, completed)
	}

	// Decrement active (increments completed)
	agent.DecrementActiveTasks()
	active, completed = agent.GetStats()
	if active != 0 || completed != 1 {
		t.Errorf("Expected stats (0,1), got (%d,%d)", active, completed)
	}
}
