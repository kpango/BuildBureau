package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"buildbureau/pkg/adk"
	"buildbureau/pkg/agent"
	"buildbureau/pkg/config"
	"buildbureau/pkg/protocol"
	"buildbureau/pkg/tui"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded config for organization with %d layers", len(cfg.Organization.Layers))

	// 2. Initialize ADK
	adkClient, err := adk.NewClient(context.Background(), adk.Config{
		ProjectID: "buildbureau",
		Location:  "us-central1",
	})
	if err != nil {
		log.Fatalf("Failed to create ADK client: %v", err)
	}

	// Initialize Slack Service
	slackSvc := &agent.SlackService{
		Enabled:  cfg.Slack.Enabled,
		Token:    string(cfg.Slack.Token),
		Channels: cfg.Slack.Channels,
		NotifyOn: cfg.Slack.NotifyOn,
	}

	// 3. Bootstrap Agents
	agents := make(map[string]*agent.Agent)
	basePort := 50051

	// Helper to create agent
	createAgent := func(name, role, configFile string) *agent.Agent {
		agentCfg, err := config.LoadAgentConfig(configFile)
		sysPrompt := ""
		if err == nil {
			sysPrompt = agentCfg.SystemPrompt
		}
		a := agent.NewAgent(name, role, basePort, sysPrompt, adkClient, slackSvc)
		basePort++
		agents[name] = a
		go func() {
			if err := a.Start(); err != nil {
				log.Fatalf("Agent %s failed: %v", name, err)
			}
		}()
		return a
	}

	// Instantiate hierarchy (Simplified tree)
	// President
	pres := createAgent("President", "President", "./agents/president.yaml")
	presSec := createAgent("PresidentSecretary", "Secretary", "./agents/secretary.yaml")

	// Director
	dir := createAgent("Director", "Director", "./agents/director.yaml")
	dirSec := createAgent("DirectorSecretary", "Secretary", "./agents/secretary.yaml")

	// Manager
	mgr := createAgent("Manager", "Manager", "./agents/manager.yaml")
	mgrSec := createAgent("ManagerSecretary", "Secretary", "./agents/secretary.yaml")

	// Engineer
	eng := createAgent("Engineer", "Engineer", "./agents/engineer.yaml")

	// Start Remote/Sub Agents
	for _, sub := range cfg.SubAgents {
		u, err := url.Parse(sub.Remote.Endpoint)
		if err != nil {
			log.Printf("Invalid endpoint for %s: %v", sub.Name, err)
			continue
		}
		portStr := u.Port()
		if portStr == "" {
			portStr = "8080"
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Printf("Invalid port for %s: %v", sub.Name, err)
			continue
		}

		a := agent.NewAgent(sub.Name, "Worker", port, fmt.Sprintf("You are a specialized worker: %v", sub.Remote.Capabilities), adkClient, slackSvc)
		go func() {
			if err := a.Start(); err != nil {
				log.Printf("Worker %s failed: %v", sub.Name, err)
			}
		}()
		agents[sub.Name] = a
		log.Printf("Started SubAgent %s on port %d", sub.Name, port)
	}

	// Allow servers to start
	time.Sleep(1 * time.Second)

	// 4. Wire them up
	connect := func(superior, subordinate *agent.Agent) {
		addr := fmt.Sprintf("localhost:%d", subordinate.Port)
		if err := superior.ConnectToSubordinate(subordinate.Name, addr); err != nil {
			log.Fatalf("Failed to connect %s to %s: %v", superior.Name, subordinate.Name, err)
		}

		supAddr := fmt.Sprintf("localhost:%d", superior.Port)
		if err := subordinate.ConnectToSuperior(supAddr); err != nil {
			log.Fatalf("Failed to connect %s to superior %s: %v", subordinate.Name, superior.Name, err)
		}
	}

	// President -> President Secretary
	connect(pres, presSec)

	// President Secretary -> Director Secretary (This is a peer/delegation link, but we'll treat as subordinate for flow)
	connect(presSec, dirSec)

	// Director Secretary -> Director
	connect(dirSec, dir)

	// Director -> Manager Secretary
	connect(dir, mgrSec)

	// Manager Secretary -> Manager
	connect(mgrSec, mgr)

	// Manager -> Engineer
	connect(mgr, eng)

	log.Println("All agents started and connected.")

	// 5. Connect TUI to President
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", pres.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype("json")),
	)
	if err != nil {
		log.Fatalf("Failed to dial President: %v", err)
	}
	presClient := protocol.NewAgentServiceClient(conn)

	// 6. Run TUI
	if err := tui.Start(presClient); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
