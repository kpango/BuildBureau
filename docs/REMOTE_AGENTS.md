# Remote Agent API Guide

This guide explains how to set up and use Remote Agents with BuildBureau for
integrating non-native LLM providers like Claude, Codex, and Qwen.

## Overview

BuildBureau natively supports Gemini through the ADK-go library. For other LLM
providers, we use the Remote Agent API pattern, which allows agents to
communicate with external LLM services via HTTP/gRPC endpoints.

## Architecture

```
┌─────────────────┐
│ BuildBureau     │
│   Engineer      │
│                 │
│  ┌──────────┐   │      ┌──────────────────┐
│  │ Sub-Agent│───┼─────▶│ Remote Agent API │
│  │  Proxy   │   │      │   (Claude)       │
│  └──────────┘   │      └──────────────────┘
│                 │
│  ┌──────────┐   │      ┌──────────────────┐
│  │ Sub-Agent│───┼─────▶│ Remote Agent API │
│  │  Proxy   │   │      │   (Codex)        │
│  └──────────┘   │      └──────────────────┘
└─────────────────┘
```

## Setting Up Remote Agents

### 1. Configure Agent with Remote Sub-Agents

Edit your agent configuration file (e.g., `agents/engineer.yaml`):

```yaml
name: Engineer
role: Engineer
description: Engineer with remote LLM capabilities
model: gemini
system_prompt: |
  You are an Engineer with access to specialized AI assistants...

sub_agents:
  - name: CodexWorker
    remote:
      endpoint: "http://localhost:8081"
      capabilities:
        - code-generation
        - python
        - javascript

  - name: ClaudeWorker
    remote:
      endpoint: "http://localhost:8082"
      capabilities:
        - analysis
        - documentation
```

### 2. Example Python Remote Agent Service

Create `claude_worker.py`:

```python
from flask import Flask, request, jsonify
import anthropic
import os

app = Flask(__name__)
client = anthropic.Anthropic(api_key=os.environ.get("CLAUDE_API_KEY"))

@app.route('/v1/generate', methods=['POST'])
def generate():
    data = request.json

    message = client.messages.create(
        model="claude-3-sonnet-20240229",
        max_tokens=data.get('max_tokens', 1024),
        temperature=data.get('temperature', 0.7),
        system=data.get('system_prompt', ''),
        messages=[
            {"role": "user", "content": data['prompt']}
        ]
    )

    return jsonify({
        "result": message.content[0].text,
        "model": "claude-3-sonnet",
        "usage": {
            "prompt_tokens": message.usage.input_tokens,
            "completion_tokens": message.usage.output_tokens
        }
    })

@app.route('/v1/status', methods=['GET'])
def status():
    return jsonify({
        "status": "ready",
        "model": "claude-3-sonnet",
        "capabilities": ["analysis", "documentation", "code-review"]
    })

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8082)
```

### 3. Example Go Remote Agent Service

Create `codex_worker.go`:

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
)

type GenerateRequest struct {
    Prompt       string  `json:"prompt"`
    Model        string  `json:"model"`
    Temperature  float64 `json:"temperature"`
    MaxTokens    int     `json:"max_tokens"`
    SystemPrompt string  `json:"system_prompt"`
}

type GenerateResponse struct {
    Result string `json:"result"`
    Model  string `json:"model"`
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
    var req GenerateRequest
    json.NewDecoder(r.Body).Decode(&req)

    // Call OpenAI Codex API here
    result := "Generated code based on: " + req.Prompt

    resp := GenerateResponse{
        Result: result,
        Model:  "codex",
    }

    json.NewEncoder(w).Encode(resp)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
    status := map[string]interface{}{
        "status":       "ready",
        "model":        "codex",
        "capabilities": []string{"code-generation"},
    }
    json.NewEncoder(w).Encode(status)
}

func main() {
    http.HandleFunc("/v1/generate", generateHandler)
    http.HandleFunc("/v1/status", statusHandler)
    log.Fatal(http.ListenAndServe(":8081", nil))
}
```

## Running the System

1. Start remote agent services:

```bash
python claude_worker.py  # Terminal 1
go run codex_worker.go   # Terminal 2
```

2. Start BuildBureau:

```bash
./buildbureau
```

## API Specification

### POST /v1/generate

Generate text using the LLM.

**Request:**

```json
{
  "prompt": "Write a Python function",
  "model": "claude-3",
  "temperature": 0.7,
  "max_tokens": 1000,
  "system_prompt": "You are a code assistant"
}
```

**Response:**

```json
{
  "result": "def example():\n    pass",
  "model": "claude-3",
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 5
  }
}
```

### GET /v1/status

Check service status.

**Response:**

```json
{
  "status": "ready",
  "model": "claude-3",
  "capabilities": ["analysis", "documentation"]
}
```

## Benefits

- **Separation of Concerns**: LLM code separate from core system
- **Language Flexibility**: Agents in any language
- **Scalability**: Distributed deployment
- **Hot Swapping**: Update without restart

## Security

- Use HTTPS in production
- Implement authentication
- Validate inputs
- Rate limiting
- Network isolation

See the full guide for more details on monitoring, troubleshooting, and advanced
configuration.
