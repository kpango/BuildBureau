package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	LLM    LLMConfig    `yaml:"llm"`
	Agents AgentsConfig `yaml:"agents"`
	Slack  SlackConfig  `yaml:"slack"`
	UI     UIConfig     `yaml:"ui"`
}

// LLMConfig contains LLM provider settings
type LLMConfig struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"api_key"`
	Model    string `yaml:"model"`
}

// AgentsConfig contains agent-specific configurations
type AgentsConfig struct {
	CEO      AgentConfig `yaml:"ceo"`
	Manager  AgentConfig `yaml:"manager"`
	Lead     AgentConfig `yaml:"lead"`
	Employee AgentConfig `yaml:"employee"`
}

// AgentConfig represents configuration for a single agent
type AgentConfig struct {
	Name    string `yaml:"name"`
	Enabled bool   `yaml:"enabled"`
}

// SlackConfig contains Slack integration settings
type SlackConfig struct {
	Enabled       bool               `yaml:"enabled"`
	Token         string             `yaml:"token"`
	Channels      ChannelMapping     `yaml:"channels"`
	Notifications NotificationConfig `yaml:"notifications"`
}

// ChannelMapping maps event types to Slack channels
type ChannelMapping struct {
	Management string `yaml:"management"`
	Updates    string `yaml:"updates"`
	Dev        string `yaml:"dev"`
}

// NotificationConfig specifies when to send notifications
type NotificationConfig struct {
	NotifyOnTaskAssigned  []string `yaml:"notify_on_task_assigned"`
	NotifyOnTaskCompleted []string `yaml:"notify_on_task_completed"`
	NotifyOnError         bool     `yaml:"notify_on_error"`
}

// UIConfig contains UI settings
type UIConfig struct {
	ShowAgentLogs bool `yaml:"show_agent_logs"`
	ColorCoding   bool `yaml:"color_coding"`
	RefreshRateMs int  `yaml:"refresh_rate_ms"`
}

// Load reads and parses the configuration file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Expand environment variables
	content := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, err
	}

	// Post-process environment variables that yaml didn't expand
	cfg.LLM.APIKey = expandEnv(cfg.LLM.APIKey)
	cfg.Slack.Token = expandEnv(cfg.Slack.Token)

	return &cfg, nil
}

// expandEnv expands ${VAR} or $VAR in strings
func expandEnv(s string) string {
	return os.ExpandEnv(s)
}

// Default returns a default configuration
func Default() *Config {
	return &Config{
		LLM: LLMConfig{
			Provider: "gemini",
			Model:    "gemini-1.5-pro",
		},
		Agents: AgentsConfig{
			CEO:      AgentConfig{Name: "CEO Agent", Enabled: true},
			Manager:  AgentConfig{Name: "Manager Agent", Enabled: true},
			Lead:     AgentConfig{Name: "Lead Agent", Enabled: true},
			Employee: AgentConfig{Name: "Employee Agent", Enabled: true},
		},
		Slack: SlackConfig{
			Enabled: false,
		},
		UI: UIConfig{
			ShowAgentLogs: true,
			ColorCoding:   true,
			RefreshRateMs: 100,
		},
	}
}
