package adk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"buildbureau/pkg/a2a"
	"buildbureau/pkg/config"
)

type Agent[Req any, Resp any] struct {
	ID           string
	Role         string
	ModelID      string
	SystemPrompt string
	Bus          *a2a.Bus
	LLM          LLMClient
	Config       config.AgentConfig

	// Mock implementation for testing/demo without real LLM
	MockImpl func(ctx context.Context, req Req) (Resp, error)
}

func NewAgent[Req any, Resp any](
	id string,
	cfg config.AgentConfig,
	bus *a2a.Bus,
	llm LLMClient,
) *Agent[Req, Resp] {
	return &Agent[Req, Resp]{
		ID:           id,
		Role:         cfg.Role,
		ModelID:      cfg.Model,
		SystemPrompt: cfg.SystemPrompt,
		Bus:          bus,
		LLM:          llm,
		Config:       cfg,
	}
}

func (a *Agent[Req, Resp]) Process(ctx context.Context, req Req) (Resp, error) {
	var resp Resp

	// 1. Notify Start via A2A
	a.Bus.Send(ctx, a2a.Message{
		ID:        fmt.Sprintf("%s-%d", a.ID, time.Now().UnixNano()),
		From:      a.ID,
		To:        "LOG", // Broadcast/Log
		Type:      "START",
		Payload:   fmt.Sprintf("Agent %s started processing", a.Role),
		Timestamp: time.Now(),
	})

	// 2. Check if MockImpl is set (for deterministic testing/demo)
	if a.MockImpl != nil {
		r, err := a.MockImpl(ctx, req)
		if err == nil {
			a.sendComplete(ctx)
			return r, nil
		}
		// If mock fails or returns error, maybe fall through? No, return error.
		a.sendError(ctx, err)
		return resp, err
	}

	// 3. Prepare Prompt
	reqBytes, err := json.Marshal(req)
	if err != nil {
		a.sendError(ctx, err)
		return resp, fmt.Errorf("failed to marshal request: %w", err)
	}
	userPrompt := string(reqBytes)

	// 4. Call LLM
	llmResp, err := a.LLM.Generate(ctx, a.SystemPrompt, userPrompt, a.ModelID)
	if err != nil {
		a.sendError(ctx, err)
		return resp, fmt.Errorf("llm generation failed: %w", err)
	}

	// 5. Unmarshal Response
	// The LLM is expected to return JSON. We might need to sanitize the string (remove markdown code blocks).
	cleanedResp := cleanJSON(llmResp)

	if err := json.Unmarshal([]byte(cleanedResp), &resp); err != nil {
		// If strict parsing fails, maybe we try to treat the whole text as a field?
		// But for now, let's assume strict JSON.
		// If it's a mock string, this will fail.
		a.sendError(ctx, fmt.Errorf("unmarshal failed: %v. Raw: %s", err, llmResp))
		return resp, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	// 6. Notify Success
	a.sendComplete(ctx)

	return resp, nil
}

func (a *Agent[Req, Resp]) sendComplete(ctx context.Context) {
	a.Bus.Send(ctx, a2a.Message{
		ID:        fmt.Sprintf("%s-%d", a.ID, time.Now().UnixNano()),
		From:      a.ID,
		To:        "LOG",
		Type:      "COMPLETE",
		Payload:   fmt.Sprintf("Agent %s completed task", a.Role),
		Timestamp: time.Now(),
	})
}

func (a *Agent[Req, Resp]) sendError(ctx context.Context, err error) {
	a.Bus.Send(ctx, a2a.Message{
		ID:        fmt.Sprintf("%s-%d", a.ID, time.Now().UnixNano()),
		From:      a.ID,
		To:        "ERROR",
		Type:      "ERROR",
		Payload:   err.Error(),
		Timestamp: time.Now(),
	})
}

// cleanJSON helper to strip ```json ... ``` if present
func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
