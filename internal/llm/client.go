package llm

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
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

// GoogleADKClient implements Google Generative AI integration
type GoogleADKClient struct {
	client *genai.Client
	model  string
}

// NewGoogleADKClient creates a new Google ADK client
func NewGoogleADKClient(ctx context.Context, apiKey, model string) (*GoogleADKClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	
	// Use default model if not specified
	if model == "" {
		model = "gemini-2.0-flash-exp"
	}
	
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}
	
	return &GoogleADKClient{
		client: client,
		model:  model,
	}, nil
}

// Generate generates a response using Google Generative AI
func (c *GoogleADKClient) Generate(ctx context.Context, req *Request) (*Response, error) {
	if c.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}
	
	// Convert our messages to genai format
	var contents []*genai.Content
	var systemInstruction *genai.Content
	
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			// System messages are handled separately as system instruction
			systemInstruction = genai.NewContentFromText(msg.Content, genai.RoleUser)
		} else {
			// Convert role
			role := genai.RoleUser
			if msg.Role == "assistant" {
				role = genai.RoleModel
			}
			contents = append(contents, genai.NewContentFromText(msg.Content, genai.Role(role)))
		}
	}
	
	// Configure generation
	genConfig := &genai.GenerateContentConfig{}
	if req.Temperature > 0 {
		temp := float32(req.Temperature)
		genConfig.Temperature = &temp
	}
	if req.MaxTokens > 0 {
		genConfig.MaxOutputTokens = int32(req.MaxTokens)
	}
	if systemInstruction != nil {
		genConfig.SystemInstruction = systemInstruction
	}
	
	// Generate content
	resp, err := c.client.Models.GenerateContent(ctx, c.model, contents, genConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}
	
	// Parse response
	if resp == nil || len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response candidates returned")
	}
	
	candidate := resp.Candidates[0]
	if candidate.Content == nil {
		return nil, fmt.Errorf("no content in response")
	}
	
	// Extract text from parts
	var contentBuilder strings.Builder
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			contentBuilder.WriteString(part.Text)
		}
	}
	
	// Determine finish reason
	finishReason := "stop"
	if candidate.FinishReason != "" {
		finishReason = string(candidate.FinishReason)
	}
	
	// Calculate tokens used
	tokensUsed := 0
	if resp.UsageMetadata != nil {
		tokensUsed = int(resp.UsageMetadata.TotalTokenCount)
	}
	
	return &Response{
		Content:      contentBuilder.String(),
		FinishReason: finishReason,
		TokensUsed:   tokensUsed,
	}, nil
}

// StreamGenerate generates a streaming response using Google Generative AI
func (c *GoogleADKClient) StreamGenerate(ctx context.Context, req *Request) (<-chan string, <-chan error) {
	contentCh := make(chan string, 10)
	errCh := make(chan error, 1)
	
	go func() {
		defer close(contentCh)
		defer close(errCh)
		
		if c.client == nil {
			errCh <- fmt.Errorf("client not initialized")
			return
		}
		
		// Convert our messages to genai format
		var contents []*genai.Content
		var systemInstruction *genai.Content
		
		for _, msg := range req.Messages {
			if msg.Role == "system" {
				systemInstruction = genai.NewContentFromText(msg.Content, genai.RoleUser)
			} else {
				role := genai.RoleUser
				if msg.Role == "assistant" {
					role = genai.RoleModel
				}
				contents = append(contents, genai.NewContentFromText(msg.Content, genai.Role(role)))
			}
		}
		
		// Configure generation
		genConfig := &genai.GenerateContentConfig{}
		if req.Temperature > 0 {
			temp := float32(req.Temperature)
			genConfig.Temperature = &temp
		}
		if req.MaxTokens > 0 {
			genConfig.MaxOutputTokens = int32(req.MaxTokens)
		}
		if systemInstruction != nil {
			genConfig.SystemInstruction = systemInstruction
		}
		
		// Stream content - returns an iterator
		iter := c.client.Models.GenerateContentStream(ctx, c.model, contents, genConfig)
		
		// Iterate over the stream
		for resp, err := range iter {
			if err != nil {
				errCh <- fmt.Errorf("stream error: %w", err)
				return
			}
			
			if resp == nil || len(resp.Candidates) == 0 {
				continue
			}
			
			candidate := resp.Candidates[0]
			if candidate.Content == nil {
				continue
			}
			
			// Extract and send text from parts
			for _, part := range candidate.Content.Parts {
				if part.Text != "" {
					contentCh <- part.Text
				}
			}
		}
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
	return f.CreateWithContext(context.Background())
}

// CreateWithContext creates an appropriate LLM client with context
func (f *ClientFactory) CreateWithContext(ctx context.Context) (Client, error) {
	switch f.provider {
	case "google":
		return NewGoogleADKClient(ctx, f.apiKey, f.model)
	case "mock":
		return NewMockClient(nil), nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", f.provider)
	}
}
