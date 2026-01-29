package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kpango/BuildBureau/internal/llm"
)

// This example demonstrates using the Google ADK integration
// Set GOOGLE_AI_API_KEY environment variable to run this example
func main() {
	apiKey := os.Getenv("GOOGLE_AI_API_KEY")
	if apiKey == "" {
		log.Fatal("GOOGLE_AI_API_KEY environment variable must be set")
	}

	ctx := context.Background()

	// Create Google ADK client
	client, err := llm.NewGoogleADKClient(ctx, apiKey, "gemini-2.0-flash-exp")
	if err != nil {
		log.Fatalf("Failed to create Google ADK client: %v", err)
	}

	fmt.Println("=== Google ADK Integration Example ===")
	fmt.Println()

	// Example 1: Simple text generation
	fmt.Println("Example 1: Simple text generation")
	req := &llm.Request{
		Messages: []llm.Message{
			{Role: "system", Content: "You are a helpful assistant specialized in software architecture."},
			{Role: "user", Content: "Explain the benefits of a hierarchical multi-agent system in 2-3 sentences."},
		},
		Temperature: 0.7,
		MaxTokens:   200,
	}

	resp, err := client.Generate(ctx, req)
	if err != nil {
		log.Fatalf("Generation failed: %v", err)
	}

	fmt.Printf("Response: %s\n", resp.Content)
	fmt.Printf("Tokens used: %d\n", resp.TokensUsed)
	fmt.Printf("Finish reason: %s\n", resp.FinishReason)
	fmt.Println()

	// Example 2: Streaming response
	fmt.Println("Example 2: Streaming response")
	req2 := &llm.Request{
		Messages: []llm.Message{
			{Role: "user", Content: "Count from 1 to 5, one number per line."},
		},
		Temperature: 0.1,
		MaxTokens:   50,
	}

	contentCh, errCh := client.StreamGenerate(ctx, req2)

	fmt.Print("Streaming: ")
	for chunk := range contentCh {
		fmt.Print(chunk)
	}

	if err := <-errCh; err != nil {
		log.Fatalf("Streaming failed: %v", err)
	}

	fmt.Println()
	fmt.Println()
	fmt.Println("=== Example Complete ===")
}
