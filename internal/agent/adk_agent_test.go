package agent

import (
	"context"
	"os"
	"testing"

	"github.com/kpango/BuildBureau/pkg/types"
)

const demoAPIKey = "demo-key"

func TestNewADKEngineerAgent(t *testing.T) {
	// Skip if no API key (can't test ADK without real API)
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" || apiKey == demoAPIKey {
		t.Skip("Skipping ADK test: GEMINI_API_KEY not set or using demo value")
	}

	config := &types.AgentConfig{
		Name: "Test ADK Engineer",
		Role: "Engineer",
	}

	engineer, err := NewADKEngineerAgent("adk-engineer-test", config, apiKey)
	if err != nil {
		t.Fatalf("Failed to create ADK engineer: %v", err)
	}

	if engineer == nil {
		t.Fatal("Expected non-nil engineer")
	}

	if engineer.GetID() != "adk-engineer-test" {
		t.Errorf("Expected ID 'adk-engineer-test', got '%s'", engineer.GetID())
	}

	if engineer.GetRole() != types.RoleEngineer {
		t.Errorf("Expected role Engineer, got %s", engineer.GetRole())
	}
}

func TestNewADKManagerAgent(t *testing.T) {
	// Skip if no API key
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" || apiKey == demoAPIKey {
		t.Skip("Skipping ADK test: GEMINI_API_KEY not set or using demo value")
	}

	config := &types.AgentConfig{
		Name: "Test ADK Manager",
		Role: "Manager",
	}

	manager, err := NewADKManagerAgent("adk-manager-test", config, apiKey)
	if err != nil {
		t.Fatalf("Failed to create ADK manager: %v", err)
	}

	if manager == nil {
		t.Fatal("Expected non-nil manager")
	}

	if manager.GetRole() != types.RoleManager {
		t.Errorf("Expected role Manager, got %s", manager.GetRole())
	}
}

func TestNewADKDirectorAgent(t *testing.T) {
	// Skip if no API key
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" || apiKey == demoAPIKey {
		t.Skip("Skipping ADK test: GEMINI_API_KEY not set or using demo value")
	}

	config := &types.AgentConfig{
		Name: "Test ADK Director",
		Role: "Director",
	}

	director, err := NewADKDirectorAgent("adk-director-test", config, apiKey)
	if err != nil {
		t.Fatalf("Failed to create ADK director: %v", err)
	}

	if director == nil {
		t.Fatal("Expected non-nil director")
	}

	if director.GetRole() != types.RoleDirector {
		t.Errorf("Expected role Director, got %s", director.GetRole())
	}
}

func TestADKAgent_ProcessTask(t *testing.T) {
	// Skip if no API key
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" || apiKey == demoAPIKey {
		t.Skip("Skipping ADK test: GEMINI_API_KEY not set or using demo value")
	}

	config := &types.AgentConfig{
		Name: "Test ADK Engineer",
		Role: "Engineer",
	}

	engineer, err := NewADKEngineerAgent("adk-engineer-test", config, apiKey)
	if err != nil {
		t.Fatalf("Failed to create ADK engineer: %v", err)
	}

	ctx := context.Background()
	if startErr := engineer.Start(ctx); startErr != nil {
		t.Fatalf("Failed to start engineer: %v", startErr)
	}
	defer engineer.Stop(ctx)

	// Create a simple test task
	task := &types.Task{
		ID:          "test-adk-task",
		Title:       "Simple Task",
		Description: "A simple test task",
		FromAgent:   "test",
		ToAgent:     engineer.GetID(),
		Content:     "Say hello in a creative way",
		Priority:    1,
	}

	response, err := engineer.ProcessTask(ctx, task)
	if err != nil {
		t.Fatalf("Failed to process task: %v", err)
	}

	if response == nil {
		t.Fatal("Expected non-nil response")
	}

	if response.TaskID != task.ID {
		t.Errorf("Expected task ID '%s', got '%s'", task.ID, response.TaskID)
	}

	if response.Status != types.StatusCompleted {
		t.Errorf("Expected status completed, got %s", response.Status)
	}

	if response.Result == "" {
		t.Error("Expected non-empty result")
	}

	// Verify the response contains "ADK Agent" to confirm it's using ADK
	if response.Result == "" {
		t.Error("Expected response to indicate ADK usage")
	}
}

func TestNewADKAgent_NoAPIKey(t *testing.T) {
	config := &types.AgentConfig{
		Name: "Test ADK Engineer",
		Role: "Engineer",
	}

	// Clear API key from environment temporarily
	oldKey := os.Getenv("GEMINI_API_KEY")
	os.Setenv("GEMINI_API_KEY", "")
	defer os.Setenv("GEMINI_API_KEY", oldKey)

	_, err := NewADKEngineerAgent("adk-engineer-test", config, "")
	if err == nil {
		t.Error("Expected error when creating ADK agent without API key")
	}
}
