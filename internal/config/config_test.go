package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")

	configContent := `
agents:
  president:
    count: 1
    model: "test-model"
    instruction: "Test instruction"
    allowTools: true
    tools: ["tool1"]
    timeout: 60
    retryCount: 3
  president_secretary:
    count: 1
    model: "test-model"
    instruction: "Test instruction"
    allowTools: false
    timeout: 30
    retryCount: 2
  department_manager:
    count: 1
    model: "test-model"
    instruction: "Test instruction"
    allowTools: true
    timeout: 60
    retryCount: 3
  department_secretary:
    count: 1
    model: "test-model"
    instruction: "Test instruction"
    allowTools: false
    timeout: 30
    retryCount: 2
  section_manager:
    count: 2
    model: "test-model"
    instruction: "Test instruction"
    allowTools: true
    timeout: 60
    retryCount: 3
  section_secretary:
    count: 2
    model: "test-model"
    instruction: "Test instruction"
    allowTools: false
    timeout: 30
    retryCount: 2
  employee:
    count: 4
    model: "test-model"
    instruction: "Test instruction"
    allowTools: true
    timeout: 90
    retryCount: 3

llm:
  provider: "test-provider"
  apiEndpoint: "https://api.test.com"
  defaultModel: "test-model"
  maxTokens: 4096
  temperature: 0.7
  topP: 0.9

grpc:
  port: 50051
  maxMessageSize: 1048576
  timeout: 60
  enableReflection: true

slack:
  enabled: false
  token: ""
  channelID: ""
  retryCount: 3
  timeout: 10
  notifications:
    projectStart:
      enabled: true
      message: "Project started"
    taskComplete:
      enabled: true
      message: "Task complete"
    error:
      enabled: true
      message: "Error occurred"
    projectComplete:
      enabled: true
      message: "Project complete"

ui:
  enableTUI: true
  refreshRate: 100
  theme: "default"
  showProgress: true
  logLevel: "info"

system:
  workDir: "./work"
  logDir: "./logs"
  cacheDir: "./cache"
  maxConcurrentTasks: 5
  globalTimeout: 3600
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Test loading the config
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Validate loaded config
	if cfg.Agents.President.Count != 1 {
		t.Errorf("Expected president count 1, got %d", cfg.Agents.President.Count)
	}

	if cfg.LLM.Provider != "test-provider" {
		t.Errorf("Expected provider 'test-provider', got '%s'", cfg.LLM.Provider)
	}

	if cfg.GRPC.Port != 50051 {
		t.Errorf("Expected gRPC port 50051, got %d", cfg.GRPC.Port)
	}

	if cfg.Slack.Enabled {
		t.Error("Expected Slack to be disabled")
	}

	if !cfg.UI.EnableTUI {
		t.Error("Expected TUI to be enabled")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "valid config",
			config: Config{
				Agents: AgentsConfig{
					President: AgentConfig{Count: 1},
					DepartmentManager: AgentConfig{Count: 1},
				},
				LLM: LLMConfig{
					Provider: "test",
				},
				GRPC: GRPCConfig{
					Port: 50051,
				},
				Slack: SlackConfig{
					Enabled: false,
				},
			},
			expectError: false,
		},
		{
			name: "missing president",
			config: Config{
				Agents: AgentsConfig{
					President: AgentConfig{Count: 0},
					DepartmentManager: AgentConfig{Count: 1},
				},
				LLM: LLMConfig{
					Provider: "test",
				},
				GRPC: GRPCConfig{
					Port: 50051,
				},
			},
			expectError: true,
		},
		{
			name: "missing LLM provider",
			config: Config{
				Agents: AgentsConfig{
					President: AgentConfig{Count: 1},
					DepartmentManager: AgentConfig{Count: 1},
				},
				LLM: LLMConfig{
					Provider: "",
				},
				GRPC: GRPCConfig{
					Port: 50051,
				},
			},
			expectError: true,
		},
		{
			name: "invalid gRPC port",
			config: Config{
				Agents: AgentsConfig{
					President: AgentConfig{Count: 1},
					DepartmentManager: AgentConfig{Count: 1},
				},
				LLM: LLMConfig{
					Provider: "test",
				},
				GRPC: GRPCConfig{
					Port: -1,
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestEnvironmentVariableExpansion(t *testing.T) {
	// Set test environment variables
	os.Setenv("TEST_TOKEN", "test-token-value")
	os.Setenv("TEST_CHANNEL", "test-channel-id")
	defer func() {
		os.Unsetenv("TEST_TOKEN")
		os.Unsetenv("TEST_CHANNEL")
	}()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")

	configContent := `
agents:
  president:
    count: 1
    model: "test-model"
    instruction: "Test"
    timeout: 60
    retryCount: 3
  president_secretary:
    count: 1
    model: "test-model"
    instruction: "Test"
    timeout: 30
    retryCount: 2
  department_manager:
    count: 1
    model: "test-model"
    instruction: "Test"
    timeout: 60
    retryCount: 3
  department_secretary:
    count: 1
    model: "test-model"
    instruction: "Test"
    timeout: 30
    retryCount: 2
  section_manager:
    count: 1
    model: "test-model"
    instruction: "Test"
    timeout: 60
    retryCount: 3
  section_secretary:
    count: 1
    model: "test-model"
    instruction: "Test"
    timeout: 30
    retryCount: 2
  employee:
    count: 1
    model: "test-model"
    instruction: "Test"
    timeout: 60
    retryCount: 3

llm:
  provider: "test"
  apiEndpoint: "https://test.com"
  defaultModel: "test"
  maxTokens: 1000
  temperature: 0.5
  topP: 0.9

grpc:
  port: 50051
  maxMessageSize: 1000000
  timeout: 60
  enableReflection: false

slack:
  enabled: true
  token: "${TEST_TOKEN}"
  channelID: "${TEST_CHANNEL}"
  retryCount: 3
  timeout: 10
  notifications:
    projectStart:
      enabled: true
      message: "Test"
    taskComplete:
      enabled: true
      message: "Test"
    error:
      enabled: true
      message: "Test"
    projectComplete:
      enabled: true
      message: "Test"

ui:
  enableTUI: false
  refreshRate: 100
  theme: "default"
  showProgress: true
  logLevel: "info"

system:
  workDir: "./work"
  logDir: "./logs"
  cacheDir: "./cache"
  maxConcurrentTasks: 5
  globalTimeout: 3600
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Slack.Token != "test-token-value" {
		t.Errorf("Expected token 'test-token-value', got '%s'", cfg.Slack.Token)
	}

	if cfg.Slack.ChannelID != "test-channel-id" {
		t.Errorf("Expected channel ID 'test-channel-id', got '%s'", cfg.Slack.ChannelID)
	}
}
