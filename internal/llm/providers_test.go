package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemoteProvider_Generate(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/v1/generate" {
			t.Errorf("Expected path /v1/generate, got %s", r.URL.Path)
		}

		// Check authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-api-key" {
			t.Errorf("Expected Authorization header 'Bearer test-api-key', got '%s'", authHeader)
		}

		// Decode request
		var req RemoteGenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		// Verify request content
		if req.Prompt == "" {
			t.Error("Expected non-empty prompt")
		}

		// Send response
		resp := RemoteGenerateResponse{
			Result: "Generated code: def hello(): print('Hello World')",
			Model:  "test-model",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create provider
	provider, err := NewRemoteProvider("test-model", server.URL, "test-api-key")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Test generation
	ctx := context.Background()
	result, err := provider.Generate(ctx, "Write a hello world function in Python", &GenerateOptions{
		Temperature: 0.7,
		MaxTokens:   100,
	})
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	if result != "Generated code: def hello(): print('Hello World')" {
		t.Errorf("Unexpected result: %s", result)
	}
}

func TestRemoteProvider_Generate_Error(t *testing.T) {
	// Create a mock HTTP server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}))
	defer server.Close()

	// Create provider
	provider, err := NewRemoteProvider("test-model", server.URL, "test-api-key")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Test generation
	ctx := context.Background()
	_, err = provider.Generate(ctx, "Write a hello world function", nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestRemoteProvider_Generate_EmptyResult(t *testing.T) {
	// Create a mock HTTP server that returns empty result
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := RemoteGenerateResponse{
			Result: "",
			Model:  "test-model",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create provider
	provider, err := NewRemoteProvider("test-model", server.URL, "test-api-key")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Test generation
	ctx := context.Background()
	_, err = provider.Generate(ctx, "Write a hello world function", nil)

	if err == nil {
		t.Error("Expected error for empty result, got nil")
	}
}

func TestNewRemoteProvider_NoEndpoint(t *testing.T) {
	_, err := NewRemoteProvider("test", "", "key")
	if err == nil {
		t.Error("Expected error when endpoint is empty, got nil")
	}
}
