package adk

import (
	"context"
	"fmt"
)

// LLMClient abstracts the interaction with different LLM providers.
type LLMClient interface {
	Generate(ctx context.Context, systemPrompt string, userPrompt string, modelID string) (string, error)
}

// MockLLMClient is a placeholder for testing without real APIs.
type MockLLMClient struct {
	// We can add some preset responses here if needed
}

func NewMockLLMClient() *MockLLMClient {
	return &MockLLMClient{}
}

func (m *MockLLMClient) Generate(ctx context.Context, systemPrompt string, userPrompt string, modelID string) (string, error) {
	// Simple mock behavior: return a dummy JSON based on the prompt content or just a generic success message.
	return fmt.Sprintf("Mocked LLM response for model %s", modelID), nil
}

// RealLLMClient would implement the actual API calls using the keys.
type RealLLMClient struct {
	APIKeys map[string]string // modelID -> key
}

func NewRealLLMClient(keys map[string]string) *RealLLMClient {
	return &RealLLMClient{APIKeys: keys}
}

func (c *RealLLMClient) Generate(ctx context.Context, systemPrompt string, userPrompt string, modelID string) (string, error) {
	apiKey, ok := c.APIKeys[modelID]
	if !ok || apiKey == "" {
		return "", fmt.Errorf("no API key found for model %s", modelID)
	}

	// Logic to call OpenAI / Gemini / Claude based on modelID.
	return fmt.Sprintf("Real LLM Call (Simulated) for %s. System: %s... User: %s...", modelID, systemPrompt[:10], userPrompt[:10]), nil
}
