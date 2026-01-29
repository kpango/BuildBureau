package adk

import (
	"context"
	"testing"
	"time"

	"buildbureau/pkg/a2a"
	"buildbureau/pkg/config"
)

type TestReq struct {
	Input string
}

type TestResp struct {
	Output string
}

func TestAgentProcess(t *testing.T) {
	bus := a2a.NewBus()

	cfg := config.AgentConfig{
		Role: "Tester",
		Model: "gpt4",
		SystemPrompt: "You are a tester.",
	}

	// Pass empty API key to force no-runner mode
	agent := NewAgent[TestReq, TestResp]("test-agent", cfg, bus, "", "")

	// Inject MockImpl
	agent.MockImpl = func(ctx context.Context, req TestReq) (TestResp, error) {
		return TestResp{Output: "Processed: " + req.Input}, nil
	}

	// Subscribe to logs
	logCh := bus.Subscribe("LOG")

	// Run Process
	go func() {
		resp, err := agent.Process(context.Background(), TestReq{Input: "Hello"})
		if err != nil {
			t.Errorf("Process failed: %v", err)
		}
		if resp.Output != "Processed: Hello" {
			t.Errorf("Unexpected output: %s", resp.Output)
		}
	}()

	// Check logs
	select {
	case msg := <-logCh:
		if msg.Type != "START" {
			t.Errorf("Expected START message, got %s", msg.Type)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for START log")
	}

	select {
	case msg := <-logCh:
		if msg.Type != "COMPLETE" {
			t.Errorf("Expected COMPLETE message, got %s", msg.Type)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for COMPLETE log")
	}
}
