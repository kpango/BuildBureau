# Multi-Provider LLM Support

BuildBureau now supports multiple LLM providers natively through their official
SDKs:

## Supported Providers

### ✅ Native Integrations (Official SDKs)

| Provider             | SDK                                      | Models                         | Status        |
| -------------------- | ---------------------------------------- | ------------------------------ | ------------- |
| **Google Gemini**    | `google.golang.org/genai`                | Gemini 2.0 Flash, Pro          | ✅ Production |
| **OpenAI**           | `github.com/sashabaranov/go-openai`      | GPT-4, GPT-3.5, etc.           | ✅ Production |
| **Anthropic Claude** | `github.com/liushuangls/go-anthropic/v2` | Claude 3.5 Sonnet, Opus, Haiku | ✅ Production |

### 🔌 Remote Agent API

For providers without Go SDKs, use the Remote Agent API (HTTP/gRPC):

- Custom models
- Cohere, Mistral, Llama
- Self-hosted models
- Internal APIs

---

## Quick Start

### 1. Set API Keys

```bash
export GEMINI_API_KEY="your-gemini-key"
export OPENAI_API_KEY="your-openai-key"
export CLAUDE_API_KEY="your-claude-key"
```

### 2. Update Configuration

`config.yaml`:

```yaml
llms:
  default_model: gemini # or openai, claude
  api_keys:
    gemini: { env: GEMINI_API_KEY }
    openai: { env: OPENAI_API_KEY }
    claude: { env: CLAUDE_API_KEY }
```

### 3. Run BuildBureau

```bash
./buildbureau
```

The system will automatically initialize all providers with valid API keys!

---

## Provider Details

### 🔵 Google Gemini

**Default Model:** `gemini-2.0-flash-exp`

**Features:**

- Fast and cost-effective
- Large context window
- Native multimodal support

**Usage:**

```go
provider, _ := llm.NewGeminiProvider(apiKey)
response, _ := provider.Generate(ctx, "Your prompt", &llm.GenerateOptions{
    Temperature: 0.7,
    MaxTokens: 2048,
    SystemPrompt: "You are a helpful assistant.",
})
```

**Get API Key:** https://aistudio.google.com/app/apikey

---

### 🟢 OpenAI

**Default Model:** `gpt-4-turbo-preview`

**Supported Models:**

- `gpt-4-turbo-preview` - Most capable
- `gpt-4` - Stable GPT-4
- `gpt-3.5-turbo` - Fast and affordable

**Usage:**

```go
provider, _ := llm.NewOpenAIProvider(apiKey, "gpt-4-turbo-preview")
response, _ := provider.Generate(ctx, "Your prompt", &llm.GenerateOptions{
    Temperature: 0.7,
    MaxTokens: 2048,
    SystemPrompt: "You are a helpful assistant.",
})
```

**Model Override:**

```bash
export OPENAI_MODEL="gpt-3.5-turbo"
```

**Get API Key:** https://platform.openai.com/api-keys

---

### 🟣 Anthropic Claude

**Default Model:** `claude-3-5-sonnet-20241022`

**Supported Models:**

- `claude-3-5-sonnet-20241022` - Best overall (default)
- `claude-3-opus-20240229` - Most capable
- `claude-3-haiku-20240307` - Fast and affordable

**Usage:**

```go
provider, _ := llm.NewClaudeProvider(apiKey, "claude-3-5-sonnet-20241022")
response, _ := provider.Generate(ctx, "Your prompt", &llm.GenerateOptions{
    Temperature: 0.7,
    MaxTokens: 2048,
    SystemPrompt: "You are a helpful assistant.",
})
```

**Model Override:**

```bash
export CLAUDE_MODEL="claude-3-haiku-20240307"
```

**Get API Key:** https://console.anthropic.com/

---

## Using Multiple Providers

BuildBureau's LLM Manager automatically initializes all providers with valid API
keys:

```go
// Manager automatically detects and initializes providers
manager, _ := llm.NewManager(config.LLMs)

// Use Gemini
response1, _ := manager.Generate(ctx, "gemini", "Write a function", opts)

// Use OpenAI
response2, _ := manager.Generate(ctx, "openai", "Write a function", opts)

// Use Claude
response3, _ := manager.Generate(ctx, "claude", "Write a function", opts)
```

---

## Example: Test All Providers

```bash
# Run the multi-provider example
go run examples/test_multiple_providers/main.go
```

This will test all configured providers and show their responses side-by-side.

---

## Provider Comparison

| Feature               | Gemini      | OpenAI      | Claude      |
| --------------------- | ----------- | ----------- | ----------- |
| **Speed**             | ⚡⚡⚡ Fast | ⚡⚡ Medium | ⚡⚡ Medium |
| **Cost**              | 💰 Low      | 💰💰 Medium | 💰💰 Medium |
| **Context Window**    | 32K-1M      | 128K        | 200K        |
| **Code Quality**      | ⭐⭐⭐      | ⭐⭐⭐⭐    | ⭐⭐⭐⭐    |
| **Creative Writing**  | ⭐⭐⭐      | ⭐⭐⭐⭐    | ⭐⭐⭐⭐⭐  |
| **Structured Output** | ⭐⭐⭐⭐    | ⭐⭐⭐⭐    | ⭐⭐⭐      |
| **Latency**           | ~1s         | ~2s         | ~2s         |

