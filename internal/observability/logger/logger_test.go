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

package logger

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{LevelFatal, "FATAL"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.level.String(); got != tt.want {
				t.Errorf("Level.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStandardLogger(t *testing.T) {
	logger := NewStandardLogger()
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestStandardLogger_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewStandardLogger(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithFormat(FormatText),
	)

	logger.Info("test message", String("key", "value"))

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("expected output to contain 'INFO', got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, "key=value") {
		t.Errorf("expected output to contain 'key=value', got: %s", output)
	}
}

func TestStandardLogger_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewStandardLogger(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithFormat(FormatJSON),
	)

	logger.Info("test message", String("key", "value"))

	output := buf.String()
	if !strings.Contains(output, `"level":"INFO"`) {
		t.Errorf("expected JSON output to contain level, got: %s", output)
	}
	if !strings.Contains(output, `"message":"test message"`) {
		t.Errorf("expected JSON output to contain message, got: %s", output)
	}
	if !strings.Contains(output, `"key":"value"`) {
		t.Errorf("expected JSON output to contain key=value, got: %s", output)
	}
}

func TestStandardLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewStandardLogger(
		WithOutput(&buf),
		WithLevel(LevelWarn),
	)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")

	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Error("debug message should be filtered out")
	}
	if strings.Contains(output, "info message") {
		t.Error("info message should be filtered out")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("warn message should not be filtered out")
	}
}

func TestStandardLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewStandardLogger(
		WithOutput(&buf),
		WithLevel(LevelInfo),
	)

	childLogger := logger.WithFields(String("component", "test"))
	childLogger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "component=test") {
		t.Errorf("expected output to contain component field, got: %s", output)
	}
}

func TestField_Constructors(t *testing.T) {
	tests := []struct {
		name     string
		field    Field
		wantKey  string
		hasValue bool
	}{
		{
			name:     "String field",
			field:    String("key", "value"),
			wantKey:  "key",
			hasValue: true,
		},
		{
			name:     "Int field",
			field:    Int("count", 42),
			wantKey:  "count",
			hasValue: true,
		},
		{
			name:     "Error field",
			field:    Error(errors.New("test error")),
			wantKey:  "error",
			hasValue: true,
		},
		{
			name:     "Any field",
			field:    Any("data", map[string]string{"key": "value"}),
			wantKey:  "data",
			hasValue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field.Key != tt.wantKey {
				t.Errorf("field.Key = %v, want %v", tt.field.Key, tt.wantKey)
			}
			if tt.hasValue && tt.field.Value == nil {
				t.Error("expected non-nil value")
			}
		})
	}
}

func TestGlobalLogger(t *testing.T) {
	var buf bytes.Buffer
	Init(
		WithOutput(&buf),
		WithLevel(LevelInfo),
	)

	Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got: %s", output)
	}
}

func TestStandardLogger_AllLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := NewStandardLogger(
		WithOutput(&buf),
		WithLevel(LevelDebug),
	)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	expectedMessages := []string{"debug message", "info message", "warn message", "error message"}
	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("expected output to contain '%s', got: %s", msg, output)
		}
	}
}

func TestStandardLogger_Close(t *testing.T) {
	logger := NewStandardLogger()
	if err := logger.Close(); err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}
