package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kpango/BuildBureau/internal/agent"
	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/internal/slack"
	"github.com/kpango/BuildBureau/internal/ui"
)

var (
	configPath = flag.String("config", "configs/config.yaml", "Path to configuration file")
	version    = "1.0.0"
)

func main() {
	flag.Parse()

	// Setup logging
	logger := log.New(os.Stdout, "[BuildBureau] ", log.LstdFlags)
	logger.Printf("Starting BuildBureau v%s\n", version)

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v\n", err)
	}
	logger.Println("Configuration loaded successfully")

	// Setup file logging if configured
	if cfg.System.Logging.EnableFileLogging {
		// Create logs directory if it doesn't exist
		os.MkdirAll("logs", 0755)
		
		logFile, err := os.OpenFile(cfg.System.Logging.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			logger.Printf("Warning: Failed to open log file: %v\n", err)
		} else {
			defer logFile.Close()
			logger.SetOutput(logFile)
		}
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Println("Shutdown signal received, cleaning up...")
		cancel()
	}()

	// Initialize Slack notifier
	notifier, err := slack.NewNotifier(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize Slack notifier: %v\n", err)
	}
	logger.Println("Slack notifier initialized")

	// Initialize agent system
	agentSystem, err := agent.NewAgentSystem(ctx, cfg, notifier, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize agent system: %v\n", err)
	}
	defer agentSystem.Close()
	logger.Println("Agent system initialized")

	// Create knowledge base directory
	if cfg.System.KnowledgeBase.Type == "file" {
		os.MkdirAll(cfg.System.KnowledgeBase.Path, 0755)
	}

	// Check if UI is enabled
	if !cfg.System.UI.Enabled {
		logger.Println("UI is disabled, running in headless mode")
		// In headless mode, you could implement a different interface
		// For now, we'll just wait for signals
		<-ctx.Done()
		return
	}

	// Initialize UI
	logger.Println("Starting Terminal UI...")

	// Create a handler for user input
	handleInput := func(input string) error {
		logger.Printf("Processing user input: %s\n", input)
		
		// Process the request through the agent system
		response, err := agentSystem.ProcessClientRequest(ctx, "client-1", input)
		if err != nil {
			logger.Printf("Error processing request: %v\n", err)
			return err
		}
		
		logger.Printf("Request processed successfully: %s\n", response)
		return nil
	}

	// Create UI model
	model := ui.NewModel(cfg, handleInput)

	// Start the Bubble Tea program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		logger.Fatalf("Error running UI: %v\n", err)
	}

	logger.Println("BuildBureau shutdown complete")
}
