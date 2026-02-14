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

// Package errors provides error handling utilities and domain-specific error types
// inspired by vdaas/vald error handling patterns
package errors

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
)

// Standard error wrapping functions
var (
	// Wrap wraps an error with a message
	Wrap = func(err error, msg string) error {
		if err == nil {
			return nil
		}
		return fmt.Errorf("%s: %w", msg, err)
	}

	// Wrapf wraps an error with a formatted message
	Wrapf = func(err error, format string, args ...interface{}) error {
		if err == nil {
			return nil
		}
		return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
	}

	// New creates a new error
	New = errors.New

	// Is reports whether any error in err's chain matches target
	Is = errors.Is

	// As finds the first error in err's chain that matches target
	As = errors.As

	// Unwrap returns the result of calling the Unwrap method on err
	Unwrap = errors.Unwrap

	// Join returns an error that wraps the given errors
	Join = errors.Join
)

// Agent-specific errors
var (
	// ErrAgentNotFound indicates that an agent with the specified ID was not found
	ErrAgentNotFound = func(id string) error {
		return fmt.Errorf("agent not found: %s", id)
	}

	// ErrAgentAlreadyExists indicates that an agent with the specified ID already exists
	ErrAgentAlreadyExists = func(id string) error {
		return fmt.Errorf("agent already exists: %s", id)
	}

	// ErrAgentNotStarted indicates that an agent is not started
	ErrAgentNotStarted = func(id string) error {
		return fmt.Errorf("agent not started: %s", id)
	}

	// ErrAgentAlreadyStarted indicates that an agent is already started
	ErrAgentAlreadyStarted = func(id string) error {
		return fmt.Errorf("agent already started: %s", id)
	}

	// ErrAgentTaskFailed indicates that an agent task failed
	ErrAgentTaskFailed = func(id string, err error) error {
		return Wrapf(err, "agent %s task failed", id)
	}

	// ErrAgentTimeout indicates that an agent operation timed out
	ErrAgentTimeout = func(id string) error {
		return fmt.Errorf("agent %s operation timed out", id)
	}
)

// LLM-specific errors
var (
	// ErrLLMProviderNotFound indicates that an LLM provider was not found
	ErrLLMProviderNotFound = func(name string) error {
		return fmt.Errorf("LLM provider not found: %s", name)
	}

	// ErrLLMGenerationFailed indicates that LLM generation failed
	ErrLLMGenerationFailed = func(provider string, err error) error {
		return Wrapf(err, "LLM generation failed for provider %s", provider)
	}

	// ErrLLMInvalidResponse indicates that LLM returned an invalid response
	ErrLLMInvalidResponse = func(provider string) error {
		return fmt.Errorf("LLM provider %s returned invalid response", provider)
	}

	// ErrLLMRateLimitExceeded indicates that LLM rate limit was exceeded
	ErrLLMRateLimitExceeded = func(provider string) error {
		return fmt.Errorf("LLM provider %s rate limit exceeded", provider)
	}

	// ErrLLMAPIKeyMissing indicates that LLM API key is missing
	ErrLLMAPIKeyMissing = func(provider string) error {
		return fmt.Errorf("LLM provider %s API key is missing", provider)
	}
)

// Memory-specific errors
var (
	// ErrMemoryStoreNotInitialized indicates that memory store is not initialized
	ErrMemoryStoreNotInitialized = errors.New("memory store not initialized")

	// ErrMemoryStoreFailed indicates that a memory store operation failed
	ErrMemoryStoreFailed = func(operation string, err error) error {
		return Wrapf(err, "memory store %s operation failed", operation)
	}

	// ErrMemoryNotFound indicates that a memory entry was not found
	ErrMemoryNotFound = func(id string) error {
		return fmt.Errorf("memory entry not found: %s", id)
	}

	// ErrMemoryInvalidType indicates that a memory entry has an invalid type
	ErrMemoryInvalidType = func(expected, actual string) error {
		return fmt.Errorf("invalid memory type: expected %s, got %s", expected, actual)
	}
)

// Configuration errors
var (
	// ErrConfigLoadFailed indicates that configuration loading failed
	ErrConfigLoadFailed = func(err error) error {
		return Wrap(err, "failed to load configuration")
	}

	// ErrConfigInvalid indicates that configuration is invalid
	ErrConfigInvalid = func(field string, reason string) error {
		return fmt.Errorf("invalid configuration field %s: %s", field, reason)
	}

	// ErrConfigMissingRequired indicates that a required configuration field is missing
	ErrConfigMissingRequired = func(field string) error {
		return fmt.Errorf("required configuration field missing: %s", field)
	}
)

// Communication errors
var (
	// ErrGRPCConnectionFailed indicates that gRPC connection failed
	ErrGRPCConnectionFailed = func(addr string, err error) error {
		return Wrapf(err, "failed to connect to gRPC server at %s", addr)
	}

	// ErrGRPCCallFailed indicates that a gRPC call failed
	ErrGRPCCallFailed = func(method string, err error) error {
		return Wrapf(err, "gRPC call to %s failed", method)
	}

	// ErrSlackNotificationFailed indicates that Slack notification failed
	ErrSlackNotificationFailed = func(err error) error {
		return Wrap(err, "failed to send Slack notification")
	}
)

// Option errors (inspired by Vald's option pattern)
var (
	// ErrOptionFailed indicates that an option function failed
	ErrOptionFailed = func(err error, ref reflect.Value) error {
		var str string
		if ref.IsValid() && ref.Kind() == reflect.Func {
			str = runtime.FuncForPC(ref.Pointer()).Name()
		}
		if str != "" {
			return Wrapf(err, "failed to setup option: %s", str)
		}
		return Wrap(err, "failed to setup option")
	}

	// ErrInvalidOption indicates that an option value is invalid
	ErrInvalidOption = func(name string, value interface{}) error {
		return fmt.Errorf("invalid option %s: %v", name, value)
	}
)

// Type conversion errors
var (
	// ErrInvalidTypeConversion indicates that a type conversion is invalid
	ErrInvalidTypeConversion = func(from, to string) error {
		return fmt.Errorf("invalid type conversion from %s to %s", from, to)
	}
)

// Generic errors
var (
	// ErrNotImplemented indicates that a feature is not implemented
	ErrNotImplemented = errors.New("not implemented")

	// ErrInvalidArgument indicates that an argument is invalid
	ErrInvalidArgument = func(name string, value interface{}) error {
		return fmt.Errorf("invalid argument %s: %v", name, value)
	}

	// ErrTimeout indicates that an operation timed out
	ErrTimeout = errors.New("operation timed out")

	// ErrContextCanceled indicates that the context was canceled
	ErrContextCanceled = errors.New("context canceled")
)
