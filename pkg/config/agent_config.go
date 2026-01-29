package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AgentConfig struct {
	Role         string `yaml:"role"`
	SystemPrompt string `yaml:"system_prompt"`
}

func LoadAgentConfig(path string) (*AgentConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg AgentConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
