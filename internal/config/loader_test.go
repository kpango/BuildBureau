package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `
organization:
  layers:
    - name: President
      agent: ./agents/president.yaml
    - name: Engineer
      count: 2
      agent: ./agents/engineer.yaml

slack:
  enabled: false
  token: { env: SLACK_TOKEN }
  channels: ["#test"]
  notify_on: ["task_assigned"]

llms:
  default_model: gemini
  api_keys:
    gemini: { env: GEMINI_API_KEY }
`

	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, writeErr := tmpfile.WriteString(configContent); writeErr != nil {
		t.Fatal(writeErr)
	}
	_ = tmpfile.Close()

	// Set required environment variables
	_ = os.Setenv("GEMINI_API_KEY", "test-key")
	defer os.Unsetenv("GEMINI_API_KEY")

	// Test loading
	loader := NewLoader()
	cfg, err := loader.Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.LLMs.DefaultModel != "gemini" {
		t.Errorf("Expected default model 'gemini', got '%s'", cfg.LLMs.DefaultModel)
	}

	if len(cfg.Organization.Layers) != 2 {
		t.Errorf("Expected 2 layers, got %d", len(cfg.Organization.Layers))
	}

	if cfg.Slack != nil && cfg.Slack.Enabled {
		t.Error("Expected Slack to be disabled")
	}
}

func TestLoadConfigMissingEnvVar(t *testing.T) {
	configContent := `
organization:
  layers: []

llms:
  default_model: gemini
  api_keys:
    gemini: { env: MISSING_API_KEY }
`

	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, writeErr := tmpfile.WriteString(configContent); writeErr != nil {
		t.Fatal(writeErr)
	}
	_ = tmpfile.Close()

	// Make sure the env var is not set
	os.Unsetenv("MISSING_API_KEY")

	loader := NewLoader()
	_, err = loader.Load(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for missing environment variable, got nil")
	}
}

func TestGetEnvValue(t *testing.T) {
	_ = os.Setenv("TEST_VAR", "test-value") //nolint:gosec // G104: Test setup, error not critical
	defer os.Unsetenv("TEST_VAR")

	// Test getting environment variable value
	value := os.Getenv("TEST_VAR")
	if value != "test-value" {
		t.Errorf("Expected 'test-value', got '%s'", value)
	}
}
