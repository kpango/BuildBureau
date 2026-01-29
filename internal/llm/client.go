package llm

import (
	"context"
	"fmt"
)

// Message represents a chat message
type Message struct {
	Role    string // system, user, assistant
	Content string
}

// Request represents an LLM request
type Request struct {
	Messages    []Message
	Model       string
	Temperature float64
	MaxTokens   int
}

// Response represents an LLM response
type Response struct {
	Content      string
	FinishReason string
	TokensUsed   int
}

// Client is the interface for LLM providers
type Client interface {
	// Generate generates a response from the LLM
	Generate(ctx context.Context, req *Request) (*Response, error)

	// StreamGenerate generates a streaming response
	StreamGenerate(ctx context.Context, req *Request) (<-chan string, <-chan error)
}

// MockClient is a mock implementation for testing
type MockClient struct {
	responses []string
	callCount int
}

// NewMockClient creates a new mock client
func NewMockClient(responses []string) *MockClient {
	if len(responses) == 0 {
		responses = []string{"Mock response"}
	}
	return &MockClient{
		responses: responses,
	}
}

// Generate returns a mock response
func (c *MockClient) Generate(ctx context.Context, req *Request) (*Response, error) {
	if c.callCount >= len(c.responses) {
		c.callCount = 0
	}
	response := c.responses[c.callCount]
	c.callCount++

	return &Response{
		Content:      response,
		FinishReason: "stop",
		TokensUsed:   len(response) / 4, // rough estimation
	}, nil
}

// StreamGenerate returns a mock streaming response
func (c *MockClient) StreamGenerate(ctx context.Context, req *Request) (<-chan string, <-chan error) {
	contentCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(contentCh)
		defer close(errCh)

		if c.callCount >= len(c.responses) {
			c.callCount = 0
		}
		response := c.responses[c.callCount]
		c.callCount++

		contentCh <- response
	}()

	return contentCh, errCh
}

// GoogleADKClient is a placeholder for Google ADK integration
type GoogleADKClient struct {
	apiKey string
	model  string
}

// NewGoogleADKClient creates a new Google ADK client (placeholder)
func NewGoogleADKClient(apiKey, model string) *GoogleADKClient {
	return &GoogleADKClient{
		apiKey: apiKey,
		model:  model,
	}
}

// Generate generates a response using Google ADK (placeholder)
func (c *GoogleADKClient) Generate(ctx context.Context, req *Request) (*Response, error) {
	// Placeholder implementation
	// TODO: Implement actual Google ADK integration
	return nil, fmt.Errorf("Google ADK integration not yet implemented")
}

// StreamGenerate generates a streaming response using Google ADK (placeholder)
func (c *GoogleADKClient) StreamGenerate(ctx context.Context, req *Request) (<-chan string, <-chan error) {
	contentCh := make(chan string)
	errCh := make(chan error, 1)

	go func() {
		defer close(contentCh)
		defer close(errCh)
		errCh <- fmt.Errorf("Google ADK integration not yet implemented")
	}()

	return contentCh, errCh
}

// ClientFactory creates LLM clients based on configuration
type ClientFactory struct {
	provider string
	apiKey   string
	model    string
}

// NewClientFactory creates a new client factory
func NewClientFactory(provider, apiKey, model string) *ClientFactory {
	return &ClientFactory{
		provider: provider,
		apiKey:   apiKey,
		model:    model,
	}
}

// Create creates an appropriate LLM client
func (f *ClientFactory) Create() (Client, error) {
	switch f.provider {
	case "google":
		return NewGoogleADKClient(f.apiKey, f.model), nil
	case "mock":
		return NewMockClient(nil), nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", f.provider)
	}
}
