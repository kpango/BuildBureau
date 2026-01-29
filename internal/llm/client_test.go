package llm

import (
	"context"
	"os"
	"testing"
)

func TestMockClient(t *testing.T) {
	client := NewMockClient([]string{"Test response"})
	
	ctx := context.Background()
	req := &Request{
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		Model:       "test",
		Temperature: 0.7,
	}
	
	resp, err := client.Generate(ctx, req)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	
	if resp.Content != "Test response" {
		t.Errorf("Expected 'Test response', got '%s'", resp.Content)
	}
	
	if resp.FinishReason != "stop" {
		t.Errorf("Expected finish reason 'stop', got '%s'", resp.FinishReason)
	}
}

func TestMockClientStreaming(t *testing.T) {
	client := NewMockClient([]string{"Streaming response"})
	
	ctx := context.Background()
	req := &Request{
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}
	
	contentCh, errCh := client.StreamGenerate(ctx, req)
	
	var content string
	for chunk := range contentCh {
		content += chunk
	}
	
	if err := <-errCh; err != nil {
		t.Fatalf("Stream error: %v", err)
	}
	
	if content != "Streaming response" {
		t.Errorf("Expected 'Streaming response', got '%s'", content)
	}
}

func TestGoogleADKClient_NoAPIKey(t *testing.T) {
	ctx := context.Background()
	_, err := NewGoogleADKClient(ctx, "", "gemini-pro")
	
	if err == nil {
		t.Error("Expected error when API key is empty")
	}
}

func TestGoogleADKClient_Integration(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_AI_API_KEY")
	if apiKey == "" {
		t.Skip("GOOGLE_AI_API_KEY not set, skipping integration test")
	}
	
	ctx := context.Background()
	client, err := NewGoogleADKClient(ctx, apiKey, "gemini-2.0-flash-exp")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	req := &Request{
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Say 'Hello, BuildBureau!' and nothing else."},
		},
		Temperature: 0.1,
		MaxTokens:   50,
	}
	
	resp, err := client.Generate(ctx, req)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	
	if resp.Content == "" {
		t.Error("Expected non-empty response content")
	}
	
	t.Logf("Response: %s", resp.Content)
	t.Logf("Tokens used: %d", resp.TokensUsed)
	t.Logf("Finish reason: %s", resp.FinishReason)
}

func TestGoogleADKClient_Streaming(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_AI_API_KEY")
	if apiKey == "" {
		t.Skip("GOOGLE_AI_API_KEY not set, skipping integration test")
	}
	
	ctx := context.Background()
	client, err := NewGoogleADKClient(ctx, apiKey, "gemini-2.0-flash-exp")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	req := &Request{
		Messages: []Message{
			{Role: "user", Content: "Count from 1 to 5."},
		},
		Temperature: 0.1,
		MaxTokens:   50,
	}
	
	contentCh, errCh := client.StreamGenerate(ctx, req)
	
	var fullContent string
	for chunk := range contentCh {
		fullContent += chunk
		t.Logf("Chunk: %s", chunk)
	}
	
	if err := <-errCh; err != nil {
		t.Fatalf("Stream error: %v", err)
	}
	
	if fullContent == "" {
		t.Error("Expected non-empty streamed content")
	}
	
	t.Logf("Full content: %s", fullContent)
}

func TestClientFactory(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		wantErr  bool
	}{
		{
			name:     "mock provider",
			provider: "mock",
			wantErr:  false,
		},
		{
			name:     "unsupported provider",
			provider: "unsupported",
			wantErr:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewClientFactory(tt.provider, "", "")
			_, err := factory.Create()
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientFactory_Google(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_AI_API_KEY")
	if apiKey == "" {
		t.Skip("GOOGLE_AI_API_KEY not set, skipping test")
	}
	
	factory := NewClientFactory("google", apiKey, "gemini-2.0-flash-exp")
	client, err := factory.Create()
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}
	
	if client == nil {
		t.Error("Expected non-nil client")
	}
}
