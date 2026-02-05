package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kpango/BuildBureau/internal/llm"
)

const demoAPIKey = "demo-key"

func main() {
	fmt.Println("=== BuildBureau Multi-Provider LLM Example ===\n")

	ctx := context.Background()
	prompt := "Write a simple 'Hello, World!' program in Python."

	// Test Gemini
	testGemini(ctx, prompt)
	fmt.Println()

	// Test OpenAI
	testOpenAI(ctx, prompt)
	fmt.Println()

	// Test Claude
	testClaude(ctx, prompt)
	fmt.Println()

	// Test all providers with comparison
	compareProviders(ctx)
}

func testGemini(ctx context.Context, prompt string) {
	fmt.Println("--- Testing Gemini ---")
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" || apiKey == demoAPIKey {
		fmt.Println("⚠️  GEMINI_API_KEY not set. Skipping.")
		return
	}

	provider, err := llm.NewGeminiProvider(apiKey)
	if err != nil {
		log.Printf("Failed to create Gemini provider: %v", err)
		return
	}
	defer provider.Close()

	opts := &llm.GenerateOptions{
		Temperature:  0.7,
		MaxTokens:    200,
		SystemPrompt: "You are a helpful coding assistant.",
	}

	response, err := provider.Generate(ctx, prompt, opts)
	if err != nil {
		log.Printf("Failed to generate with Gemini: %v", err)
		return
	}

	fmt.Printf("✅ Gemini Response:\n%s\n", response)
}

func testOpenAI(ctx context.Context, prompt string) {
	fmt.Println("--- Testing OpenAI ---")
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" || apiKey == "demo-key" {
		fmt.Println("⚠️  OPENAI_API_KEY not set. Skipping.")
		return
	}

	// You can specify a model, or leave empty for default (GPT-4 Turbo)
	model := os.Getenv("OPENAI_MODEL") // e.g., "gpt-4-turbo-preview" or "gpt-3.5-turbo"

	provider, err := llm.NewOpenAIProvider(apiKey, model)
	if err != nil {
		log.Printf("Failed to create OpenAI provider: %v", err)
		return
	}
	defer provider.Close()

	opts := &llm.GenerateOptions{
		Temperature:  0.7,
		MaxTokens:    200,
		SystemPrompt: "You are a helpful coding assistant.",
	}

	response, err := provider.Generate(ctx, prompt, opts)
	if err != nil {
		log.Printf("Failed to generate with OpenAI: %v", err)
		return
	}

	fmt.Printf("✅ OpenAI Response (model: %s):\n%s\n", provider.Name(), response)
}

func testClaude(ctx context.Context, prompt string) {
	fmt.Println("--- Testing Claude ---")
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" || apiKey == "demo-key" {
		fmt.Println("⚠️  CLAUDE_API_KEY not set. Skipping.")
		return
	}

	// You can specify a model, or leave empty for default (Claude 3.5 Sonnet)
	model := os.Getenv("CLAUDE_MODEL") // e.g., "claude-3-5-sonnet-20241022" or "claude-3-haiku-20240307"

	provider, err := llm.NewClaudeProvider(apiKey, model)
	if err != nil {
		log.Printf("Failed to create Claude provider: %v", err)
		return
	}
	defer provider.Close()

	opts := &llm.GenerateOptions{
		Temperature:  0.7,
		MaxTokens:    200,
		SystemPrompt: "You are a helpful coding assistant.",
	}

	response, err := provider.Generate(ctx, prompt, opts)
	if err != nil {
		log.Printf("Failed to generate with Claude: %v", err)
		return
	}

	fmt.Printf("✅ Claude Response:\n%s\n", response)
}

func compareProviders(ctx context.Context) {
	fmt.Println("--- Provider Comparison ---")
	fmt.Println("Testing all available providers with the same prompt:\n")

	prompt := "What is 2 + 2? Answer with just the number."

	providers := []struct {
		name     string
		envKey   string
		createFn func(string, string) (llm.Provider, error)
		modelEnv string
	}{
		{"Gemini", "GEMINI_API_KEY", func(key, _ string) (llm.Provider, error) { return llm.NewGeminiProvider(key) }, ""},
		{"OpenAI", "OPENAI_API_KEY", func(key, model string) (llm.Provider, error) { return llm.NewOpenAIProvider(key, model) }, "OPENAI_MODEL"},
		{"Claude", "CLAUDE_API_KEY", func(key, model string) (llm.Provider, error) { return llm.NewClaudeProvider(key, model) }, "CLAUDE_MODEL"},
	}

	opts := &llm.GenerateOptions{
		Temperature: 0.1, // Low temperature for consistent responses
		MaxTokens:   10,
	}

	for _, p := range providers {
		apiKey := os.Getenv(p.envKey)
		if apiKey == "" || apiKey == "demo-key" {
			fmt.Printf("⚠️  %s: API key not set\n", p.name)
			continue
		}

		model := ""
		if p.modelEnv != "" {
			model = os.Getenv(p.modelEnv)
		}

		provider, err := p.createFn(apiKey, model)
		if err != nil {
			fmt.Printf("❌ %s: Failed to create provider: %v\n", p.name, err)
			continue
		}

		response, err := provider.Generate(ctx, prompt, opts)
		if err != nil {
			fmt.Printf("❌ %s: Failed to generate: %v\n", p.name, err)
			continue
		}

		fmt.Printf("✅ %s: %s\n", p.name, response)

		if closer, ok := provider.(interface{ Close() error }); ok {
			closer.Close()
		}
	}

	fmt.Println("\n📊 All providers tested successfully!")
}
