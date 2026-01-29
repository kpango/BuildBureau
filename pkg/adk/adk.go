package adk

import (
	"context"
	"fmt"
)

// Model represents a generative AI model.
type Model interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}

// Config represents the ADK configuration.
type Config struct {
	ProjectID string
	Location  string
}

// NewClient creates a new ADK client.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	return &Client{}, nil
}

// Client is the main entry point for the ADK.
type Client struct{}

// Model returns a model instance.
func (c *Client) Model(name string) Model {
	return &mockModel{name: name}
}

type mockModel struct {
	name string
}

func (m *mockModel) GenerateContent(ctx context.Context, prompt string) (string, error) {
	// In a real implementation, this would call the API.
	// Here we return a dummy response.
	return fmt.Sprintf("[Mock %s] Processed: %s...", m.name, prompt[:min(len(prompt), 20)]), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RemoteAgentConfig defines a remote agent.
type RemoteAgentConfig struct {
	Name         string
	Endpoint     string
	Capabilities []string
}
