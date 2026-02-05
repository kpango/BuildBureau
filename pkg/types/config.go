package types

// Config represents the main configuration structure for BuildBureau.
type Config struct {
	LLMs         LLMConfig          `yaml:"llms"`
	Slack        *SlackConfig       `yaml:"slack,omitempty"`
	Memory       *MemoryConfig      `yaml:"memory,omitempty"`
	Organization OrganizationConfig `yaml:"organization"`
}

// OrganizationConfig defines the agent hierarchy.
type OrganizationConfig struct {
	Layers []LayerConfig `yaml:"layers"`
}

// LayerConfig defines a layer in the organization.
type LayerConfig struct {
	Name     string   `yaml:"name"`
	Agent    string   `yaml:"agent,omitempty"`
	AttachTo []string `yaml:"attach_to,omitempty"`
	Count    int      `yaml:"count,omitempty"`
}

// SlackConfig defines Slack notification settings.
type SlackConfig struct {
	Token    EnvironmentVariable `yaml:"token"`
	Channels []string            `yaml:"channels"`
	NotifyOn []string            `yaml:"notify_on"`
	Enabled  bool                `yaml:"enabled"`
}

// LLMConfig defines LLM configuration.
type LLMConfig struct {
	APIKeys      map[string]EnvironmentVariable `yaml:"api_keys"`
	DefaultModel string                         `yaml:"default_model"`
}

// EnvironmentVariable represents a value that comes from an environment variable.
type EnvironmentVariable struct {
	Env string `yaml:"env"`
}

// AgentConfig represents the configuration for an individual agent.
type AgentConfig struct {
	Name         string           `yaml:"name"`
	Role         string           `yaml:"role"`
	Description  string           `yaml:"description"`
	Model        string           `yaml:"model,omitempty"`
	SystemPrompt string           `yaml:"system_prompt"`
	SubAgents    []SubAgentConfig `yaml:"sub_agents,omitempty"`
	Capabilities []string         `yaml:"capabilities,omitempty"`
}

// SubAgentConfig represents a sub-agent configuration (for remote agents).
type SubAgentConfig struct {
	Name         string        `yaml:"name"`
	Remote       *RemoteConfig `yaml:"remote,omitempty"`
	Capabilities []string      `yaml:"capabilities,omitempty"`
}

// RemoteConfig defines a remote agent endpoint.
type RemoteConfig struct {
	Endpoint     string   `yaml:"endpoint"`
	Capabilities []string `yaml:"capabilities,omitempty"`
}

// MemoryConfig represents memory storage configuration.
type MemoryConfig struct {
	SQLite    SQLiteConfig    `yaml:"sqlite"`
	Vald      ValdConfig      `yaml:"vald"`
	Retention RetentionConfig `yaml:"retention"`
	Enabled   bool            `yaml:"enabled"`
}

// SQLiteConfig represents SQLite database configuration.
type SQLiteConfig struct {
	Path     string `yaml:"path"`
	Enabled  bool   `yaml:"enabled"`
	InMemory bool   `yaml:"in_memory"`
}

// ValdConfig represents Vald vector database configuration.
type ValdConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Dimension int    `yaml:"dimension"`
	PoolSize  int    `yaml:"pool_size"`
	Enabled   bool   `yaml:"enabled"`
}

// RetentionConfig represents memory retention policies.
type RetentionConfig struct {
	ConversationDays int `yaml:"conversation_days"` // 0 = forever
	TaskDays         int `yaml:"task_days"`
	KnowledgeDays    int `yaml:"knowledge_days"`
	MaxEntries       int `yaml:"max_entries"` // 0 = unlimited
}
