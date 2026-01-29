package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/genai"
)

// GeminiProvider implements the Provider interface for Google Gemini using the genai library.
type GeminiProvider struct {
	client *genai.Client
	model  string
}

// NewGeminiProvider creates a new Gemini provider with real API integration.
func NewGeminiProvider(apiKey string) (*GeminiProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("gemini API key is required")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiProvider{
		client: client,
		model:  "gemini-2.0-flash-exp", // Using latest flash model
	}, nil
}

// Generate sends a prompt to Gemini and returns the response.
func (p *GeminiProvider) Generate(ctx context.Context, prompt string, opts *GenerateOptions) (string, error) {
	if opts == nil {
		opts = &GenerateOptions{
			Temperature: 0.7,
			MaxTokens:   2048,
		}
	}

	// Create content with text prompt
	userContent := &genai.Content{
		Parts: []*genai.Part{{Text: prompt}},
		Role:  genai.RoleUser,
	}

	// Create config
	temp := float32(opts.Temperature)
	maxTokens := int32(opts.MaxTokens)
	config := &genai.GenerateContentConfig{
		Temperature:     &temp,
		MaxOutputTokens: maxTokens,
	}

	// Add system instruction if provided
	if opts.SystemPrompt != "" {
		config.SystemInstruction = &genai.Content{
			Parts: []*genai.Part{{Text: opts.SystemPrompt}},
		}
	}

	// Generate content
	resp, err := p.client.Models.GenerateContent(ctx, p.model, []*genai.Content{userContent}, config)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract text from response
	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	var responseText strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			responseText.WriteString(part.Text)
		}
	}

	if responseText.Len() == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	return responseText.String(), nil
}

// Name returns the provider name.
func (p *GeminiProvider) Name() string {
	return "gemini"
}

// Close closes the Gemini client.
func (p *GeminiProvider) Close() error {
	// The genai.Client doesn't have a Close method in the current version
	// The client will be garbage collected
	return nil
}

// RemoteProvider implements the Provider interface for remote LLM services
// This is used for Claude, Codex, and Qwen via the Remote Agent API.
type RemoteProvider struct {
	httpClient *http.Client
	name       string
	endpoint   string
	apiKey     string
}

// RemoteGenerateRequest represents the request to a remote LLM service.
type RemoteGenerateRequest struct {
	Prompt       string  `json:"prompt"`
	Model        string  `json:"model,omitempty"`
	SystemPrompt string  `json:"system_prompt,omitempty"`
	Temperature  float64 `json:"temperature,omitempty"`
	MaxTokens    int     `json:"max_tokens,omitempty"`
}

// RemoteGenerateResponse represents the response from a remote LLM service.
type RemoteGenerateResponse struct {
	Result string         `json:"result"`
	Model  string         `json:"model,omitempty"`
	Usage  map[string]any `json:"usage,omitempty"`
	Error  string         `json:"error,omitempty"`
}

