package adk

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"strings"
	"time"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"

	"buildbureau/pkg/a2a"
	"buildbureau/pkg/config"
)

type Agent[Req any, Resp any] struct {
	ID           string
	Role         string
	ModelID      string
	SystemPrompt string
	Bus          *a2a.Bus
	Config       config.AgentConfig

	// Real ADK Runner
	Runner *runner.Runner

	// Mock implementation for testing/demo without real LLM
	MockImpl func(ctx context.Context, req Req) (Resp, error)
}

func NewAgent[Req any, Resp any](
	id string,
	cfg config.AgentConfig,
	bus *a2a.Bus,
	apiKey string,
	modelName string, // Resolved real model name (e.g. gemini-1.5-pro)
) *Agent[Req, Resp] {
	a := &Agent[Req, Resp]{
		ID:           id,
		Role:         cfg.Role,
		ModelID:      modelName,
		SystemPrompt: cfg.SystemPrompt,
		Bus:          bus,
		Config:       cfg,
	}

	// Initialize Real ADK if key is present
	if apiKey != "" {
		ctx := context.Background()
		// Initialize Model
		// Since we are using the Gemini adapter from ADK, we must use Gemini models.
		// If modelName is not a Gemini model, this might fail at runtime.
		model, err := gemini.NewModel(ctx, modelName, &genai.ClientConfig{
			APIKey: apiKey,
		})
		if err != nil {
			// Log error but allow fallback to mock?
			// fmt.Printf("Failed to create model for agent %s: %v\n", id, err)
		} else {
			// Initialize LLM Agent
			llmAgent, err := llmagent.New(llmagent.Config{
				Name:        id,
				Model:       model,
				Description: fmt.Sprintf("Agent %s with role %s", id, cfg.Role),
				Instruction: cfg.SystemPrompt,
			})
			if err != nil {
				// fmt.Printf("Failed to create llmagent for %s: %v\n", id, err)
			} else {
				// Initialize Runner
				r, err := runner.New(runner.Config{
					Agent:          llmAgent,
					SessionService: session.InMemoryService(),
				})
				if err != nil {
					// fmt.Printf("Failed to create runner for %s: %v\n", id, err)
				} else {
					a.Runner = r
				}
			}
		}
	}

	return a
}

func (a *Agent[Req, Resp]) Process(ctx context.Context, req Req) (Resp, error) {
	var resp Resp

	// 1. Notify Start via A2A
	a.Bus.Send(ctx, a2a.Message{
		ID:        fmt.Sprintf("%s-%d", a.ID, time.Now().UnixNano()),
		From:      a.ID,
		To:        "LOG",
		Type:      "START",
		Payload:   fmt.Sprintf("Agent %s started processing", a.Role),
		Timestamp: time.Now(),
	})

	// 2. Check Mock Logic
	// If Runner is nil, we MUST use MockImpl.
	if a.Runner == nil {
		if a.MockImpl != nil {
			r, err := a.MockImpl(ctx, req)
			if err == nil {
				a.sendComplete(ctx)
				return r, nil
			}
			a.sendError(ctx, err)
			return resp, err
		}
		// No runner and no mock? Fail.
		err := fmt.Errorf("no execution engine (Runner or Mock) for agent %s", a.ID)
		a.sendError(ctx, err)
		return resp, err
	}

	// 3. Prepare Input
	reqBytes, err := json.Marshal(req)
	if err != nil {
		a.sendError(ctx, err)
		return resp, fmt.Errorf("failed to marshal request: %w", err)
	}
	userPrompt := string(reqBytes)

	// 4. Run ADK Agent
	sessionID := fmt.Sprintf("sess-%d", time.Now().UnixNano())
	userID := "user-default"

	msg := &genai.Content{
		Parts: []*genai.Part{
			{Text: userPrompt},
		},
	}

	var responseText string

	// Execute Runner
	// Note: runner.Run returns an iterator. We loop over it.
	next, stop := iter.Pull2(a.Runner.Run(ctx, userID, sessionID, msg, agent.RunConfig{}))
	defer stop()

	for {
		event, err, ok := next()
		if !ok {
			break
		}
		if err != nil {
			a.sendError(ctx, err)
			return resp, fmt.Errorf("runner execution failed: %w", err)
		}

		// Check for response content
		if event != nil && event.Content != nil {
			for _, part := range event.Content.Parts {
				if part.Text != "" {
					responseText += part.Text
				}
			}
		}
	}

	// 5. Unmarshal Response
	cleanedResp := cleanJSON(responseText)
	if cleanedResp == "" {
		err := fmt.Errorf("received empty response from agent")
		a.sendError(ctx, err)
		return resp, err
	}

	if err := json.Unmarshal([]byte(cleanedResp), &resp); err != nil {
		a.sendError(ctx, fmt.Errorf("unmarshal failed: %v. Raw: %s", err, responseText))
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

func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
