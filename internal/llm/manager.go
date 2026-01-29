package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/kpango/BuildBureau/internal/config"
	"github.com/kpango/BuildBureau/pkg/types"
)

// Provider represents an LLM provider interface.
type Provider interface {
	// Generate sends a prompt to the LLM and returns the response
	Generate(ctx context.Context, prompt string, opts *GenerateOptions) (string, error)

	// Name returns the name of the provider
	Name() string
}

// GenerateOptions contains options for generation.
type GenerateOptions struct {
	SystemPrompt string
	Temperature  float64
	MaxTokens    int
}

// Manager manages multiple LLM providers.
type Manager struct {
	providers    map[string]Provider
	defaultModel string
}

// NewManager creates a new LLM manager with real provider initialization.
func NewManager(cfg *types.LLMConfig) (*Manager, error) {
	m := &Manager{
		providers:    make(map[string]Provider),
		defaultModel: cfg.DefaultModel,
	}

	// Initialize Gemini provider if API key is available
	if geminiKey, exists := cfg.APIKeys["gemini"]; exists {
		apiKey := config.GetEnvValue(geminiKey)
		if apiKey != "" {
			provider, err := NewGeminiProvider(apiKey)
			if err != nil {
				return nil, fmt.Errorf("failed to initialize Gemini provider: %w", err)
			}
			m.providers["gemini"] = provider
		}
	}

	// Initialize OpenAI provider if API key is available
	if openaiKey, exists := cfg.APIKeys["openai"]; exists {
		apiKey := config.GetEnvValue(openaiKey)
		if apiKey != "" {
			// Use model from environment or default
			model := os.Getenv("OPENAI_MODEL")
			provider, err := NewOpenAIProvider(apiKey, model)
			if err != nil {
				fmt.Printf("Warning: failed to initialize OpenAI provider: %v\n", err)
			} else {
				m.providers["openai"] = provider
			}
		}
	}

	// Initialize Claude provider if API key is available
	if claudeKey, exists := cfg.APIKeys["claude"]; exists {
		apiKey := config.GetEnvValue(claudeKey)
		if apiKey != "" {
			// Use model from environment or default
			model := os.Getenv("CLAUDE_MODEL")
			provider, err := NewClaudeProvider(apiKey, model)
			if err != nil {
				fmt.Printf("Warning: failed to initialize Claude provider: %v\n", err)
			} else {
				m.providers["claude"] = provider
			}
		}
	}

	// Initialize remote providers for Codex, Qwen, or custom endpoints
	remoteProviders := []struct {
		name     string
		endpoint string
	}{
		{"codex", os.Getenv("CODEX_ENDPOINT")},
		{"qwen", os.Getenv("QWEN_ENDPOINT")},
		{"custom", os.Getenv("CUSTOM_LLM_ENDPOINT")},
	}

	for _, rp := range remoteProviders {
		if envVar, exists := cfg.APIKeys[rp.name]; exists {
			apiKey := config.GetEnvValue(envVar)
			if apiKey != "" && rp.endpoint != "" {
				provider, err := NewRemoteProvider(rp.name, rp.endpoint, apiKey)
				if err != nil {
					// Log but don't fail - remote providers are optional
					fmt.Printf("Warning: failed to initialize %s provider: %v\n", rp.name, err)
					continue
				}
				m.providers[rp.name] = provider
			}
		}
	}

	if len(m.providers) == 0 {
		return nil, fmt.Errorf("no LLM providers could be initialized")
	}

	return m, nil
}

// Generate sends a prompt to the specified model or default.
func (m *Manager) Generate(ctx context.Context, model, prompt string, opts *GenerateOptions) (string, error) {
	if model == "" {
		model = m.defaultModel
	}

	provider, ok := m.providers[model]
	if !ok {
		return "", fmt.Errorf("model %s not available", model)
	}

	return provider.Generate(ctx, prompt, opts)
}

// GetProvider returns a specific provider.
func (m *Manager) GetProvider(name string) (Provider, error) {
	provider, ok := m.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return provider, nil
}

// AddProvider adds a new provider.
func (m *Manager) AddProvider(name string, provider Provider) {
	m.providers[name] = provider
}

// Close closes all providers.
func (m *Manager) Close() error {
	for name, provider := range m.providers {
		if closer, ok := provider.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				fmt.Printf("Warning: failed to close provider %s: %v\n", name, err)
			}
		}
	}
	return nil
}
