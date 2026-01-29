package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// GenerateRequest matches the RemoteGenerateRequest structure.
type GenerateRequest struct {
	Prompt       string  `json:"prompt"`
	Model        string  `json:"model,omitempty"`
	SystemPrompt string  `json:"system_prompt,omitempty"`
	Temperature  float64 `json:"temperature,omitempty"`
	MaxTokens    int     `json:"max_tokens,omitempty"`
}

// GenerateResponse matches the RemoteGenerateResponse structure.
type GenerateResponse struct {
	Result string         `json:"result"`
	Model  string         `json:"model,omitempty"`
	Usage  map[string]any `json:"usage,omitempty"`
	Error  string         `json:"error,omitempty"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	modelName := os.Getenv("MODEL_NAME")
	if modelName == "" {
		modelName = "example-model"
	}

	http.HandleFunc("/v1/generate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Check authorization header (optional)
		authHeader := r.Header.Get("Authorization")
		log.Printf("Received request with auth: %s", authHeader)

		// Decode request
		var req GenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
			return
		}

		log.Printf("Generating for prompt: %s (temp: %.2f, max_tokens: %d)",
			req.Prompt, req.Temperature, req.MaxTokens)

		// Simulate generation (in a real implementation, this would call an actual LLM)
		result := fmt.Sprintf("Generated response for: %s\n\nThis is a simulated response from the %s model. "+
			"In a real implementation, this would call the actual LLM API (Claude, Codex, Qwen, etc.).",
			req.Prompt, modelName)

		// Return response
		resp := GenerateResponse{
			Result: result,
			Model:  modelName,
			Usage: map[string]any{
				"prompt_tokens":     len(req.Prompt) / 4, // Rough estimate
				"completion_tokens": len(result) / 4,
				"total_tokens":      (len(req.Prompt) + len(result)) / 4,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/v1/status", func(w http.ResponseWriter, r *http.Request) {
		status := map[string]any{
			"status":       "ready",
			"model":        modelName,
			"capabilities": []string{"text-generation"},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Remote Agent API Server - %s\n\nEndpoints:\n- POST /v1/generate\n- GET /v1/status\n", modelName)
	})

	log.Printf("Starting Remote Agent API server for model '%s' on port %s", modelName, port)
	log.Printf("Endpoints:")
	log.Printf("  - POST http://localhost:%s/v1/generate", port)
	log.Printf("  - GET  http://localhost:%s/v1/status", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
