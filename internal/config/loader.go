package config

import (
	"fmt"
	"os"

	"github.com/kpango/BuildBureau/pkg/types"
	"gopkg.in/yaml.v3"
)

// Loader handles loading and parsing configuration files.
type Loader struct{}

// NewLoader creates a new configuration loader.
func NewLoader() *Loader {
	return &Loader{}
}

// Load reads and parses a YAML configuration file.
func (l *Loader) Load(path string) (*types.Config, error) {
	data, err := os.ReadFile(path) //nolint:gosec // G304: Config file path is from trusted source
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config types.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Resolve environment variables
	if err := l.resolveEnvVars(&config); err != nil {
		return nil, fmt.Errorf("failed to resolve environment variables: %w", err)
	}

	return &config, nil
}

// LoadAgentConfig loads an individual agent configuration file.
func (l *Loader) LoadAgentConfig(path string) (*types.AgentConfig, error) {
	data, err := os.ReadFile(path) //nolint:gosec // G304: Agent config file path is from trusted source
	if err != nil {
		return nil, fmt.Errorf("failed to read agent config file: %w", err)
	}

	var agentConfig types.AgentConfig
	if err := yaml.Unmarshal(data, &agentConfig); err != nil {
		return nil, fmt.Errorf("failed to parse agent config: %w", err)
	}

	return &agentConfig, nil
}

// resolveEnvVars resolves environment variables in the configuration.
func (l *Loader) resolveEnvVars(config *types.Config) error {
	// Check LLM API keys availability (but don't require ALL of them)
	// At least ONE provider must be available - this is validated in LLM Manager
	availableProviders := 0
	for key, envVar := range config.LLMs.APIKeys {
		if envVar.Env != "" {
			if value := os.Getenv(envVar.Env); value != "" {
				availableProviders++
			} else {
				// Warn but don't fail - users can run with just one provider
				fmt.Printf("Warning: environment variable %s (for %s provider) is not set - this provider will be unavailable\n", envVar.Env, key)
			}
		}
	}

	// Provide helpful message if no providers are available
	if availableProviders == 0 {
		return fmt.Errorf("no LLM provider API keys are set - at least one is required (GEMINI_API_KEY, OPENAI_API_KEY, CLAUDE_API_KEY, etc.)")
	}

	// Resolve Slack token
	if config.Slack != nil && config.Slack.Enabled {
		if config.Slack.Token.Env != "" {
			if value := os.Getenv(config.Slack.Token.Env); value == "" {
				return fmt.Errorf("environment variable %s (for Slack token) is not set", config.Slack.Token.Env)
			}
		}
	}

	return nil
}

// GetEnvValue retrieves the actual value from an EnvironmentVariable.
func GetEnvValue(envVar types.EnvironmentVariable) string {
	if envVar.Env != "" {
		return os.Getenv(envVar.Env)
	}
	return ""
}
