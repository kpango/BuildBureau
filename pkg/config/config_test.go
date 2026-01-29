package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	yamlContent := `
organization:
  layers:
    - name: President
      agent: ./agents/president.yaml
slack:
  enabled: true
  token: { env: TEST_SLACK_TOKEN }
llms:
  default_model: gemini
  api_keys:
    gemini: { env: TEST_GEMINI_KEY }
`
	tmpFile, err := os.CreateTemp("", "config_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	os.Setenv("TEST_SLACK_TOKEN", "my-slack-token")
	os.Setenv("TEST_GEMINI_KEY", "my-gemini-key")

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Slack.Token != "my-slack-token" {
		t.Errorf("Expected slack token 'my-slack-token', got '%s'", cfg.Slack.Token)
	}
	if cfg.LLMs.APIKeys["gemini"] != "my-gemini-key" {
		t.Errorf("Expected gemini key 'my-gemini-key', got '%s'", cfg.LLMs.APIKeys["gemini"])
	}
	if len(cfg.Organization.Layers) != 1 {
		t.Errorf("Expected 1 layer, got %d", len(cfg.Organization.Layers))
	}
}
