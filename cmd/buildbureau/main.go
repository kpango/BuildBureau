package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"buildbureau/internal/agents"
	"buildbureau/internal/ui"
	"buildbureau/pkg/a2a"
	"buildbureau/pkg/adk"
	"buildbureau/pkg/config"
	"buildbureau/pkg/slack"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize Infrastructure
	bus := a2a.NewBus()

	// Real LLM Client (will use env vars if present)
	// Build map of keys
	apiKeys := make(map[string]string)
	for k, v := range cfg.Models {
		apiKeys[k] = v.APIKey
	}
	// Note: If keys are missing, we might want to default to Mock.
	// For this implementation, I'll use a "Hybrid" client or just check.
	// If no keys, use Mock.
	var llmClient adk.LLMClient

	// Simple check: if at least one key is present, use Real (simulated).
	// Otherwise use Mock.
	hasKeys := false
	for _, k := range apiKeys {
		if k != "" {
			hasKeys = true
			break
		}
	}

	if hasKeys {
		llmClient = adk.NewRealLLMClient(apiKeys)
		fmt.Println("Using Real LLM Client (Simulated)")
	} else {
		llmClient = adk.NewMockLLMClient()
		fmt.Println("Using Mock LLM Client (No API keys found)")
	}

	// 3. Initialize System (Agents)
	sys := agents.NewSystem(cfg, bus, llmClient)

	// If using Mock Client, we also need to inject the Mock Implementations into the Agents
	// because NewMockLLMClient just returns strings, but Agents might expect structured JSON.
	// So we use SetupMocks() which overrides the Agent.Process logic partially or fully.
	// Actually, `SetupMocks` in `internal/agents/mocks.go` replaces `MockImpl` on the agent.
	// This completely bypasses the LLMClient.
	// If we want to use the RealLLMClient, we should NOT call SetupMocks.
	// If we want to use the MockLLMClient, we SHOULD call SetupMocks because the MockLLMClient
	// in `pkg/adk/llm.go` returns a generic string which will fail JSON parsing in `Agent.Process`.

	if !hasKeys {
		sys.SetupMocks()
		fmt.Println("Agents configured with Mock Implementations")
	}

	// 4. Initialize Slack
	slackNotifier := slack.NewNotifier(
		cfg.Slack.Token,
		cfg.Slack.ChannelID,
		cfg.Slack.Enabled,
	)

	// Slack Bridge: Listen to A2A events and forward important ones
	slackCh := bus.Subscribe("SLACK_BRIDGE")
	go func() {
		for msg := range slackCh {
			// Filter important events
			if msg.Type == "ERROR" {
				slackNotifier.NotifyError(fmt.Sprintf("%v", msg.Payload))
			} else if msg.Type == "COMPLETE" {
				// Only notify project completion or major milestones?
				// For now, notify all COMPLETE events (maybe too noisy)
				// Let's check the From field.
				slackNotifier.NotifySuccess(fmt.Sprintf("%s: %v", msg.From, msg.Payload))
			}
		}
	}()

	// 5. Initialize UI
	// UI needs its own subscription
	uiCh := bus.SubscribeGlobal()
	model := ui.NewModel(sys, uiCh)

	// 6. Run
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