---

## Configuration Options

### Per-Provider Model Selection

Set specific models via environment variables:

```bash
# Use GPT-3.5 instead of GPT-4
export OPENAI_MODEL="gpt-3.5-turbo"

# Use Claude Haiku instead of Sonnet
export CLAUDE_MODEL="claude-3-haiku-20240307"

# Gemini model is currently fixed to 2.0 Flash
```

### Default Provider

Set which provider to use by default in `config.yaml`:

```yaml
llms:
  default_model: openai # gemini, openai, or claude
```

### Provider-Specific Settings

Each provider supports:

- **Temperature** (0.0-1.0): Creativity vs consistency
- **MaxTokens**: Response length limit
- **SystemPrompt**: Role/behavior instructions

---

## Remote Agent API (HTTP)

For providers without native SDK support, use the Remote Agent API:

```yaml
# config.yaml
llms:
  api_keys:
    custom: { env: CUSTOM_API_KEY }
```

```bash
export CUSTOM_LLM_ENDPOINT="http://localhost:8080"
export CUSTOM_API_KEY="your-key"
```

Create a simple HTTP server implementing:

**POST /v1/generate**

```json
{
  "prompt": "Your prompt",
  "temperature": 0.7,
  "max_tokens": 2048,
  "system_prompt": "System instructions"
}
```

Response:

```json
{
  "result": "Generated text",
  "model": "custom-model"
}
```

See `docs/REMOTE_AGENTS.md` for complete Remote Agent API documentation.

---

## Troubleshooting

### Provider Not Initialized

**Symptom:** "model X not available"

**Solution:**

1. Check API key is set: `echo $OPENAI_API_KEY`
2. Verify key in config: `api_keys.openai.env` matches environment variable name
3. Check key is valid (not "demo-key")

### API Errors

**Symptom:** "failed to create chat completion" or "failed to create message"

**Common Causes:**

- Invalid API key
- Insufficient credits/quota
- Rate limiting
- Network issues

**Solution:**

1. Verify API key is valid
2. Check account has credits
3. Add retry logic
4. Check provider status page

### Model Not Found

**Symptom:** "model does not exist"

**Solution:** Use correct model names:

- OpenAI: `gpt-4-turbo-preview`, `gpt-3.5-turbo`
- Claude: `claude-3-5-sonnet-20241022`, `claude-3-haiku-20240307`
- Gemini: Model is fixed (no override needed)

---

## Best Practices

### 1. Use Appropriate Models

- **Code Generation:** OpenAI GPT-4, Gemini 2.0
- **Analysis:** Claude 3.5 Sonnet, GPT-4
- **Speed/Cost:** Gemini Flash, GPT-3.5, Claude Haiku

### 2. Handle Errors Gracefully

```go
response, err := provider.Generate(ctx, prompt, opts)
if err != nil {
    // Fallback to another provider
    response, err = fallbackProvider.Generate(ctx, prompt, opts)
}
```

### 3. Set Appropriate Limits

```go
opts := &llm.GenerateOptions{
    Temperature: 0.7,  // 0=deterministic, 1=creative
    MaxTokens:   2048, // Limit response length
    SystemPrompt: "Your role",
}
```

### 4. Cache Responses

For repeated queries, cache LLM responses to:

- Reduce costs
- Improve latency
- Avoid rate limits

---

## Cost Optimization

| Provider | Model         | Cost (per 1M tokens)           |
| -------- | ------------- | ------------------------------ |
| Gemini   | 2.0 Flash     | $0.075 (input), $0.30 (output) |
| OpenAI   | GPT-3.5 Turbo | $0.50 (input), $1.50 (output)  |
| OpenAI   | GPT-4 Turbo   | $10 (input), $30 (output)      |
| Claude   | Haiku         | $0.25 (input), $1.25 (output)  |
| Claude   | Sonnet 3.5    | $3 (input), $15 (output)       |

**Tips:**

- Use Gemini Flash or GPT-3.5 for simple tasks
- Use Claude Haiku for fast responses
- Reserve GPT-4/Claude Opus for complex reasoning

---

## Migration Guide

### From RemoteProvider to Native

**Before:**

```yaml
llms:
  api_keys:
    openai: { env: OPENAI_API_KEY }

# Required separate HTTP server
export OPENAI_ENDPOINT="http://localhost:8080"
```

**After:**

```yaml
llms:
  api_keys:
    openai: { env: OPENAI_API_KEY }

# No endpoint needed - uses native SDK!
```

Native providers are:

- ✅ Faster (no HTTP overhead)
- ✅ More reliable (no separate server)
- ✅ Easier to configure
- ✅ Better error handling

---

## Further Reading

- [OpenAI API Documentation](https://platform.openai.com/docs/api-reference)
- [Claude API Documentation](https://docs.anthropic.com/claude/reference)
- [Gemini API Documentation](https://ai.google.dev/docs)
- [Remote Agents Guide](./REMOTE_AGENTS.md)
- [ADK Integration](./ADK_INTEGRATION.md)

---

## Support

For issues or questions:

1. Check this documentation
2. Review example code in `examples/`
3. Check provider status pages
4. Open a GitHub issue

**Pro Tip:** Start with Gemini (free tier) or GPT-3.5 (affordable) to test your
integration before using premium models!