// NewRemoteProvider creates a new remote provider.
func NewRemoteProvider(name, endpoint, apiKey string) (*RemoteProvider, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint is required for remote provider")
	}

	return &RemoteProvider{
		name:     name,
		endpoint: endpoint,
		apiKey:   apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// Generate sends a prompt to the remote provider via HTTP.
func (p *RemoteProvider) Generate(ctx context.Context, prompt string, opts *GenerateOptions) (string, error) {
	if opts == nil {
		opts = &GenerateOptions{
			Temperature: 0.7,
			MaxTokens:   2048,
		}
	}

	// Create request body
	reqBody := RemoteGenerateRequest{
		Prompt:       prompt,
		Model:        p.name,
		Temperature:  opts.Temperature,
		MaxTokens:    opts.MaxTokens,
		SystemPrompt: opts.SystemPrompt,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint+"/v1/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	// Send request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("remote provider returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result RemoteGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error != "" {
		return "", fmt.Errorf("remote provider error: %s", result.Error)
	}

	if result.Result == "" {
		return "", fmt.Errorf("empty result from remote provider")
	}

	return result.Result, nil
}

// Name returns the provider name.
func (p *RemoteProvider) Name() string {
	return p.name
}

// Close closes the HTTP client.
func (p *RemoteProvider) Close() error {
	p.httpClient.CloseIdleConnections()
	return nil
}

// OpenAIProvider implements the Provider interface for OpenAI using the official SDK.
type OpenAIProvider struct {
	client *openai.Client
	model  string
}

// NewOpenAIProvider creates a new OpenAI provider with real API integration.
func NewOpenAIProvider(apiKey string, model string) (*OpenAIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	client := openai.NewClient(apiKey)

	// Default to GPT-4 if no model specified
	if model == "" {
		model = openai.GPT4TurboPreview
	}

	return &OpenAIProvider{
		client: client,
		model:  model,
	}, nil
}

// Generate sends a prompt to OpenAI and returns the response.
func (p *OpenAIProvider) Generate(ctx context.Context, prompt string, opts *GenerateOptions) (string, error) {
	if opts == nil {
		opts = &GenerateOptions{
			Temperature: 0.7,
			MaxTokens:   2048,
		}
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	// Add system message if provided
	if opts.SystemPrompt != "" {
		messages = append([]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: opts.SystemPrompt,
			},
		}, messages...)
	}

	req := openai.ChatCompletionRequest{
		Model:       p.model,
		Messages:    messages,
		Temperature: float32(opts.Temperature),
		MaxTokens:   opts.MaxTokens,
	}

	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return resp.Choices[0].Message.Content, nil
}

// Name returns the provider name.
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// Close closes the OpenAI client.
func (p *OpenAIProvider) Close() error {
	// The OpenAI client doesn't require explicit cleanup
	return nil
}

// ClaudeProvider implements the Provider interface for Anthropic Claude using the official SDK.
type ClaudeProvider struct {
	client *anthropic.Client
	model  string
}

// NewClaudeProvider creates a new Claude provider with real API integration.
func NewClaudeProvider(apiKey string, model string) (*ClaudeProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("claude API key is required")
	}

	client := anthropic.NewClient(apiKey)

	// Default to Claude 3.5 Sonnet if no model specified
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}

	return &ClaudeProvider{
		client: client,
		model:  model,
	}, nil
}

// Generate sends a prompt to Claude and returns the response.
func (p *ClaudeProvider) Generate(ctx context.Context, prompt string, opts *GenerateOptions) (string, error) {
	if opts == nil {
		opts = &GenerateOptions{
			Temperature: 0.7,
			MaxTokens:   2048,
		}
	}

	req := anthropic.MessagesRequest{
		Model:       anthropic.Model(p.model),
		MaxTokens:   opts.MaxTokens,
		Temperature: new(float32(opts.Temperature)),
		Messages: []anthropic.Message{
			{
				Role: anthropic.RoleUser,
				Content: []anthropic.MessageContent{
					anthropic.NewTextMessageContent(prompt),
				},
			},
		},
	}

	// Add system message if provided
	if opts.SystemPrompt != "" {
		req.System = opts.SystemPrompt
	}

	resp, err := p.client.CreateMessages(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create message: %w", err)
	}

	if len(resp.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	// Extract text from content blocks
	var result strings.Builder
	for _, content := range resp.Content {
		if content.Type == "text" && content.Text != nil {
			result.WriteString(*content.Text)
		}
	}

	if result.Len() == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	return result.String(), nil
}

// Name returns the provider name.
func (p *ClaudeProvider) Name() string {
	return "claude"
}

// Close closes the Claude client.
func (p *ClaudeProvider) Close() error {
	// The Claude client doesn't require explicit cleanup
	return nil
}

// Helper function to convert float32 to *float32.
//
//go:fix inline
func toFloat32Ptr(f float32) *float32 {
	return new(f)
}
