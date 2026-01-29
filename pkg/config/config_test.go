package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	yamlContent := `
system:
  timeout_seconds: 30
slack:
  enabled: true
  token_env: "TEST_SLACK_TOKEN"
models:
  test_model:
    name: "gpt-test"
    api_key_env: "TEST_API_KEY"
agents:
  president:
    role: "President"
`
	tmpfile, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(yamlContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Set environment variables
	os.Setenv("TEST_SLACK_TOKEN", "xoxb-1234")
	os.Setenv("TEST_API_KEY", "sk-test-key")
	defer os.Unsetenv("TEST_SLACK_TOKEN")
	defer os.Unsetenv("TEST_API_KEY")

	// Load config
	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify values
	if cfg.System.TimeoutSeconds != 30 {
		t.Errorf("Expected TimeoutSeconds 30, got %d", cfg.System.TimeoutSeconds)
	}

	if cfg.Slack.Token != "xoxb-1234" {
		t.Errorf("Expected Slack Token 'xoxb-1234', got '%s'", cfg.Slack.Token)
	}

	if key := cfg.GetModelAPIKey("test_model"); key != "sk-test-key" {
		t.Errorf("Expected API Key 'sk-test-key', got '%s'", key)
	}
}
