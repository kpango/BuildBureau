package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the complete BuildBureau configuration
type Config struct {
	Hierarchy HierarchyConfig        `yaml:"hierarchy"`
	Agents    map[string]AgentConfig `yaml:"agents"`
	Slack     SlackConfig            `yaml:"slack"`
	System    SystemConfig           `yaml:"system"`
}

// HierarchyConfig defines the organizational structure
type HierarchyConfig struct {
	Departments           int      `yaml:"departments"`
	ManagersPerDepartment int      `yaml:"managers_per_department"`
	ManagerSpecialties    []string `yaml:"manager_specialties"`
	WorkersPerManager     int      `yaml:"workers_per_manager"`
}

// AgentConfig defines configuration for an agent
type AgentConfig struct {
	Model       string   `yaml:"model"`
	Instruction string   `yaml:"instruction"`
	Tools       []string `yaml:"tools"`
	Temperature float64  `yaml:"temperature"`
}

// SlackConfig defines Slack integration settings
type SlackConfig struct {
	Enabled       bool                    `yaml:"enabled"`
	Token         string                  `yaml:"token"`
	Channels      SlackChannels           `yaml:"channels"`
	NotifyOn      map[string]NotifyConfig `yaml:"notify_on"`
	MessageFormat MessageFormatConfig     `yaml:"message_format"`
}

// SlackChannels defines channel mappings
type SlackChannels struct {
	Main        string `yaml:"main"`
	Management  string `yaml:"management"`
	Development string `yaml:"development"`
	Errors      string `yaml:"errors"`
}

// NotifyConfig defines when and how to notify
type NotifyConfig struct {
	Enabled bool     `yaml:"enabled"`
	Roles   []string `yaml:"roles"`
	Channel string   `yaml:"channel"`
}

// MessageFormatConfig defines message formatting options
type MessageFormatConfig struct {
	Prefix           string `yaml:"prefix"`
	IncludeTimestamp bool   `yaml:"include_timestamp"`
	IncludeAgentName bool   `yaml:"include_agent_name"`
	UseThreads       bool   `yaml:"use_threads"`
}

// SystemConfig defines system-level settings
type SystemConfig struct {
	Logging       LoggingConfig       `yaml:"logging"`
	APIKeys       APIKeysConfig       `yaml:"api_keys"`
	GRPC          GRPCConfig          `yaml:"grpc"`
	Timeouts      TimeoutsConfig      `yaml:"timeouts"`
	KnowledgeBase KnowledgeBaseConfig `yaml:"knowledge_base"`
	UI            UIConfig            `yaml:"ui"`
}

// LoggingConfig defines logging settings
type LoggingConfig struct {
	Level             string `yaml:"level"`
	Output            string `yaml:"output"`
	File              string `yaml:"file"`
	EnableFileLogging bool   `yaml:"enable_file_logging"`
}

// APIKeysConfig stores API keys
type APIKeysConfig struct {
	Gemini string `yaml:"gemini"`
	Search string `yaml:"search"`
}

// GRPCConfig defines gRPC server settings
type GRPCConfig struct {
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	EnableReflection bool   `yaml:"enable_reflection"`
}

// TimeoutsConfig defines timeout settings
type TimeoutsConfig struct {
	AgentResponseSeconds int `yaml:"agent_response_seconds"`
	TaskExecutionSeconds int `yaml:"task_execution_seconds"`
	MaxRetries           int `yaml:"max_retries"`
	RetryBackoffSeconds  int `yaml:"retry_backoff_seconds"`
}

// KnowledgeBaseConfig defines knowledge base settings
type KnowledgeBaseConfig struct {
	Type             string `yaml:"type"`
	Path             string `yaml:"path"`
	EnableVersioning bool   `yaml:"enable_versioning"`
}

// UIConfig defines UI settings
type UIConfig struct {
	Enabled            bool   `yaml:"enabled"`
	Theme              string `yaml:"theme"`
	UpdateIntervalMS   int    `yaml:"update_interval_ms"`
	ShowAgentHierarchy bool   `yaml:"show_agent_hierarchy"`
	EnableColors       bool   `yaml:"enable_colors"`
	Language           string `yaml:"language"`
	MaxHistoryLines    int    `yaml:"max_history_lines"`
}

// Load reads and parses the configuration file
func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables in the config
	expandedData := expandEnvVars(string(data))

	var config Config
	if err := yaml.Unmarshal([]byte(expandedData), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// expandEnvVars expands environment variables in the format ${VAR_NAME}
func expandEnvVars(s string) string {
	return os.Expand(s, func(key string) string {
		return os.Getenv(key)
	})
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate hierarchy
	if c.Hierarchy.Departments < 1 {
		return fmt.Errorf("departments must be at least 1")
	}
	if c.Hierarchy.ManagersPerDepartment < 1 {
		return fmt.Errorf("managers_per_department must be at least 1")
	}
	if len(c.Hierarchy.ManagerSpecialties) != c.Hierarchy.ManagersPerDepartment {
		return fmt.Errorf("number of manager specialties must match managers_per_department")
	}
	if c.Hierarchy.WorkersPerManager < 1 {
		return fmt.Errorf("workers_per_manager must be at least 1")
	}

	// Validate required agents
	requiredAgents := []string{"ceo", "ceo_secretary", "dept_head", "dept_head_secretary", "manager", "manager_secretary", "worker"}
	for _, agent := range requiredAgents {
		if _, ok := c.Agents[agent]; !ok {
			return fmt.Errorf("missing required agent configuration: %s", agent)
		}
	}

	// Validate Slack config if enabled
	if c.Slack.Enabled {
		if c.Slack.Token == "" {
			return fmt.Errorf("slack token is required when slack is enabled")
		}
	}

	// Validate system config
	if c.System.APIKeys.Gemini == "" {
		return fmt.Errorf("gemini API key is required")
	}

	return nil
}

// GetChannelForEvent returns the appropriate Slack channel for an event
func (c *Config) GetChannelForEvent(eventType string) string {
	notifyConfig, ok := c.Slack.NotifyOn[eventType]
	if !ok {
		return c.Slack.Channels.Main
	}

	channelName := notifyConfig.Channel
	switch strings.ToLower(channelName) {
	case "main":
		return c.Slack.Channels.Main
	case "management":
		return c.Slack.Channels.Management
	case "development":
		return c.Slack.Channels.Development
	case "errors":
		return c.Slack.Channels.Errors
	default:
		return c.Slack.Channels.Main
	}
}

// ShouldNotify checks if an event should trigger a notification for a role
func (c *Config) ShouldNotify(eventType, role string) bool {
	if !c.Slack.Enabled {
		return false
	}

	notifyConfig, ok := c.Slack.NotifyOn[eventType]
	if !ok || !notifyConfig.Enabled {
		return false
	}

	for _, r := range notifyConfig.Roles {
		if strings.EqualFold(r, role) {
			return true
		}
	}

	return false
}
