package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"buildbureau/internal/agents"
	"buildbureau/internal/ui"
	"buildbureau/pkg/a2a"
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

	// Check if we have keys to decide on Mock vs Real
	hasKeys := false
	for _, m := range cfg.Models {
		if m.APIKey != "" {
			hasKeys = true
			break
		}
	}

	// 3. Initialize System (Agents)
	sys := agents.NewSystem(cfg, bus)

	if !hasKeys {
		sys.SetupMocks()
		fmt.Println("Agents configured with Mock Implementations (No API keys found)")
	} else {
		fmt.Println("Agents configured with Real ADK (Keys found)")
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
