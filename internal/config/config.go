package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the complete system configuration
type Config struct {
	Agents AgentsConfig `yaml:"agents"`
	LLM    LLMConfig    `yaml:"llm"`
	GRPC   GRPCConfig   `yaml:"grpc"`
	Slack  SlackConfig  `yaml:"slack"`
	UI     UIConfig     `yaml:"ui"`
	System SystemConfig `yaml:"system"`
}

// AgentsConfig contains configuration for all agent types
type AgentsConfig struct {
	President           AgentConfig `yaml:"president"`
	PresidentSecretary  AgentConfig `yaml:"president_secretary"`
	DepartmentManager   AgentConfig `yaml:"department_manager"`
	DepartmentSecretary AgentConfig `yaml:"department_secretary"`
	SectionManager      AgentConfig `yaml:"section_manager"`
	SectionSecretary    AgentConfig `yaml:"section_secretary"`
	Employee            AgentConfig `yaml:"employee"`
}

// AgentConfig represents configuration for a single agent type
type AgentConfig struct {
	Count       int      `yaml:"count"`
	Model       string   `yaml:"model"`
	Instruction string   `yaml:"instruction"`
	AllowTools  bool     `yaml:"allowTools"`
	Tools       []string `yaml:"tools"`
	Timeout     int      `yaml:"timeout"`
	RetryCount  int      `yaml:"retryCount"`
}

// LLMConfig represents LLM provider configuration
type LLMConfig struct {
	Provider     string  `yaml:"provider"`
	APIEndpoint  string  `yaml:"apiEndpoint"`
	DefaultModel string  `yaml:"defaultModel"`
	MaxTokens    int     `yaml:"maxTokens"`
	Temperature  float64 `yaml:"temperature"`
	TopP         float64 `yaml:"topP"`
}

// GRPCConfig represents gRPC server configuration
type GRPCConfig struct {
	Port             int  `yaml:"port"`
	MaxMessageSize   int  `yaml:"maxMessageSize"`
	Timeout          int  `yaml:"timeout"`
	EnableReflection bool `yaml:"enableReflection"`
}

// SlackConfig represents Slack integration configuration
type SlackConfig struct {
	Enabled       bool                       `yaml:"enabled"`
	Token         string                     `yaml:"token"`
	ChannelID     string                     `yaml:"channelID"`
	Notifications NotificationsConfig        `yaml:"notifications"`
	RetryCount    int                        `yaml:"retryCount"`
	Timeout       int                        `yaml:"timeout"`
}

// NotificationsConfig represents notification settings
type NotificationsConfig struct {
	ProjectStart    NotificationConfig `yaml:"projectStart"`
	TaskComplete    NotificationConfig `yaml:"taskComplete"`
	Error           NotificationConfig `yaml:"error"`
	ProjectComplete NotificationConfig `yaml:"projectComplete"`
}

// NotificationConfig represents a single notification type
type NotificationConfig struct {
	Enabled bool   `yaml:"enabled"`
	Message string `yaml:"message"`
}

// UIConfig represents Terminal UI configuration
type UIConfig struct {
	EnableTUI   bool   `yaml:"enableTUI"`
	RefreshRate int    `yaml:"refreshRate"`
	Theme       string `yaml:"theme"`
	ShowProgress bool  `yaml:"showProgress"`
	LogLevel    string `yaml:"logLevel"`
}

// SystemConfig represents system-wide configuration
type SystemConfig struct {
	WorkDir           string `yaml:"workDir"`
	LogDir            string `yaml:"logDir"`
	CacheDir          string `yaml:"cacheDir"`
	MaxConcurrentTasks int   `yaml:"maxConcurrentTasks"`
	GlobalTimeout     int    `yaml:"globalTimeout"`
}

// Load reads and parses the configuration file
func Load(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables in the config
	expandedData := os.ExpandEnv(string(data))

	var config Config
	if err := yaml.Unmarshal([]byte(expandedData), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Agents.President.Count < 1 {
		return fmt.Errorf("president agent count must be at least 1")
	}
	if c.Agents.DepartmentManager.Count < 1 {
		return fmt.Errorf("department manager agent count must be at least 1")
	}
	if c.LLM.Provider == "" {
		return fmt.Errorf("LLM provider must be specified")
	}
	if c.GRPC.Port <= 0 {
		return fmt.Errorf("gRPC port must be positive")
	}
	if c.Slack.Enabled && c.Slack.Token == "" {
		return fmt.Errorf("Slack token must be provided when Slack is enabled")
	}
	return nil
}
