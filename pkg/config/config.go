package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the root configuration.
type Config struct {
	Organization Organization `yaml:"organization"`
	Slack        Slack        `yaml:"slack"`
	LLMs         LLMs         `yaml:"llms"`
	SubAgents    []SubAgent   `yaml:"sub_agents"`
}

type Organization struct {
	Layers []Layer `yaml:"layers"`
}

type Layer struct {
	Name     string   `yaml:"name"`
	Agent    string   `yaml:"agent"`
	Count    int      `yaml:"count,omitempty"`
	AttachTo []string `yaml:"attach_to,omitempty"`
}

type Slack struct {
	Enabled  bool     `yaml:"enabled"`
	Token    Secret   `yaml:"token"`
	Channels []string `yaml:"channels"`
	NotifyOn []string `yaml:"notify_on"`
}

type LLMs struct {
	DefaultModel string            `yaml:"default_model"`
	APIKeys      map[string]Secret `yaml:"api_keys"`
}

type SubAgent struct {
	Name   string      `yaml:"name"`
	Remote RemoteAgent `yaml:"remote"`
}

type RemoteAgent struct {
	Endpoint     string   `yaml:"endpoint"`
	Capabilities []string `yaml:"capabilities"`
}

// Secret handles values that can be literal strings or environment variable references.
type Secret string

func (s *Secret) UnmarshalYAML(value *yaml.Node) error {
	var tempStr string
	if err := value.Decode(&tempStr); err == nil {
		*s = Secret(tempStr)
		return nil
	}

	var tempMap map[string]string
	if err := value.Decode(&tempMap); err == nil {
		if envKey, ok := tempMap["env"]; ok {
			*s = Secret(os.Getenv(envKey))
			return nil
		}
	}

	return fmt.Errorf("invalid secret format")
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
