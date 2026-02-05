package llm

import (
	"context"
	"os"
	"testing"
)

const demoAPIKey = "demo-key"

func TestOpenAIProvider(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" || apiKey == demoAPIKey {
		t.Skip("Skipping OpenAI test: OPENAI_API_KEY not set")
	}

	provider, err := NewOpenAIProvider(apiKey, "")
	if err != nil {
		t.Fatalf("Failed to create OpenAI provider: %v", err)
	}

	ctx := context.Background()
	prompt := "Say 'Hello from OpenAI' and nothing else."

	opts := &GenerateOptions{
		Temperature: 0.7,
		MaxTokens:   50,
	}

	response, err := provider.Generate(ctx, prompt, opts)
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	if response == "" {
		t.Error("Expected non-empty response")
	}

	t.Logf("OpenAI Response: %s", response)

	// Test with system prompt
	opts.SystemPrompt = "You are a helpful assistant."
	response2, err := provider.Generate(ctx, "What is 2+2?", opts)
	if err != nil {
		t.Fatalf("Failed to generate with system prompt: %v", err)
	}

	if response2 == "" {
		t.Error("Expected non-empty response with system prompt")
	}

	t.Logf("OpenAI Response with system prompt: %s", response2)
}

func TestOpenAIProviderWithCustomModel(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" || apiKey == "demo-key" {
		t.Skip("Skipping OpenAI test: OPENAI_API_KEY not set")
	}

	// Test with GPT-3.5
	provider, err := NewOpenAIProvider(apiKey, "gpt-3.5-turbo")
	if err != nil {
		t.Fatalf("Failed to create OpenAI provider: %v", err)
	}

	ctx := context.Background()
	prompt := "Say 'Hello from GPT-3.5' briefly."

	response, err := provider.Generate(ctx, prompt, nil)
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	if response == "" {
		t.Error("Expected non-empty response")
	}

	t.Logf("GPT-3.5 Response: %s", response)
}

func TestClaudeProvider(t *testing.T) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" || apiKey == "demo-key" {
		t.Skip("Skipping Claude test: CLAUDE_API_KEY not set")
	}

	provider, err := NewClaudeProvider(apiKey, "")
	if err != nil {
		t.Fatalf("Failed to create Claude provider: %v", err)
	}

	ctx := context.Background()
	prompt := "Say 'Hello from Claude' and nothing else."

	opts := &GenerateOptions{
		Temperature: 0.7,
		MaxTokens:   50,
	}

	response, err := provider.Generate(ctx, prompt, opts)
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	if response == "" {
		t.Error("Expected non-empty response")
	}

	t.Logf("Claude Response: %s", response)

	// Test with system prompt
	opts.SystemPrompt = "You are a helpful assistant."
	response2, err := provider.Generate(ctx, "What is the capital of France?", opts)
	if err != nil {
		t.Fatalf("Failed to generate with system prompt: %v", err)
	}

	if response2 == "" {
		t.Error("Expected non-empty response with system prompt")
	}

	t.Logf("Claude Response with system prompt: %s", response2)
}

func TestClaudeProviderWithCustomModel(t *testing.T) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" || apiKey == "demo-key" {
		t.Skip("Skipping Claude test: CLAUDE_API_KEY not set")
	}

	// Test with Claude 3 Haiku (faster, cheaper model)
	provider, err := NewClaudeProvider(apiKey, "claude-3-haiku-20240307")
	if err != nil {
		t.Fatalf("Failed to create Claude provider: %v", err)
	}

	ctx := context.Background()
	prompt := "Say 'Hello from Claude Haiku' briefly."

	response, err := provider.Generate(ctx, prompt, nil)
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	if response == "" {
		t.Error("Expected non-empty response")
	}

	t.Logf("Claude Haiku Response: %s", response)
}

func TestProviderComparison(t *testing.T) {
	// This test compares responses from different providers
	// Skip if API keys aren't available

	geminiKey := os.Getenv("GEMINI_API_KEY")
	openaiKey := os.Getenv("OPENAI_API_KEY")
	claudeKey := os.Getenv("CLAUDE_API_KEY")

	hasGemini := geminiKey != "" && geminiKey != "demo-key"
	hasOpenAI := openaiKey != "" && openaiKey != "demo-key"
	hasClaude := claudeKey != "" && claudeKey != "demo-key"

	if !hasGemini && !hasOpenAI && !hasClaude {
		t.Skip("Skipping comparison test: no real API keys available")
	}

	ctx := context.Background()
	prompt := "Write a one-line function in Python that adds two numbers."

	t.Log("Comparing provider responses for:", prompt)
	t.Log("")

	if hasGemini {
		provider, _ := NewGeminiProvider(geminiKey)
		response, err := provider.Generate(ctx, prompt, &GenerateOptions{Temperature: 0.7, MaxTokens: 100})
		if err == nil {
			t.Logf("Gemini: %s", response)
		} else {
			t.Logf("Gemini error: %v", err)
		}
	}

	if hasOpenAI {
		provider, _ := NewOpenAIProvider(openaiKey, "")
		response, err := provider.Generate(ctx, prompt, &GenerateOptions{Temperature: 0.7, MaxTokens: 100})
		if err == nil {
			t.Logf("OpenAI: %s", response)
		} else {
			t.Logf("OpenAI error: %v", err)
		}
	}

	if hasClaude {
		provider, _ := NewClaudeProvider(claudeKey, "")
		response, err := provider.Generate(ctx, prompt, &GenerateOptions{Temperature: 0.7, MaxTokens: 100})
		if err == nil {
			t.Logf("Claude: %s", response)
		} else {
			t.Logf("Claude error: %v", err)
		}
	}
}

func TestOpenAIProviderNoAPIKey(t *testing.T) {
	_, err := NewOpenAIProvider("", "")
	if err == nil {
		t.Error("Expected error when creating OpenAI provider without API key")
	}
}

func TestClaudeProviderNoAPIKey(t *testing.T) {
	_, err := NewClaudeProvider("", "")
	if err == nil {
		t.Error("Expected error when creating Claude provider without API key")
	}
}
