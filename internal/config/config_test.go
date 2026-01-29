package config

import (
	"os"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.LLM.Provider != "gemini" {
		t.Errorf("Expected provider 'gemini', got %s", cfg.LLM.Provider)
	}

	if cfg.LLM.Model != "gemini-1.5-pro" {
		t.Errorf("Expected model 'gemini-1.5-pro', got %s", cfg.LLM.Model)
	}

	if !cfg.Agents.CEO.Enabled {
		t.Error("Expected CEO agent to be enabled")
	}

	if cfg.UI.RefreshRateMs != 100 {
		t.Errorf("Expected refresh rate 100, got %d", cfg.UI.RefreshRateMs)
	}
}

func TestLoad(t *testing.T) {
	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "config-test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	configContent := `
llm:
  provider: "test-provider"
  api_key: "test-key"
  model: "test-model"

agents:
  ceo:
    name: "Test CEO"
    enabled: true

slack:
  enabled: false

ui:
  show_agent_logs: true
  color_coding: true
  refresh_rate_ms: 200
`

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.LLM.Provider != "test-provider" {
		t.Errorf("Expected provider 'test-provider', got %s", cfg.LLM.Provider)
	}

	if cfg.LLM.APIKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got %s", cfg.LLM.APIKey)
	}

	if cfg.Agents.CEO.Name != "Test CEO" {
		t.Errorf("Expected CEO name 'Test CEO', got %s", cfg.Agents.CEO.Name)
	}

	if cfg.UI.RefreshRateMs != 200 {
		t.Errorf("Expected refresh rate 200, got %d", cfg.UI.RefreshRateMs)
	}
}

func TestLoadWithEnvVar(t *testing.T) {
	// Set environment variable
	os.Setenv("TEST_API_KEY", "env-api-key")
	defer os.Unsetenv("TEST_API_KEY")

	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "config-test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	configContent := `
llm:
  provider: "test-provider"
  api_key: "${TEST_API_KEY}"
  model: "test-model"
`

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.LLM.APIKey != "env-api-key" {
		t.Errorf("Expected API key 'env-api-key', got %s", cfg.LLM.APIKey)
	}
}
