// Copyright (C) 2024 BuildBureau team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestWrap(t *testing.T) {
	baseErr := errors.New("base error")
	wrapped := Wrap(baseErr, "wrapped")
	
	if wrapped == nil {
		t.Fatal("expected non-nil error")
	}
	
	if !strings.Contains(wrapped.Error(), "wrapped") {
		t.Errorf("expected error message to contain 'wrapped', got: %s", wrapped.Error())
	}
	
	if !strings.Contains(wrapped.Error(), "base error") {
		t.Errorf("expected error message to contain 'base error', got: %s", wrapped.Error())
	}
}

func TestWrapf(t *testing.T) {
	baseErr := errors.New("base error")
	wrapped := Wrapf(baseErr, "wrapped with %s", "format")
	
	if wrapped == nil {
		t.Fatal("expected non-nil error")
	}
	
	if !strings.Contains(wrapped.Error(), "wrapped with format") {
		t.Errorf("expected error message to contain 'wrapped with format', got: %s", wrapped.Error())
	}
}

func TestWrapNil(t *testing.T) {
	wrapped := Wrap(nil, "message")
	if wrapped != nil {
		t.Errorf("expected nil when wrapping nil error, got: %v", wrapped)
	}
}

func TestAgentErrors(t *testing.T) {
	tests := []struct {
		name     string
		errFunc  func() error
		contains string
	}{
		{
			name:     "AgentNotFound",
			errFunc:  func() error { return ErrAgentNotFound("agent-123") },
			contains: "agent not found: agent-123",
		},
		{
			name:     "AgentAlreadyExists",
			errFunc:  func() error { return ErrAgentAlreadyExists("agent-123") },
			contains: "agent already exists: agent-123",
		},
		{
			name:     "AgentNotStarted",
			errFunc:  func() error { return ErrAgentNotStarted("agent-123") },
			contains: "agent not started: agent-123",
		},
		{
			name:     "AgentTimeout",
			errFunc:  func() error { return ErrAgentTimeout("agent-123") },
			contains: "agent agent-123 operation timed out",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc()
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			if !strings.Contains(err.Error(), tt.contains) {
				t.Errorf("expected error to contain '%s', got: %s", tt.contains, err.Error())
			}
		})
	}
}

func TestLLMErrors(t *testing.T) {
	tests := []struct {
		name     string
		errFunc  func() error
		contains string
	}{
		{
			name:     "LLMProviderNotFound",
			errFunc:  func() error { return ErrLLMProviderNotFound("gemini") },
			contains: "LLM provider not found: gemini",
		},
		{
			name:     "LLMInvalidResponse",
			errFunc:  func() error { return ErrLLMInvalidResponse("openai") },
			contains: "LLM provider openai returned invalid response",
		},
		{
			name:     "LLMRateLimitExceeded",
			errFunc:  func() error { return ErrLLMRateLimitExceeded("claude") },
			contains: "LLM provider claude rate limit exceeded",
		},
		{
			name:     "LLMAPIKeyMissing",
			errFunc:  func() error { return ErrLLMAPIKeyMissing("gemini") },
			contains: "LLM provider gemini API key is missing",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc()
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			if !strings.Contains(err.Error(), tt.contains) {
				t.Errorf("expected error to contain '%s', got: %s", tt.contains, err.Error())
			}
		})
	}
}

func TestMemoryErrors(t *testing.T) {
	tests := []struct {
		name     string
		errFunc  func() error
		contains string
	}{
		{
			name:     "MemoryStoreNotInitialized",
			errFunc:  func() error { return ErrMemoryStoreNotInitialized },
			contains: "memory store not initialized",
		},
		{
			name:     "MemoryNotFound",
			errFunc:  func() error { return ErrMemoryNotFound("mem-123") },
			contains: "memory entry not found: mem-123",
		},
		{
			name:     "MemoryInvalidType",
			errFunc:  func() error { return ErrMemoryInvalidType("conversation", "task") },
			contains: "invalid memory type: expected conversation, got task",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc()
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			if !strings.Contains(err.Error(), tt.contains) {
				t.Errorf("expected error to contain '%s', got: %s", tt.contains, err.Error())
			}
		})
	}
}

func TestConfigErrors(t *testing.T) {
	baseErr := errors.New("parse error")
	
	tests := []struct {
		name     string
		errFunc  func() error
		contains string
	}{
		{
			name:     "ConfigLoadFailed",
			errFunc:  func() error { return ErrConfigLoadFailed(baseErr) },
			contains: "failed to load configuration",
		},
		{
			name:     "ConfigInvalid",
			errFunc:  func() error { return ErrConfigInvalid("port", "must be positive") },
			contains: "invalid configuration field port: must be positive",
		},
		{
			name:     "ConfigMissingRequired",
			errFunc:  func() error { return ErrConfigMissingRequired("api_key") },
			contains: "required configuration field missing: api_key",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc()
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			if !strings.Contains(err.Error(), tt.contains) {
				t.Errorf("expected error to contain '%s', got: %s", tt.contains, err.Error())
			}
		})
	}
}

func TestErrorIs(t *testing.T) {
	baseErr := errors.New("base")
	wrapped := Wrap(baseErr, "wrapped")
	
	if !Is(wrapped, baseErr) {
		t.Error("expected Is to return true for wrapped error")
	}
}

func TestGenericErrors(t *testing.T) {
	if ErrNotImplemented == nil {
		t.Error("ErrNotImplemented should not be nil")
	}
	
	if ErrTimeout == nil {
		t.Error("ErrTimeout should not be nil")
	}
	
	if ErrContextCanceled == nil {
		t.Error("ErrContextCanceled should not be nil")
	}
}
