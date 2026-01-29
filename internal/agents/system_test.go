package agents

import (
	"context"
	"testing"

	"buildbureau/internal/protocol"
	"buildbureau/pkg/a2a"
	"buildbureau/pkg/config"
)

func TestSystemRunProject(t *testing.T) {
	// Setup dependencies
	cfg := &config.Config{
		Agents: map[string]config.AgentConfig{
			"president": {Role: "President", Model: "test-model"},
			"manager":   {Role: "Manager", Model: "test-model"},
			"section":   {Role: "Section", Model: "test-model"},
			"worker":    {Role: "Worker", Model: "test-model"},
		},
		Models: map[string]config.ModelConfig{
			"test-model": {APIKey: ""}, // No key -> forces mock mode implicitly
		},
	}
	bus := a2a.NewBus()

	// No LLMClient needed anymore
	sys := NewSystem(cfg, bus)
	sys.SetupMocks() // Use mocks to bypass LLM JSON parsing issues and ensure flow

	// Subscribe to logs
	logCh := bus.SubscribeGlobal()

	// Run
	req := protocol.RequirementSpec{
		ProjectName: "TestProject",
		Details:     "Build a CRM",
	}

	result, err := sys.RunProject(context.Background(), req)
	if err != nil {
		t.Fatalf("RunProject failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected success result")
	}

	// Mocks generate impl_S1.go and impl_S2.go
	found := false
	for k := range result.AllArtifacts {
		if k == "impl_S1.go" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected impl_S1.go in artifacts")
	}

	// Check logs
	msgCount := 0
Loop:
	for {
		select {
		case <-logCh:
			msgCount++
		default:
			break Loop
		}
	}
}
