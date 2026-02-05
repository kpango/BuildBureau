package agent

import (
	"context"
	"testing"

	"github.com/kpango/BuildBureau/pkg/types"
)

func TestBaseAgent(t *testing.T) {
	config := &types.AgentConfig{
		Name: "TestAgent",
		Role: "test",
	}

	agent := NewBaseAgent("test-1", types.RoleEngineer, config)

	if agent.GetID() != "test-1" {
		t.Errorf("Expected ID 'test-1', got '%s'", agent.GetID())
	}

	if agent.GetRole() != types.RoleEngineer {
		t.Errorf("Expected role 'Engineer', got '%s'", agent.GetRole())
	}

	if agent.IsRunning() {
		t.Error("Expected agent to not be running initially")
	}

	// Test start
	ctx := context.Background()
	if err := agent.Start(ctx); err != nil {
		t.Errorf("Failed to start agent: %v", err)
	}

	if !agent.IsRunning() {
		t.Error("Expected agent to be running after start")
	}

	// Test double start
	if err := agent.Start(ctx); err == nil {
		t.Error("Expected error when starting already running agent")
	}

	// Test stop
	if err := agent.Stop(ctx); err != nil {
		t.Errorf("Failed to stop agent: %v", err)
	}

	if agent.IsRunning() {
		t.Error("Expected agent to not be running after stop")
	}
}

func TestEngineerAgent(t *testing.T) {
	config := &types.AgentConfig{
		Name: "TestEngineer",
		Role: "Engineer",
	}

	engineer := NewEngineerAgent("engineer-1", config, nil) // nil LLM manager for test

	ctx := context.Background()
	if err := engineer.Start(ctx); err != nil {
		t.Fatalf("Failed to start engineer: %v", err)
	}
	defer engineer.Stop(ctx)

	task := &types.Task{
		ID:          "test-task-1",
		Title:       "Test Task",
		Description: "Implement test feature",
		FromAgent:   "manager-1",
		ToAgent:     engineer.GetID(),
		Content:     "Write a function that adds two numbers",
		Priority:    1,
	}

	response, err := engineer.ProcessTask(ctx, task)
	if err != nil {
		t.Fatalf("Failed to process task: %v", err)
	}

	if response.TaskID != task.ID {
		t.Errorf("Expected task ID '%s', got '%s'", task.ID, response.TaskID)
	}

	if response.Status != types.StatusCompleted {
		t.Errorf("Expected status 'completed', got '%s'", response.Status)
	}

	if response.Result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestPresidentWithSecretary(t *testing.T) {
	presidentCfg := &types.AgentConfig{
		Name: "TestPresident",
		Role: "President",
	}

	secretaryCfg := &types.AgentConfig{
		Name: "TestSecretary",
		Role: "Secretary",
	}

	president := NewPresidentAgent("president-1", presidentCfg)
	secretary := NewSecretaryAgent("secretary-1", secretaryCfg)

	president.SetSecretary(secretary)

	ctx := context.Background()
	president.Start(ctx)
	secretary.Start(ctx)
	defer president.Stop(ctx)
	defer secretary.Stop(ctx)

	task := &types.Task{
		ID:          "client-task-1",
		Title:       "Client Request",
		Description: "Build a web application",
		FromAgent:   "client",
		ToAgent:     president.GetID(),
		Content:     "Build a web application",
		Priority:    1,
	}

	response, err := president.ProcessTask(ctx, task)
	if err != nil {
		t.Fatalf("Failed to process task: %v", err)
	}

	if response.Status != types.StatusCompleted {
		t.Errorf("Expected status 'completed', got '%s'", response.Status)
	}

	// The result should mention both president and secretary
	if response.Result == "" {
		t.Error("Expected non-empty result")
	}
}
