# Using BuildBureau with a Single LLM Provider

BuildBureau now supports running with just **one LLM provider** instead of
requiring all providers to be configured!

## Quick Start

You only need **ONE** API key to get started:

```bash
# Option 1: Use Gemini (Recommended - Free tier available)
export GEMINI_API_KEY="your-gemini-api-key"

# Option 2: Use OpenAI
export OPENAI_API_KEY="your-openai-api-key"

# Option 3: Use Claude
export CLAUDE_API_KEY="your-claude-api-key"

# Run BuildBureau
./build/buildbureau
```

## What You'll See

### With One Provider (Success)

```
Warning: environment variable OPENAI_API_KEY (for openai provider) is not set - this provider will be unavailable
Warning: environment variable CLAUDE_API_KEY (for claude provider) is not set - this provider will be unavailable
✓ Configuration loaded successfully with 1 provider(s) available
✓ Using default provider: gemini
```

### With No Providers (Error)

```
Error: no LLM provider API keys are set - at least one is required (GEMINI_API_KEY, OPENAI_API_KEY, CLAUDE_API_KEY, etc.)
```

## Supported Providers

| Provider   | Environment Variable | Notes                            |
| ---------- | -------------------- | -------------------------------- |
| **Gemini** | `GEMINI_API_KEY`     | ✅ Recommended (free tier, fast) |
| **OpenAI** | `OPENAI_API_KEY`     | GPT-4, GPT-3.5 models            |
| **Claude** | `CLAUDE_API_KEY`     | Claude 3.5 Sonnet, Opus, Haiku   |
| **Codex**  | `CODEX_API_KEY`      | Custom endpoint required         |
| **Qwen**   | `QWEN_API_KEY`       | Custom endpoint required         |

## Getting API Keys

### Google Gemini (Free Tier Available)

1. Visit [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Create an API key
3. Set environment variable:
   ```bash
   export GEMINI_API_KEY="your-key-here"
   ```

### OpenAI

1. Visit [OpenAI Platform](https://platform.openai.com/api-keys)
2. Create an API key
3. Set environment variable:
   ```bash
   export OPENAI_API_KEY="your-key-here"
   ```

### Anthropic Claude

1. Visit [Anthropic Console](https://console.anthropic.com/)
2. Create an API key
3. Set environment variable:
   ```bash
   export CLAUDE_API_KEY="your-key-here"
   ```

## Configuration

Your `config.yaml` can list all providers:

```yaml
llms:
  default_model: gemini
  api_keys:
    gemini: { env: GEMINI_API_KEY }
    openai: { env: OPENAI_API_KEY }
    claude: { env: CLAUDE_API_KEY }
```

**Only the providers with set API keys will be initialized.** Missing keys will
show warnings but won't stop BuildBureau from running.

## Using Multiple Providers (Optional)

You can set multiple API keys to enable provider fallback:

```bash
export GEMINI_API_KEY="key1"
export OPENAI_API_KEY="key2"
export CLAUDE_API_KEY="key3"
```

Benefits:

- Automatic fallback if one provider fails
- Choose different providers for different tasks
- Cost optimization by routing to cheaper providers

## Docker Usage

With Docker, you can pass just one API key:

```bash
# Docker run with single provider
docker run -e GEMINI_API_KEY="your-key" buildbureau

# Docker Compose
docker-compose run -e GEMINI_API_KEY="your-key" buildbureau
```

## Examples

### Example 1: Development with Gemini Only

```bash
export GEMINI_API_KEY="AIza..."
make build
./build/buildbureau
```

### Example 2: Production with OpenAI

```bash
export OPENAI_API_KEY="sk-..."
make docker-build
make docker-run
```

### Example 3: Testing with Claude

```bash
export CLAUDE_API_KEY="sk-ant..."
make test/llm-integration
```

## Troubleshooting

### "No LLM provider API keys are set"

**Problem**: No API keys configured **Solution**: Set at least one API key:

```bash
export GEMINI_API_KEY="your-key"
```

### Provider Warnings

**Problem**: Seeing warnings about missing providers **Solution**: This is
normal! Warnings inform you which providers are unavailable. BuildBureau will
work fine with the available provider(s).

### "Model not available"

**Problem**: Trying to use a provider that isn't configured **Solution**:
Either:

1. Set the API key for that provider, or
2. Change `default_model` in config.yaml to an available provider

## Best Practices

1. **Start Simple**: Begin with just Gemini (free tier)
2. **Add As Needed**: Add more providers when you need specific capabilities
3. **Monitor Costs**: Each provider has different pricing
4. **Use Fallbacks**: Multiple providers increase reliability
5. **Secure Keys**: Never commit API keys to git

## Migration from Previous Version

If you were using the old version that required all keys:

**Before** (Old - Required All):

```bash
export GEMINI_API_KEY="key1"
export OPENAI_API_KEY="key2"  # Required even if not used
export CLAUDE_API_KEY="key3"  # Required even if not used
```

**After** (New - Flexible):

```bash
export GEMINI_API_KEY="key1"  # Just one is enough!
```

Your existing configurations with multiple keys will continue to work!

## Summary

✅ **One API key is enough** to run BuildBureau ✅ **Multiple providers are
optional** but supported ✅ **Clear warnings** show which providers are
available ✅ **Backward compatible** with existing setups ✅ **Cost
effective** - pay only for what you use

Happy building! 🎉
