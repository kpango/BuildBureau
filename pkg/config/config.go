package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type SystemConfig struct {
	TimeoutSeconds int    `yaml:"timeout_seconds"`
	MaxRetries     int    `yaml:"max_retries"`
	LogLevel       string `yaml:"log_level"`
}

type SlackConfig struct {
	Enabled   bool     `yaml:"enabled"`
	TokenEnv  string   `yaml:"token_env"`
	ChannelID string   `yaml:"channel_id"`
	Events    []string `yaml:"events"`
	// Resolved token
	Token string `yaml:"-"`
}

type ModelConfig struct {
	Name      string `yaml:"name"`
	APIKeyEnv string `yaml:"api_key_env"`
	// Resolved API Key
	APIKey string `yaml:"-"`
}

type AgentConfig struct {
	Role         string   `yaml:"role"`
	Count        int      `yaml:"count"`
	Model        string   `yaml:"model"`
	SystemPrompt string   `yaml:"system_prompt"`
	AllowedTools []string `yaml:"allowed_tools"`
}

type Config struct {
	System SystemConfig           `yaml:"system"`
	Slack  SlackConfig            `yaml:"slack"`
	Models map[string]ModelConfig `yaml:"models"`
	Agents map[string]AgentConfig `yaml:"agents"`
}

func Load(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := cfg.resolveSecrets(); err != nil {
		return nil, fmt.Errorf("failed to resolve secrets: %w", err)
	}

	return &cfg, nil
}

func (c *Config) resolveSecrets() error {
	// Resolve Slack Token
	if c.Slack.Enabled && c.Slack.TokenEnv != "" {
		c.Slack.Token = os.Getenv(c.Slack.TokenEnv)
		if c.Slack.Token == "" {
			// Warn or Error? For now, let's just log or ignore, maybe the user hasn't set it yet.
			// Ideally we might want to error if enabled is true.
			// fmt.Printf("Warning: Slack is enabled but environment variable %s is not set.\n", c.Slack.TokenEnv)
		}
	}

	// Resolve Model API Keys
	for k, model := range c.Models {
		if model.APIKeyEnv != "" {
			apiKey := os.Getenv(model.APIKeyEnv)
			model.APIKey = apiKey
			c.Models[k] = model // Update the map value
		}
	}

	return nil
}

func (c *Config) GetModelAPIKey(modelID string) string {
	if m, ok := c.Models[modelID]; ok {
		return m.APIKey
	}
	return ""
}
