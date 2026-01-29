# Google ADK Integration Guide

## Overview

BuildBureau now includes full integration with Google's Generative AI (ADK - Agent Development Kit) using the official `google.golang.org/genai` SDK. This enables the system to use Google's Gemini models for intelligent agent responses.

## Features

- ✅ Synchronous text generation
- ✅ Streaming responses
- ✅ System instructions support
- ✅ Temperature control
- ✅ Token limit configuration
- ✅ Token usage tracking
- ✅ Multiple Gemini model support

## Setup

### 1. Get a Google AI API Key

1. Visit [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Click "Create API Key"
3. Copy the generated key

### 2. Set Environment Variable

```bash
export GOOGLE_AI_API_KEY="your-api-key-here"
```

Or add to your `.env` file:

```bash
GOOGLE_AI_API_KEY=your-api-key-here
```

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/kpango/BuildBureau/internal/llm"
)

func main() {
    ctx := context.Background()
    apiKey := os.Getenv("GOOGLE_AI_API_KEY")
    
    // Create client with default model (gemini-2.0-flash-exp)
    client, err := llm.NewGoogleADKClient(ctx, apiKey, "")
    if err != nil {
        log.Fatal(err)
    }
    
    // Generate response
    resp, err := client.Generate(ctx, &llm.Request{
        Messages: []llm.Message{
            {Role: "user", Content: "Hello, how are you?"},
        },
        Temperature: 0.7,
        MaxTokens:   100,
    })
    
    fmt.Println(resp.Content)
}
```

### With System Instructions

```go
req := &llm.Request{
    Messages: []llm.Message{
        {
            Role:    "system", 
            Content: "You are a helpful assistant specialized in software architecture.",
        },
        {
            Role:    "user",
            Content: "Explain microservices architecture.",
        },
    },
    Temperature: 0.7,
    MaxTokens:   500,
}

resp, err := client.Generate(ctx, req)
```

### Streaming Responses

```go
req := &llm.Request{
    Messages: []llm.Message{
        {Role: "user", Content: "Write a short story."},
    },
    Temperature: 0.8,
    MaxTokens:   500,
}

contentCh, errCh := client.StreamGenerate(ctx, req)

// Read chunks as they arrive
for chunk := range contentCh {
    fmt.Print(chunk)
}

// Check for errors
if err := <-errCh; err != nil {
    log.Printf("Streaming error: %v", err)
}
```

### Using the Client Factory

```go
factory := llm.NewClientFactory("google", apiKey, "gemini-2.0-flash-exp")
client, err := factory.Create()
if err != nil {
    log.Fatal(err)
}

// Use client...
```

## Configuration

### Supported Models

- `gemini-2.0-flash-exp` (default) - Fast, efficient, multimodal
- `gemini-1.5-pro` - Advanced reasoning
- `gemini-1.5-flash` - Balanced speed and capability
- All other Gemini models supported by the API

### Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `Temperature` | float64 | Controls randomness (0.0-1.0) | 0.7 |
| `MaxTokens` | int | Maximum tokens in response | Model default |
| `Model` | string | Gemini model to use | gemini-2.0-flash-exp |

### Message Roles

- `system` - System instructions (handled as system instruction in API)
- `user` - User messages
- `assistant` - Assistant/model messages (for conversation history)

## Integration with BuildBureau

### In Configuration (config.yaml)

```yaml
llm:
  provider: google  # or "mock" for testing
  model: gemini-2.0-flash-exp
  temperature: 0.7
  maxTokens: 1000
```

### Environment Variables

The system automatically reads `GOOGLE_AI_API_KEY` from the environment when using the Google provider.

### In Agents

Agents can use the LLM client through the service layer:

```go
// In service initialization
llmClient := llm.NewMockClient(nil) // Default for testing

// Or with Google ADK if API key is available
if apiKey := os.Getenv("GOOGLE_AI_API_KEY"); apiKey != "" {
    llmClient, _ = llm.NewGoogleADKClient(ctx, apiKey, cfg.LLM.Model)
}

// Use in services
presidentService := grpc.NewPresidentService(pool, kb, registry, llmClient)
```

## Examples

### Running the Example

```bash
# Set your API key
export GOOGLE_AI_API_KEY="your-key-here"

# Run the example
cd examples/google-adk
go run main.go
```

### Example Output

```
=== Google ADK Integration Example ===

Example 1: Simple text generation
Response: A hierarchical multi-agent system breaks down complex tasks...
Tokens used: 87
Finish reason: STOP

Example 2: Streaming response
Streaming: 1
2
3
4
5

=== Example Complete ===
```

## Error Handling

### Common Errors

**API Key Not Set:**
```go
client, err := llm.NewGoogleADKClient(ctx, "", "")
// err: API key is required
```

**Invalid Model:**
```go
// The API will return an error if the model doesn't exist
resp, err := client.Generate(ctx, req)
// Check err for model-related issues
```

**Rate Limiting:**
```go
// Google's API may rate limit requests
// Implement exponential backoff for production use
```

### Best Practices

1. **Always check errors:**
   ```go
   if err != nil {
       log.Printf("Generation failed: %v", err)
       // Handle appropriately
   }
   ```

2. **Use context with timeout:**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   ```

3. **Handle streaming errors:**
   ```go
   contentCh, errCh := client.StreamGenerate(ctx, req)
   for chunk := range contentCh {
       // Process chunk
   }
   if err := <-errCh; err != nil {
       // Handle error
   }
   ```

## Testing

### Unit Tests

```bash
# Run tests without API key (skips integration tests)
go test ./internal/llm -v

# Run with API key for integration tests
export GOOGLE_AI_API_KEY="your-key"
go test ./internal/llm -v
```

### Mock Client for Testing

```go
// Use mock client in tests
mockClient := llm.NewMockClient([]string{
    "Mock response 1",
    "Mock response 2",
})

resp, err := mockClient.Generate(ctx, req)
// resp.Content == "Mock response 1"
```

## Performance Considerations

### Token Usage

- Monitor `resp.TokensUsed` to track API costs
- Set appropriate `MaxTokens` limits
- Use `temperature` to balance creativity vs consistency

### Streaming vs Synchronous

- **Use streaming** for:
  - Long responses
  - Real-time user feedback
  - Interactive applications

- **Use synchronous** for:
  - Short responses
  - Batch processing
  - Simpler code flow

### Cost Optimization

```go
// Use lower temperature for more deterministic responses
req.Temperature = 0.1

// Limit tokens for shorter responses
req.MaxTokens = 100

// Use appropriate model for task
// - gemini-2.0-flash-exp: Fast, cost-effective
// - gemini-1.5-pro: Complex reasoning tasks
```

## Troubleshooting

### Problem: "API key is required"
**Solution:** Set `GOOGLE_AI_API_KEY` environment variable

### Problem: "failed to create genai client"
**Solution:** Check network connectivity and API key validity

### Problem: "model not found"
**Solution:** Verify model name is correct and available

### Problem: Slow responses
**Solution:** 
- Use faster models (gemini-2.0-flash-exp)
- Reduce MaxTokens
- Consider streaming for better UX

## Additional Resources

- [Google AI Studio](https://aistudio.google.com/)
- [Gemini API Documentation](https://ai.google.dev/)
- [Go genai SDK](https://pkg.go.dev/google.golang.org/genai)
- [BuildBureau Examples](../examples/)

## Support

For issues or questions:
1. Check this documentation
2. Review example code in `examples/google-adk/`
3. Check test cases in `internal/llm/client_test.go`
4. Open an issue on GitHub
